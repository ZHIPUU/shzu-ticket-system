package models

import "time"

// Ticket 工单数据模型
type Ticket struct {
	ID          uint      `gorm:"primaryKey" json:"-"`
	TicketID    string    `gorm:"uniqueIndex;size:32;not null" json:"ticket_id"`
	Question    string    `gorm:"not null" json:"question"`
	UserID      string    `gorm:"index;size:128;not null" json:"user_id"`
	Source      string    `gorm:"size:32;default:hiagent_chat" json:"source"`
	RAGResult   string    `gorm:"default:''" json:"rag_result"`
	Status      string    `gorm:"size:16;default:pending;index" json:"status"`
	// pending: 待处理
	// processing: 处理中
	// answered: 已回答
	// closed: 已关闭
	CreatedAt   time.Time `json:"created_at"`
	AnsweredAt  *time.Time `json:"answered_at,omitempty"`
	Answer      *string    `json:"answer,omitempty"`
	AnsweredBy  *string    `gorm:"size:64" json:"answered_by,omitempty"`
	CloseReason *string    `json:"close_reason,omitempty"`

	// 归档与分类（v3 新增）
	Archived    bool       `gorm:"default:false;index" json:"archived"`
	ArchivedAt  *time.Time `json:"archived_at,omitempty"`
	Category    string     `gorm:"size:32;default:'';index" json:"category"` // 业务分类：教务/宿舍/招办/资助/...
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `gorm:"index" json:"deleted_at,omitempty"` // 软删除时间，nil 表示未删
}
