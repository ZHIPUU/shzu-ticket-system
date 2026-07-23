package middleware

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"

	"ticket-system/config"
	"ticket-system/models"
)

// ContextKey 上下文键
const (
	CtxUser   = "current_user"
	CtxUserID = "current_user_id"
)

// Claims JWT 载荷
type Claims struct {
	UserID   uint   `json:"uid"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// JWTAuth 验证 JWT Token
func JWTAuth(cfg *config.Config, db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从 Authorization Header 取 Bearer Token
		auth := c.GetHeader("Authorization")
		if !strings.HasPrefix(auth, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error_code":    "UNAUTHORIZED",
				"error_message": "Missing or invalid Authorization header",
				"detail":        "Use: Authorization: Bearer <token>",
			})
			return
		}
		raw := strings.TrimPrefix(auth, "Bearer ")
		claims := &Claims{}
		_, err := jwt.ParseWithClaims(raw, claims, func(t *jwt.Token) (interface{}, error) {
			return []byte(cfg.JWTSecret), nil
		})
		if err != nil {
			log.Printf("[JWTAuth] parse error: %v (auth header len=%d)", err, len(auth))
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error_code":    "INVALID_TOKEN",
				"error_message": "Token invalid or expired",
				"detail":        err.Error(),
			})
			return
		}

		// 加载用户，确认账号有效
		var u models.User
		if err := db.First(&u, claims.UserID).Error; err != nil {
			log.Printf("[JWTAuth] user not found, userID=%d", claims.UserID)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error_code":    "USER_NOT_FOUND",
				"error_message": "User not found",
			})
			return
		}
		if !u.Active {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error_code":    "USER_DISABLED",
				"error_message": "User account is disabled",
			})
			return
		}

		c.Set(CtxUser, &u)
		c.Set(CtxUserID, u.ID)
		c.Next()
	}
}

// RequireRole 角色权限校验
func RequireRole(roles ...string) gin.HandlerFunc {
	allowed := make(map[string]bool, len(roles))
	for _, r := range roles {
		allowed[r] = true
	}
	return func(c *gin.Context) {
		u, ok := c.Get(CtxUser)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error_code":    "UNAUTHORIZED",
				"error_message": "Not authenticated",
			})
			return
		}
		user := u.(*models.User)
		if !allowed[user.Role] {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error_code":    "FORBIDDEN",
				"error_message": "Insufficient permissions",
				"detail":        "This action requires one of roles: " + strings.Join(roles, ", "),
			})
			return
		}
		c.Next()
	}
}

// CurrentUser 从上下文取当前用户
func CurrentUser(c *gin.Context) *models.User {
	v, ok := c.Get(CtxUser)
	if !ok {
		return nil
	}
	return v.(*models.User)
}
