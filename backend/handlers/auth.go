package handlers

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"ticket-system/config"
	"ticket-system/middleware"
	"ticket-system/models"
)

// AuthHandler 认证相关 handler
type AuthHandler struct {
	DB  *gorm.DB
	Cfg *config.Config
}

func NewAuth(db *gorm.DB, cfg *config.Config) *AuthHandler {
	return &AuthHandler{DB: db, Cfg: cfg}
}

// ─── DTO ───

type LoginRequest struct {
	Username string `json:"username" binding:"required,min=3,max=64"`
	Password string `json:"password" binding:"required,min=6,max=128"`
}

type LoginResponse struct {
	Token             string    `json:"token"`
	User              UserPublic `json:"user"`
	ExpiresAt         time.Time `json:"expires_at"`
	MustChangePassword bool     `json:"must_change_password"`
}

type UserPublic struct {
	ID                uint       `json:"id"`
	Username          string     `json:"username"`
	Role              string     `json:"role"`
	DisplayName       string     `json:"display_name"`
	Email             string     `json:"email"`
	Active            bool       `json:"active"`
	MustChangePassword bool      `json:"must_change_password"`
	LastLoginAt       *time.Time `json:"last_login_at,omitempty"`
	CreatedAt         time.Time  `json:"created_at"`
}

func toPublic(u *models.User) UserPublic {
	return UserPublic{
		ID:                u.ID,
		Username:          u.Username,
		Role:              u.Role,
		DisplayName:       u.DisplayName,
		Email:             u.Email,
		Active:            u.Active,
		MustChangePassword: u.MustChangePwd,
		LastLoginAt:       u.LastLoginAt,
		CreatedAt:         u.CreatedAt,
	}
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required,min=6"`
	NewPassword string `json:"new_password" binding:"required,min=8,max=128"`
}

// ─── Handlers ───

// Login 登录
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error_code":    "VALIDATION_ERROR",
			"error_message": "Invalid login payload",
			"detail":        err.Error(),
		})
		return
	}

	var u models.User
	if err := h.DB.Where("username = ?", strings.TrimSpace(req.Username)).First(&u).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error_code":    "INVALID_CREDENTIALS",
			"error_message": "Username or password incorrect",
		})
		return
	}

	if !u.Active {
		c.JSON(http.StatusForbidden, gin.H{
			"error_code":    "USER_DISABLED",
			"error_message": "User account is disabled",
		})
		return
	}

	// 校验密码
	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error_code":    "INVALID_CREDENTIALS",
			"error_message": "Username or password incorrect",
		})
		return
	}

	// 更新最后登录时间
	now := time.Now()
	u.LastLoginAt = &now
	h.DB.Save(&u)

	// 生成 JWT
	exp := time.Now().Add(h.Cfg.JWTExpiry)
	claims := &middleware.Claims{
		UserID:   u.ID,
		Username: u.Username,
		Role:     u.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(exp),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "ticket-system",
			Subject:   u.Username,
		},
	}
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := tok.SignedString([]byte(h.Cfg.JWTSecret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error_code":    "TOKEN_SIGN_FAILED",
			"error_message": "Failed to sign token",
			"detail":        err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, LoginResponse{
		Token:              signed,
		User:               toPublic(&u),
		ExpiresAt:          exp,
		MustChangePassword: u.MustChangePwd,
	})
}

// Me 获取当前登录用户信息
func (h *AuthHandler) Me(c *gin.Context) {
	u := middleware.CurrentUser(c)
	if u == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error_code": "UNAUTHORIZED"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"user": toPublic(u)})
}

// ChangePassword 修改自己密码
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	u := middleware.CurrentUser(c)
	if u == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error_code": "UNAUTHORIZED"})
		return
	}

	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error_code":    "VALIDATION_ERROR",
			"error_message": "Invalid payload",
			"detail":        err.Error(),
		})
		return
	}

	if req.OldPassword == req.NewPassword {
		c.JSON(http.StatusBadRequest, gin.H{
			"error_code":    "PASSWORD_REUSED",
			"error_message": "New password must be different from old password",
		})
		return
	}

	// 校验强度
	if err := validatePasswordStrength(req.NewPassword); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error_code":    "WEAK_PASSWORD",
			"error_message": err.Error(),
		})
		return
	}

	// 校验旧密码
	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(req.OldPassword)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error_code":    "OLD_PASSWORD_WRONG",
			"error_message": "Old password is incorrect",
		})
		return
	}

	// 哈希新密码
	hash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	u.PasswordHash = string(hash)
	u.MustChangePwd = false
	h.DB.Save(u)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Password changed successfully",
	})
}

// validatePasswordStrength 密码强度校验
func validatePasswordStrength(p string) error {
	if len(p) < 8 {
		return errWeak("Password must be at least 8 characters")
	}
	hasLetter := false
	hasDigit := false
	for _, r := range p {
		switch {
		case r >= 'a' && r <= 'z', r >= 'A' && r <= 'Z':
			hasLetter = true
		case r >= '0' && r <= '9':
			hasDigit = true
		}
	}
	if !hasLetter || !hasDigit {
		return errWeak("Password must contain both letters and digits")
	}
	return nil
}

type errWeak string

func (e errWeak) Error() string { return string(e) }
