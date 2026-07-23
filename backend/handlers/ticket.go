package handlers

import (
	"crypto/rand"
	"encoding/csv"
	"fmt"
	"math/big"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"ticket-system/config"
	"ticket-system/middleware"
	"ticket-system/models"
)

// Handler 持有 DB 和配置
type Handler struct {
	DB  *gorm.DB
	Cfg *config.Config
}

func New(db *gorm.DB, cfg *config.Config) *Handler {
	return &Handler{DB: db, Cfg: cfg}
}

// ─── DTO ───

type SubmitRequest struct {
	Question  string     `json:"question" binding:"required,max=500"`
	UserID    string     `json:"user_id" binding:"omitempty,max=128"`
	Timestamp *time.Time `json:"timestamp"`
	Source    string     `json:"source" binding:"omitempty,oneof=hiagent_chat wechat_service wechat_subscribe feishu yiban"`
	RAGResult string     `json:"rag_result" binding:"max=2000"`
}

type SubmitResponse struct {
	TicketID string `json:"ticket_id"`
	Status   string `json:"status"`
	Message  string `json:"message"`
}

type AnswerRequest struct {
	Answer   string `json:"answer" binding:"required,max=5000"`
	Operator string `json:"operator" binding:"omitempty,max=64"`
	SyncToKB bool   `json:"sync_to_kb"`
}

type CloseRequest struct {
	Reason string `json:"reason" binding:"max=500"`
}

type AnswerResponse struct {
	Success  bool   `json:"success"`
	TicketID string `json:"ticket_id"`
	Message  string `json:"message"`
}

type ListResponse struct {
	Total    int64            `json:"total"`
	Page     int              `json:"page"`
	PageSize int              `json:"page_size"`
	Items    []models.Ticket  `json:"items"`
}

type BatchDeleteRequest struct {
	TicketIDs []string `json:"ticket_ids" binding:"required,min=1,max=100"`
}

type ArchiveRequest struct {
	Category string `json:"category" binding:"omitempty,max=32"`
}

// ─── 工具函数 ───

const ticketChars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// generateTicketID 生成 T20260723-XXXXXX 格式工单号
func generateTicketID(prefix string) string {
	dateStr := time.Now().Format("20060102")
	b := make([]byte, 6)
	max := big.NewInt(int64(len(ticketChars)))
	for i := range b {
		n, _ := rand.Int(rand.Reader, max)
		b[i] = ticketChars[n.Int64()]
	}
	return fmt.Sprintf("%s%s-%s", prefix, dateStr, string(b))
}

// notDeleted 排除软删除的工单
func notDeleted(db *gorm.DB) *gorm.DB {
	return db.Where("deleted_at IS NULL")
}

// ─── 端点实现 ───

// SubmitTicket 提交工单
func (h *Handler) SubmitTicket(c *gin.Context) {
	var req SubmitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error_code":    "VALIDATION_ERROR",
			"error_message": "Request validation failed",
			"detail":        err.Error(),
		})
		return
	}

	if req.Source == "" {
		req.Source = "hiagent_chat"
	}
	if req.UserID == "" {
		req.UserID = "anonymous"
	}

	// 5 分钟内同一 user+question 去重
	fiveMinAgo := time.Now().Add(-5 * time.Minute)
	var recent models.Ticket
	err := h.DB.Where("user_id = ? AND question = ? AND created_at >= ? AND deleted_at IS NULL",
		req.UserID, req.Question, fiveMinAgo).
		First(&recent).Error
	if err == nil {
		c.JSON(http.StatusOK, SubmitResponse{
			TicketID: recent.TicketID,
			Status:   "created",
			Message:  fmt.Sprintf("已有工单（%s），请勿重复提交", recent.TicketID),
		})
		return
	}

	t := models.Ticket{
		TicketID:  generateTicketID(h.Cfg.TicketPrefix),
		Question:  strings.TrimSpace(req.Question),
		UserID:    req.UserID,
		Source:    req.Source,
		RAGResult: req.RAGResult,
		Status:    "pending",
	}
	if req.Timestamp != nil {
		t.CreatedAt = *req.Timestamp
	}

	if err := h.DB.Create(&t).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error_code":    "DB_ERROR",
			"error_message": "Failed to create ticket",
			"detail":        err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SubmitResponse{
		TicketID: t.TicketID,
		Status:   "created",
		Message:  fmt.Sprintf("工单已提交（%s），请等待答复", t.TicketID),
	})
}

// ListTickets 查询工单列表（增强筛选：archived / category / search）
func (h *Handler) ListTickets(c *gin.Context) {
	page := atoiDefault(c.Query("page"), 1)
	pageSize := atoiDefault(c.Query("page_size"), 20)
	if pageSize > 100 {
		pageSize = 100
	}
	if pageSize < 1 {
		pageSize = 20
	}
	if page < 1 {
		page = 1
	}

	q := h.DB.Model(&models.Ticket{}).Where("deleted_at IS NULL")

	if s := c.Query("status"); s != "" {
		q = q.Where("status = ?", s)
	}
	if src := c.Query("source"); src != "" {
		q = q.Where("source = ?", src)
	}
	if sd := c.Query("start_date"); sd != "" {
		if t, err := time.Parse("2006-01-02", sd); err == nil {
			q = q.Where("created_at >= ?", t)
		}
	}
	if ed := c.Query("end_date"); ed != "" {
		if t, err := time.Parse("2006-01-02", ed); err == nil {
			q = q.Where("created_at < ?", t.Add(24*time.Hour))
		}
	}
	// 新增筛选
	if arch := c.Query("archived"); arch != "" {
		q = q.Where("archived = ?", arch == "true")
	}
	if cat := c.Query("category"); cat != "" {
		q = q.Where("category = ?", cat)
	}
	if kw := c.Query("search"); kw != "" {
		like := "%" + kw + "%"
		q = q.Where("question LIKE ? OR ticket_id LIKE ?", like, like)
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var items []models.Ticket
	if err := q.Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&items).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, ListResponse{
		Total:    total,
		Page:     page,
		PageSize: pageSize,
		Items:    items,
	})
}

// GetTicket 查询单个工单详情
func (h *Handler) GetTicket(c *gin.Context) {
	id := c.Param("ticket_id")
	var t models.Ticket
	if err := h.DB.Where("ticket_id = ? AND deleted_at IS NULL", id).First(&t).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error_code":    "NOT_FOUND",
				"error_message": "Ticket not found",
				"detail":        "工单号 " + id + " 不存在或已删除",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, t)
}

// AnswerTicket 人工答复工单（已答复可覆盖重答）
func (h *Handler) AnswerTicket(c *gin.Context) {
	id := c.Param("ticket_id")
	var req AnswerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error_code":    "VALIDATION_ERROR",
			"error_message": "Request validation failed",
			"detail":        err.Error(),
		})
		return
	}

	var t models.Ticket
	if err := h.DB.Where("ticket_id = ? AND deleted_at IS NULL", id).First(&t).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error_code":    "NOT_FOUND",
				"error_message": "Ticket not found",
				"detail":        "工单号 " + id + " 不存在或已删除",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if t.Status == "closed" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error_code":    "TICKET_CLOSED",
			"error_message": "Ticket is closed",
			"detail":        "工单已关闭，无法回答",
		})
		return
	}

	// 优先用 JWT 当前用户的 display_name，其次用 req.Operator
	operatorName := req.Operator
	if u := middleware.CurrentUser(c); u != nil {
		if u.DisplayName != "" {
			operatorName = u.DisplayName
		} else {
			operatorName = u.Username
		}
	}

	now := time.Now()
	answer := req.Answer
	t.Answer = &answer
	t.AnsweredBy = &operatorName
	t.AnsweredAt = &now
	t.Status = "answered"

	if err := h.DB.Save(&t).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, AnswerResponse{
		Success:  true,
		TicketID: t.TicketID,
		Message:  fmt.Sprintf("工单 %s 已答复", t.TicketID),
	})
}

// CloseTicket 关闭工单
func (h *Handler) CloseTicket(c *gin.Context) {
	id := c.Param("ticket_id")
	var req CloseRequest
	_ = c.ShouldBindJSON(&req)

	var t models.Ticket
	if err := h.DB.Where("ticket_id = ? AND deleted_at IS NULL", id).First(&t).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error_code":    "NOT_FOUND",
				"error_message": "Ticket not found",
				"detail":        "工单号 " + id + " 不存在或已删除",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	t.Status = "closed"
	if req.Reason != "" {
		reason := req.Reason
		t.CloseReason = &reason
	}

	if err := h.DB.Save(&t).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, AnswerResponse{
		Success:  true,
		TicketID: t.TicketID,
		Message:  fmt.Sprintf("工单 %s 已关闭", t.TicketID),
	})
}

// ─── 新增端点 ───

// ArchiveTicket 归档工单（可同时设置分类）
func (h *Handler) ArchiveTicket(c *gin.Context) {
	id := c.Param("ticket_id")
	var req ArchiveRequest
	_ = c.ShouldBindJSON(&req)

	var t models.Ticket
	if err := h.DB.Where("ticket_id = ? AND deleted_at IS NULL", id).First(&t).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error_code":    "NOT_FOUND",
				"error_message": "Ticket not found",
				"detail":        "工单号 " + id + " 不存在或已删除",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	now := time.Now()
	t.Archived = true
	t.ArchivedAt = &now
	if req.Category != "" {
		t.Category = req.Category
	}

	if err := h.DB.Save(&t).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, AnswerResponse{
		Success:  true,
		TicketID: t.TicketID,
		Message:  fmt.Sprintf("工单 %s 已归档", t.TicketID),
	})
}

// UnarchiveTicket 取消归档
func (h *Handler) UnarchiveTicket(c *gin.Context) {
	id := c.Param("ticket_id")

	var t models.Ticket
	if err := h.DB.Where("ticket_id = ? AND deleted_at IS NULL", id).First(&t).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error_code":    "NOT_FOUND",
				"error_message": "Ticket not found",
				"detail":        "工单号 " + id + " 不存在或已删除",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	t.Archived = false
	t.ArchivedAt = nil

	if err := h.DB.Save(&t).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, AnswerResponse{
		Success:  true,
		TicketID: t.TicketID,
		Message:  fmt.Sprintf("工单 %s 已取消归档", t.TicketID),
	})
}

// DeleteTicket 软删除单个工单
func (h *Handler) DeleteTicket(c *gin.Context) {
	id := c.Param("ticket_id")

	var t models.Ticket
	if err := h.DB.Where("ticket_id = ? AND deleted_at IS NULL", id).First(&t).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error_code":    "NOT_FOUND",
				"error_message": "Ticket not found",
				"detail":        "工单号 " + id + " 不存在或已删除",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	now := time.Now()
	t.DeletedAt = &now
	if err := h.DB.Save(&t).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, AnswerResponse{
		Success:  true,
		TicketID: t.TicketID,
		Message:  fmt.Sprintf("工单 %s 已删除", t.TicketID),
	})
}

// BatchDelete 批量软删除工单
func (h *Handler) BatchDelete(c *gin.Context) {
	var req BatchDeleteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error_code":    "VALIDATION_ERROR",
			"error_message": "Invalid batch payload",
			"detail":        err.Error(),
		})
		return
	}

	now := time.Now()
	result := h.DB.Model(&models.Ticket{}).
		Where("ticket_id IN ? AND deleted_at IS NULL", req.TicketIDs).
		Update("deleted_at", now)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error_code":    "DB_ERROR",
			"error_message": "Failed to delete tickets",
			"detail":        result.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":      true,
		"deleted":      result.RowsAffected,
		"total_request": len(req.TicketIDs),
		"message":      fmt.Sprintf("已删除 %d 个工单", result.RowsAffected),
	})
}

// ExportTickets 导出工单（CSV 或 JSON，基于当前筛选条件）
func (h *Handler) ExportTickets(c *gin.Context) {
	format := c.Query("format")
	if format != "csv" && format != "json" {
		format = "json" // 默认 JSON
	}

	q := h.DB.Model(&models.Ticket{}).Where("deleted_at IS NULL")

	if s := c.Query("status"); s != "" {
		q = q.Where("status = ?", s)
	}
	if src := c.Query("source"); src != "" {
		q = q.Where("source = ?", src)
	}
	if sd := c.Query("start_date"); sd != "" {
		if t, err := time.Parse("2006-01-02", sd); err == nil {
			q = q.Where("created_at >= ?", t)
		}
	}
	if ed := c.Query("end_date"); ed != "" {
		if t, err := time.Parse("2006-01-02", ed); err == nil {
			q = q.Where("created_at < ?", t.Add(24*time.Hour))
		}
	}
	if arch := c.Query("archived"); arch != "" {
		q = q.Where("archived = ?", arch == "true")
	}
	if cat := c.Query("category"); cat != "" {
		q = q.Where("category = ?", cat)
	}
	if kw := c.Query("search"); kw != "" {
		like := "%" + kw + "%"
		q = q.Where("question LIKE ? OR ticket_id LIKE ?", like, like)
	}

	var items []models.Ticket
	if err := q.Order("created_at DESC").Find(&items).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if format == "csv" {
		c.Header("Content-Type", "text/csv; charset=utf-8")
		c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=tickets_%s.csv", time.Now().Format("20060102_150405")))

		// 写入 BOM 确保 Excel 正确识别 UTF-8
		c.Writer.Write([]byte{0xEF, 0xBB, 0xBF})

		writer := csv.NewWriter(c.Writer)
		writer.Write([]string{"工单号", "问题", "用户ID", "来源", "状态", "分类", "已归档", "答复内容", "答复人", "创建时间", "答复时间", "关闭原因"})

		for _, t := range items {
			answer := ""
			if t.Answer != nil {
				answer = *t.Answer
			}
			answeredBy := ""
			if t.AnsweredBy != nil {
				answeredBy = *t.AnsweredBy
			}
			closeReason := ""
			if t.CloseReason != nil {
				closeReason = *t.CloseReason
			}
			archivedStr := "否"
			if t.Archived {
				archivedStr = "是"
			}
			answeredAt := ""
			if t.AnsweredAt != nil {
				answeredAt = t.AnsweredAt.Format("2006-01-02 15:04:05")
			}

			writer.Write([]string{
				t.TicketID,
				t.Question,
				t.UserID,
				t.Source,
				t.Status,
				t.Category,
				archivedStr,
				answer,
				answeredBy,
				t.CreatedAt.Format("2006-01-02 15:04:05"),
				answeredAt,
				closeReason,
			})
		}
		writer.Flush()
		return
	}

	// JSON 导出
	c.Header("Content-Type", "application/json; charset=utf-8")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=tickets_%s.json", time.Now().Format("20060102_150405")))
	c.JSON(http.StatusOK, gin.H{
		"total":   len(items),
		"exported_at": time.Now().UTC().Format(time.RFC3339),
		"items":   items,
	})
}

// PatchTicket 部分更新工单（分类等字段）
func (h *Handler) PatchTicket(c *gin.Context) {
	id := c.Param("ticket_id")
	var req ArchiveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error_code":    "VALIDATION_ERROR",
			"error_message": "Invalid payload",
			"detail":        err.Error(),
		})
		return
	}

	var t models.Ticket
	if err := h.DB.Where("ticket_id = ? AND deleted_at IS NULL", id).First(&t).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error_code":    "NOT_FOUND",
				"error_message": "Ticket not found",
				"detail":        "工单号 " + id + " 不存在或已删除",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if req.Category != "" {
		t.Category = req.Category
	}

	if err := h.DB.Save(&t).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, AnswerResponse{
		Success:  true,
		TicketID: t.TicketID,
		Message:  fmt.Sprintf("工单 %s 已更新", t.TicketID),
	})
}

// ReopenTicket 重新打开已关闭的工单（状态恢复为 answered）
func (h *Handler) ReopenTicket(c *gin.Context) {
	id := c.Param("ticket_id")

	var t models.Ticket
	if err := h.DB.Where("ticket_id = ? AND deleted_at IS NULL", id).First(&t).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error_code":    "NOT_FOUND",
				"error_message": "Ticket not found",
				"detail":        "工单号 " + id + " 不存在或已删除",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if t.Status != "closed" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error_code":    "TICKET_NOT_CLOSED",
			"error_message": "Only closed tickets can be reopened",
			"detail":        "只能重新打开已关闭的工单",
		})
		return
	}

	t.Status = "answered"
	t.CloseReason = nil

	if err := h.DB.Save(&t).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, AnswerResponse{
		Success:  true,
		TicketID: t.TicketID,
		Message:  fmt.Sprintf("工单 %s 已重新打开", t.TicketID),
	})
}

func atoiDefault(s string, def int) int {
	if s == "" {
		return def
	}
	var n int
	for _, ch := range s {
		if ch < '0' || ch > '9' {
			return def
		}
		n = n*10 + int(ch-'0')
	}
	if n <= 0 {
		return def
	}
	return n
}
