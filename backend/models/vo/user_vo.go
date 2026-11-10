package vo

import (
	"time"

	"blog-api/constant"
)

// UserVO 登录接口返回的用户精简信息
type UserVO struct {
	ID       uint          `json:"id"`
	Username string        `json:"username"`
	Nickname string        `json:"nickname"`
	Avatar   string        `json:"avatar"`
	Role     constant.Role `json:"role"`
}

// UserRegisterVO 注册接口返回的用户信息
type UserRegisterVO struct {
	ID       uint          `json:"id"`
	Username string        `json:"username"`
	Email    string        `json:"email"`
	Nickname string        `json:"nickname"`
	Role     constant.Role `json:"role"`
}

// UserProfileVO 当前登录用户资料
type UserProfileVO struct {
	ID        uint                `json:"id"`
	Username  string              `json:"username"`
	Email     string              `json:"email"`
	Nickname  string              `json:"nickname"`
	Avatar    string              `json:"avatar"`
	Bio       string              `json:"bio"`
	Role      constant.Role       `json:"role"`
	Status    constant.UserStatus `json:"status"`
	CreatedAt time.Time           `json:"created_at"`
}

// UserBriefVO 嵌套在文章 / 评论中的作者信息
type UserBriefVO struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
}
