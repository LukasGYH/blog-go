package handlers

import (
	"time"

	"github.com/gin-gonic/gin"

	"blog-api/config"
	"blog-api/middleware"
	"blog-api/models/dto"
	"blog-api/models/vo"
	"blog-api/services"
	"blog-api/utils"
)

// AuthHandler 认证相关接口
type AuthHandler struct {
	authService *services.AuthService
	jwtSecret   []byte
	jwtExpire   time.Duration
}

// NewAuthHandler 创建认证处理器
func NewAuthHandler(authService *services.AuthService, cfg *config.Config) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		jwtSecret:   []byte(cfg.JWTSecret),
		jwtExpire:   cfg.JWTExpire,
	}
}

// Register 注册
func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, utils.ParseBindingError(err))
		return
	}

	result, err := h.authService.Register(c.Request.Context(), req)
	if err != nil {
		utils.HandleError(c, err)
		return
	}
	utils.Created(c, result)
}

// Login 登录
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, utils.ParseBindingError(err))
		return
	}

	user, err := h.authService.Login(c.Request.Context(), req)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	token, expiresAt, err := utils.GenerateToken(h.jwtSecret, user.ID, user.Username, user.Role, h.jwtExpire)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	utils.Success(c, gin.H{
		"token":      token,
		"expires_at": expiresAt,
		"user": vo.UserVO{
			ID:       user.ID,
			Username: user.Username,
			Nickname: user.Nickname,
			Avatar:   user.Avatar,
			Role:     user.Role,
		},
	})
}

// Me 查看当前用户资料
func (h *AuthHandler) Me(c *gin.Context) {
	profile, err := h.authService.GetProfile(c.Request.Context(), middleware.CurrentUserID(c))
	if err != nil {
		utils.HandleError(c, err)
		return
	}
	utils.Success(c, profile)
}

// UpdateMe 修改当前用户资料
func (h *AuthHandler) UpdateMe(c *gin.Context) {
	var req dto.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, utils.ParseBindingError(err))
		return
	}

	profile, err := h.authService.UpdateProfile(c.Request.Context(), middleware.CurrentUserID(c), req)
	if err != nil {
		utils.HandleError(c, err)
		return
	}
	utils.Success(c, profile)
}

// ChangePassword 修改密码
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	var req dto.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, utils.ParseBindingError(err))
		return
	}

	if err := h.authService.ChangePassword(c.Request.Context(), middleware.CurrentUserID(c), req); err != nil {
		utils.HandleError(c, err)
		return
	}
	utils.Success(c, gin.H{"message": "password updated"})
}
