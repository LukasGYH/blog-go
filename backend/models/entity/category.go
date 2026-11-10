package entity

import (
	"time"

	"gorm.io/gorm"
)

// Category 文章分类（对应表 categories）
type Category struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"uniqueIndex;not null;size:50"`
	Slug        string         `json:"slug" gorm:"uniqueIndex;not null;size:50"`
	Description string         `json:"description" gorm:"size:255"`
	SortOrder   int            `json:"sort_order" gorm:"default:0"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// Articles 用于统计分类下文章数（Preload 关联，不参与 JSON 输出）
	Articles []Article `json:"-" gorm:"foreignKey:CategoryID"`
}
