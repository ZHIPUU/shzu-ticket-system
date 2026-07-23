package handlers

import (
	"crypto/rand"
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
	Question  string `json:"question" binding:"required,max=500"`
	UserID    string `json:"user_id" binding:"required,max=128"`
	Timestamp *time.Time `json:"timestamp"`
	Source    string `json:"source" binding:"omitempty,oneof=hiagent_chat wechat_service wechat_subscribe feishu yiban"`
	RAGResult string `json:"rag_result" binding:"max=2000"`
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

// ─── 工具函数 ───

const ticketChars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// generateTicketID 生成 T20260723-XXXXXX 格式工单号
func generateTicketID(prefix string) string {
	dateStr := time.Now().Format("20060102")
	// 6 位随机字符
	b := make([]byte, 6)
	max := big.NewInt(int64(len(ticketChars)))
	for i := range b {
		n, _ := rand.Int(rand.Reader, max)
		b[i] = ticketChars[n.Int64()]
	}
	return fmt.Sprintf("%s%s-%s", prefix, dateStr, string(b))
}

// formatTimestamp 统一时间格式
func formatTimestamp(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.UTC().Format(time.RFC3339)
}

// ─── 端点实现 ───

// SubmitTicket 提交工单
// @Summary 提交工单（HiAgent 插件调用）
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

	// 默认值
	if req.Source == "" {
		req.Source = "hiagent_chat"
	}

	// 5 分钟内同一 user+question 去重
	fiveMinAgo := time.Now().Add(-5 * time.Minute)
	var recent models.Ticket
	err := h.DB.Where("user_id = ? AND question = ? AND created_at >= ?",
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

	// 创建工单
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

// ListTickets 查询工单列表
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

	q := h.DB.Model(&models.Ticket{})
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
			q = q.Where("created_at < ?", t)
		}
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
	if err := h.DB.Where("ticket_id = ?", id).First(&t).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error_code":    "NOT_FOUND",
				"error_message": "Ticket not found",
				"detail":        "工单号 " + id + " 不存在",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, t)
}

// AnswerTicket 人工答复工单
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
	if err := h.DB.Where("ticket_id = ?", id).First(&t).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error_code":    "NOT_FOUND",
				"error_message": "Ticket not found",
				"detail":        "工单号 " + id + " 不存在",
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

	// TODO: 当 req.SyncToKB 为 true 时，调用 HiAgent 知识库 API
	// 当前 phase 留空，由后续阶段接入

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
	// requestBody 是可选的，bind 失败不报错
	_ = c.ShouldBindJSON(&req)

	var t models.Ticket
	if err := h.DB.Where("ticket_id = ?", id).First(&t).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error_code":    "NOT_FOUND",
				"error_message": "Ticket not found",
				"detail":        "工单号 " + id + " 不存在",
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
