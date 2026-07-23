package handlers

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"ticket-system/middleware"
	"ticket-system/models"
)

// UserHandler 用户管理 handler
type UserHandler struct {
	DB *gorm.DB
}

func NewUser(db *gorm.DB) *UserHandler {
	return &UserHandler{DB: db}
}

// ─── DTO ───

type CreateUserRequest struct {
	Username    string `json:"username" binding:"required,min=3,max=64"`
	Password    string `json:"password" binding:"required,min=8,max=128"`
	Role        string `json:"role" binding:"required,oneof=admin staff"`
	DisplayName string `json:"display_name" binding:"max=64"`
	Email       string `json:"email" binding:"omitempty,email,max=128"`
}

type UpdateUserRequest struct {
	Role        *string `json:"role" binding:"omitempty,oneof=admin staff"`
	DisplayName *string `json:"display_name" binding:"omitempty,max=64"`
	Email       *string `json:"email" binding:"omitempty,email,max=128"`
	Active      *bool   `json:"active"`
	Password    string  `json:"password" binding:"omitempty,min=8,max=128"`
}

type UserListResponse struct {
	Total int64       `json:"total"`
	Items []UserPublic `json:"items"`
}

// ─── Handlers ───

// ListUsers 用户列表
func (h *UserHandler) ListUsers(c *gin.Context) {
	var users []models.User
	if err := h.DB.Order("created_at DESC").Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	items := make([]UserPublic, 0, len(users))
	for i := range users {
		items = append(items, toPublic(&users[i]))
	}
	c.JSON(http.StatusOK, UserListResponse{
		Total: int64(len(items)),
		Items: items,
	})
}

// CreateUser 新建用户
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error_code":    "VALIDATION_ERROR",
			"error_message": "Invalid payload",
			"detail":        err.Error(),
		})
		return
	}

	username := strings.TrimSpace(req.Username)
	// 唯一性
	var existing models.User
	if err := h.DB.Where("username = ?", username).First(&existing).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{
			"error_code":    "USERNAME_TAKEN",
			"error_message": "Username already exists",
		})
		return
	}

	if err := validatePasswordStrength(req.Password); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error_code":    "WEAK_PASSWORD",
			"error_message": err.Error(),
		})
		return
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	u := models.User{
		Username:      username,
		PasswordHash:  string(hash),
		Role:          req.Role,
		DisplayName:   req.DisplayName,
		Email:         req.Email,
		Active:        true,
		MustChangePwd: true, // 强制首次改密码
	}
	if err := h.DB.Create(&u).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"user":    toPublic(&u),
		"message": "User created. They must change password on first login.",
	})
}

// UpdateUser 编辑用户（admin 给自己也能用，包括改密码）
func (h *UserHandler) UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var u models.User
	if err := h.DB.Where("id = ?", id).First(&u).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error_code":    "NOT_FOUND",
			"error_message": "User not found",
		})
		return
	}

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error_code":    "VALIDATION_ERROR",
			"error_message": "Invalid payload",
			"detail":        err.Error(),
		})
		return
	}

	currentUser := middleware.CurrentUser(c)

	if req.Role != nil {
		// 防止最后一个 admin 被降级
		if u.Role == "admin" && *req.Role != "admin" {
			var adminCount int64
			h.DB.Model(&models.User{}).Where("role = ? AND active = ?", "admin", true).Count(&adminCount)
			if adminCount <= 1 {
				c.JSON(http.StatusBadRequest, gin.H{
					"error_code":    "LAST_ADMIN",
					"error_message": "Cannot demote the only active admin",
				})
				return
			}
		}
		u.Role = *req.Role
	}
	if req.DisplayName != nil {
		u.DisplayName = *req.DisplayName
	}
	if req.Email != nil {
		u.Email = *req.Email
	}
	if req.Active != nil {
		// 不能禁用自己
		if currentUser.ID == u.ID && !*req.Active {
			c.JSON(http.StatusBadRequest, gin.H{
				"error_code":    "SELF_DISABLE",
				"error_message": "Cannot disable your own account",
			})
			return
		}
		// 不能禁用最后一个 admin
		if u.Role == "admin" && u.Active && !*req.Active {
			var adminCount int64
			h.DB.Model(&models.User{}).Where("role = ? AND active = ?", "admin", true).Count(&adminCount)
			if adminCount <= 1 {
				c.JSON(http.StatusBadRequest, gin.H{
					"error_code":    "LAST_ADMIN",
					"error_message": "Cannot disable the only active admin",
				})
				return
			}
		}
		u.Active = *req.Active
	}
	if req.Password != "" {
		if err := validatePasswordStrength(req.Password); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error_code":    "WEAK_PASSWORD",
				"error_message": err.Error(),
			})
			return
		}
		hash, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		u.PasswordHash = string(hash)
		// 管理员重置别人密码时，强制下次登录改密码
		if currentUser.ID != u.ID {
			u.MustChangePwd = true
		}
	}

	u.UpdatedAt = time.Now()
	if err := h.DB.Save(&u).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"user":    toPublic(&u),
		"message": "User updated",
	})
}
