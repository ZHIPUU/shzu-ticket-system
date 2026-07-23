package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"ticket-system/config"
	"ticket-system/handlers"
	"ticket-system/middleware"
	"ticket-system/models"
)

func main() {
	// 子命令：reset-admin <new-password> — 一次性重置 admin 密码
	resetCmd := flag.NewFlagSet("reset-admin", flag.ExitOnError)
	if len(os.Args) >= 2 && os.Args[1] == "reset-admin" {
		_ = resetCmd.Parse(os.Args[2:])
		if resetCmd.NArg() < 1 {
			log.Fatal("用法: ticket-server reset-admin <new-password>")
		}
		newPwd := resetCmd.Arg(0)
		cfg := config.Load()
		db, err := gorm.Open(sqlite.Open(cfg.DatabaseURL), &gorm.Config{})
		if err != nil {
			log.Fatalf("打开 db 失败: %v", err)
		}
		// 兼容空 db
		if err := db.AutoMigrate(&models.Ticket{}, &models.User{}); err != nil {
			log.Fatalf("迁移失败: %v", err)
		}
		// 如果 admin 不存在，先创建
		var admin models.User
		if err := db.Where("username = ?", cfg.AdminUser).First(&admin).Error; err == gorm.ErrRecordNotFound {
			hash, _ := bcrypt.GenerateFromPassword([]byte(newPwd), bcrypt.DefaultCost)
			newAdmin := models.User{
				Username:      cfg.AdminUser,
				PasswordHash:  string(hash),
				Role:          "admin",
				DisplayName:   "系统管理员",
				Active:        true,
				MustChangePwd: true,
			}
			if err := db.Create(&newAdmin).Error; err != nil {
				log.Fatalf("创建 admin 失败: %v", err)
			}
			log.Printf("✅ admin 已创建: username=%s password=%s", cfg.AdminUser, newPwd)
			return
		} else if err != nil {
			log.Fatalf("查询失败: %v", err)
		}
		// admin 已存在：重置密码
		hash, _ := bcrypt.GenerateFromPassword([]byte(newPwd), bcrypt.DefaultCost)
		res := db.Model(&models.User{}).Where("username = ?", cfg.AdminUser).
			Updates(map[string]interface{}{
				"password_hash":   string(hash),
				"must_change_pwd": true,
				"active":          true,
			})
		if res.Error != nil {
			log.Fatalf("更新失败: %v", res.Error)
		}
		log.Printf("✅ admin 密码已重置为: %s  (RowsAffected=%d)", newPwd, res.RowsAffected)
		return
	}

	cfg := config.Load()

	// 初始化数据库
	db, err := gorm.Open(sqlite.Open(cfg.DatabaseURL), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	if err := db.AutoMigrate(&models.Ticket{}, &models.User{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
	log.Printf("Database ready: %s", cfg.DatabaseURL)

	// 初始化默认 admin（首次启动）
	ensureDefaultAdmin(db, cfg)

	// 初始化 handler
	ticketH := handlers.New(db, cfg)
	authH := handlers.NewAuth(db, cfg)
	userH := handlers.NewUser(db)

	// 路由
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "X-API-Key", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// 健康检查（无需鉴权）
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "service": "ticket-system"})
	})

	// ─── Auth 公开端点（无需 token）───
	api := r.Group("/api/v1")
	api.POST("/auth/login", authH.Login)

	// ─── 鉴权辅助中间件：双轨（API Key 或 JWT）───
	// 智能体走 API Key，后台用户走 JWT，二选一
	dualAuth := func(c *gin.Context) {
		// 优先检查 JWT
		auth := c.GetHeader("Authorization")
		if len(auth) > 7 && auth[:7] == "Bearer " {
			log.Printf("[dualAuth] JWT path, auth head len=%d", len(auth))
			middleware.JWTAuth(cfg, db)(c)
			return
		}
		log.Printf("[dualAuth] APIKey path, auth head=[%s]", auth)
		// fallback 到 API Key
		middleware.APIKeyAuth(cfg.APIKey)(c)
	}

	// ─── 工单 API（双轨：API Key 给智能体 / JWT 给后台）───
	tickets := api.Group("/tickets")
	tickets.Use(dualAuth)
	{
		tickets.POST("", ticketH.SubmitTicket)
		tickets.GET("", ticketH.ListTickets)
		tickets.GET("/export", ticketH.ExportTickets)
		tickets.POST("/batch-delete", ticketH.BatchDelete)
		tickets.GET("/:ticket_id", ticketH.GetTicket)
		tickets.PATCH("/:ticket_id", ticketH.PatchTicket)
		tickets.DELETE("/:ticket_id", ticketH.DeleteTicket)
		tickets.POST("/:ticket_id/answer", ticketH.AnswerTicket)
		tickets.POST("/:ticket_id/close", ticketH.CloseTicket)
		tickets.POST("/:ticket_id/reopen", ticketH.ReopenTicket)
		tickets.POST("/:ticket_id/archive", ticketH.ArchiveTicket)
		tickets.POST("/:ticket_id/unarchive", ticketH.UnarchiveTicket)
	}

	// ─── Auth 受保护端点（仅 JWT）───
	auth := api.Group("/auth")
	auth.Use(middleware.JWTAuth(cfg, db))
	{
		auth.GET("/me", authH.Me)
		auth.POST("/change-password", authH.ChangePassword)
	}

	// ─── User 管理（仅 JWT + admin）───
	users := api.Group("/users")
	users.Use(middleware.JWTAuth(cfg, db))
	users.Use(middleware.RequireRole("admin"))
	{
		users.GET("", userH.ListUsers)
		users.POST("", userH.CreateUser)
		users.PATCH("/:id", userH.UpdateUser)
	}

	addr := cfg.Host + ":" + cfg.Port
	log.Printf("Starting server on %s (API key: %s...)", addr, maskKey(cfg.APIKey))
	log.Printf("Default admin: %s / (set ADMIN_PASSWORD env to override)", cfg.AdminUser)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

// ensureDefaultAdmin 首次启动时插入默认 admin；已存在则不重置密码（保护用户已修改的密码）
func ensureDefaultAdmin(db *gorm.DB, cfg *config.Config) {
	var admin models.User
	err := db.Where("username = ?", cfg.AdminUser).First(&admin).Error
	if err == gorm.ErrRecordNotFound {
		// 首次启动：创建默认 admin
		hash, hErr := bcrypt.GenerateFromPassword([]byte(cfg.AdminPass), bcrypt.DefaultCost)
		if hErr != nil {
			log.Fatalf("Failed to hash admin password: %v", hErr)
		}
		newAdmin := models.User{
			Username:      cfg.AdminUser,
			PasswordHash:  string(hash),
			Role:          "admin",
			DisplayName:   "系统管理员",
			Active:        true,
			MustChangePwd: true,
		}
		if err := db.Create(&newAdmin).Error; err != nil {
			log.Fatalf("Failed to create default admin: %v", err)
		}
		log.Printf("✨ Default admin created: username=%s password=%s (MUST change on first login!)",
			cfg.AdminUser, cfg.AdminPass)
		return
	}
	// admin 已存在：不再强制重置密码，避免覆盖用户已修改的密码
	// 若需要重置，请通过管理后台 /api/v1/users/:id PATCH 显式操作
	if !admin.Active {
		admin.Active = true
		_ = db.Save(&admin)
	}
}

func maskKey(k string) string {
	if len(k) <= 8 {
		return k
	}
	return k[:4] + "***" + k[len(k)-4:]
}
