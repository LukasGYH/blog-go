package constant

// UserStatus 用户状态
type UserStatus string

const (
	UserStatusActive   UserStatus = "active"   // 正常
	UserStatusDisabled UserStatus = "disabled" // 已禁用
)
