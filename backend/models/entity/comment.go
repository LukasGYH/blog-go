package entity

import (
	"time"

	"gorm.io/gorm"

	"blog-api/constant"
)

// Comment 评论实体（对应表 comments，两级：顶层评论 + 回复）
type Comment struct {
	ID        uint                  `json:"id" gorm:"primaryKey"`
	ArticleID uint                  `json:"article_id" gorm:"index;not null"`
	UserID    uint                  `json:"user_id" gorm:"index;not null"`
	User      User                  `json:"user" gorm:"foreignKey:UserID"`
	ParentID  uint                  `json:"parent_id" gorm:"index;default:0"` // 0 表示顶层评论
	Content   string                `json:"content" gorm:"not null;size:500"`
	Status    constant.CommentStatus `json:"status" gorm:"not null;default:visible;size:16"`
	CreatedAt time.Time             `json:"created_at"`
	DeletedAt gorm.DeletedAt        `json:"-" gorm:"index"`

	// Replies 子回复（自关联，仅两级，不参与 JSON 输出）
	Replies []Comment `json:"-" gorm:"foreignKey:ParentID"`
}
