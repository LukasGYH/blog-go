package entity

import (
	"time"

	"gorm.io/gorm"
)

// Tag 文章标签（对应表 tags，与文章为多对多）
type Tag struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Name      string         `json:"name" gorm:"uniqueIndex;not null;size:50"`
	Slug      string         `json:"slug" gorm:"uniqueIndex;not null;size:50"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// Articles 用于统计标签下文章数（多对多 article_tags 关联，不参与 JSON 输出）
	Articles []Article `json:"-" gorm:"many2many:article_tags;"`
}
