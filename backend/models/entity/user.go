package entity

import (
	"time"

	"gorm.io/gorm"

	"blog-api/constant"
)

// User 用户实体（对应表 users）
type User struct {
	ID        uint                `json:"id" gorm:"primaryKey"`
	Username  string              `json:"username" gorm:"uniqueIndex;not null;size:50"`
	Email     string              `json:"email" gorm:"uniqueIndex;not null;size:100"`
	Password  string              `json:"-" gorm:"not null;size:100"` // bcrypt 密文，JSON 永不返回
	Nickname  string              `json:"nickname" gorm:"size:50"`
	Avatar    string              `json:"avatar" gorm:"size:255"`
	Bio       string              `json:"bio" gorm:"size:255"`
	Role      constant.Role       `json:"role" gorm:"index;not null;default:user;size:16"`
	Status    constant.UserStatus `json:"status" gorm:"not null;default:active;size:16"`
	CreatedAt time.Time           `json:"created_at"`
	UpdatedAt time.Time           `json:"updated_at"`
	DeletedAt gorm.DeletedAt      `json:"-" gorm:"index"`
}
