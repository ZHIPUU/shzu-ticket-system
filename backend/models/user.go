package models

import "time"

// User 用户表
type User struct {
	ID           uint       `gorm:"primaryKey" json:"id"`
	Username     string     `gorm:"uniqueIndex;size:64;not null" json:"username"`
	PasswordHash string     `gorm:"size:128;not null" json:"-"`
	Role         string     `gorm:"size:16;not null;default:staff" json:"role"`
	// admin: 管理员（全权限，含用户管理）
	// staff: 工作人员（仅工单操作）
	DisplayName  string     `gorm:"size:64" json:"display_name"`
	Email        string     `gorm:"size:128" json:"email"`
	Active       bool       `gorm:"default:true;index" json:"active"`
	MustChangePwd bool      `gorm:"default:false" json:"must_change_password"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	LastLoginAt  *time.Time `json:"last_login_at,omitempty"`
}
