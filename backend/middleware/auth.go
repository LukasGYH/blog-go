package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"blog-api/constant"
	"blog-api/utils"
)

// 上下文键名
const (
	ctxUserID   = "userID"
	ctxUsername = "username"
	ctxUserRole = "userRole"
)

// errMissingToken 缺失或格式错误的 Authorization 头
var errMissingToken = errors.New("missing or malformed authorization header")

// Auth 强制鉴权中间件：解析 Bearer token，将 userID / username / role 写入上下文
func Auth(secret []byte) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, err := extractClaims(c, secret)
		if err != nil {
			utils.Error(c, http.StatusUnauthorized, "unauthorized")
			c.Abort()
			return
		}
		setCurrentUser(c, claims)
		c.Next()
	}
}

// OptionalAuth 可选鉴权中间件：携带有效 token 时写入用户上下文，否则放行（游客）
func OptionalAuth(secret []byte) gin.HandlerFunc {
	return func(c *gin.Context) {
		if claims, err := extractClaims(c, secret); err == nil {
			setCurrentUser(c, claims)
		}
		c.Next()
	}
}

// AdminOnly 管理员校验中间件，需置于 Auth / OptionalAuth 之后
func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		if CurrentUserRole(c) != constant.RoleAdmin {
			utils.Error(c, http.StatusForbidden, "admin only")
			c.Abort()
			return
		}
		c.Next()
	}
}

// CurrentUserID 从上下文读取当前用户 ID，未登录返回 0
func CurrentUserID(c *gin.Context) uint {
	if value, ok := c.Get(ctxUserID); ok {
		if id, ok := value.(uint); ok {
			return id
		}
	}
	return 0
}

// CurrentUserRole 从上下文读取当前用户角色，未登录返回空字符串
func CurrentUserRole(c *gin.Context) constant.Role {
	if value, ok := c.Get(ctxUserRole); ok {
		if role, ok := value.(constant.Role); ok {
			return role
		}
	}
	return ""
}

// extractClaims 从 Authorization 头解析并校验 token
func extractClaims(c *gin.Context, secret []byte) (*utils.Claims, error) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return nil, errMissingToken
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return nil, errMissingToken
	}
	return utils.ParseToken(parts[1], secret)
}

// setCurrentUser 将用户身份写入 gin 上下文
func setCurrentUser(c *gin.Context, claims *utils.Claims) {
	c.Set(ctxUserID, claims.UserID)
	c.Set(ctxUsername, claims.Username)
	c.Set(ctxUserRole, claims.Role)
}
