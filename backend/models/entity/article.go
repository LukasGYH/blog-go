package entity

import (
	"time"

	"gorm.io/gorm"

	"blog-api/constant"
)

// Article 文章实体（对应表 articles）
type Article struct {
	ID          uint                   `json:"id" gorm:"primaryKey"`
	Title       string                 `json:"title" gorm:"not null;size:200"`
	Slug        string                 `json:"slug" gorm:"uniqueIndex;not null;size:200"`
	Summary     string                 `json:"summary" gorm:"size:500"`
	Content     string                 `json:"content" gorm:"not null;type:text"`
	CoverImage  string                 `json:"cover_image" gorm:"size:255"`
	Status      constant.ArticleStatus `json:"status" gorm:"index;index:idx_articles_status_published,priority:1;not null;default:draft;size:16"`
	ViewCount   int                    `json:"view_count" gorm:"not null;default:0"`
	AuthorID    uint                   `json:"author_id" gorm:"index;not null"`
	Author      User                   `json:"author" gorm:"foreignKey:AuthorID"`
	CategoryID  *uint                  `json:"category_id" gorm:"index"`
	Category    *Category              `json:"category" gorm:"foreignKey:CategoryID"`
	Tags        []Tag                  `json:"tags" gorm:"many2many:article_tags;"`
	PublishedAt *time.Time             `json:"published_at" gorm:"index;index:idx_articles_status_published,priority:2"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
	DeletedAt   gorm.DeletedAt         `json:"-" gorm:"index"`

	// CommentCount 评论数（非数据库字段，由仓储层按需聚合填充）
	CommentCount int64 `json:"-" gorm:"-"`
}
