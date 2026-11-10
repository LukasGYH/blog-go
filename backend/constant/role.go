package constant

// Role 用户角色
type Role string

const (
	RoleUser  Role = "user"  // 普通注册用户
	RoleAdmin Role = "admin" // 管理员
)
