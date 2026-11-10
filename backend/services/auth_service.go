package services

import (
	"context"

	"blog-api/constant"
	"blog-api/models/dto"
	"blog-api/models/entity"
	"blog-api/models/vo"
	"blog-api/repositories"
	"blog-api/utils"
)

// AuthService 认证与用户资料业务逻辑（参考样板）
type AuthService struct {
	repo *repositories.UserRepository
}

// NewAuthService 创建认证服务
func NewAuthService(repo *repositories.UserRepository) *AuthService {
	return &AuthService{repo: repo}
}

// Register 注册：唯一性预检 → bcrypt 加密 → 建用户
func (s *AuthService) Register(ctx context.Context, req dto.RegisterRequest) (*vo.UserRegisterVO, error) {
	exists, err := s.repo.ExistsUsername(ctx, req.Username)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, utils.Conflict("username already exists")
	}

	exists, err = s.repo.ExistsEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, utils.Conflict("email already exists")
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	nickname := req.Nickname
	if nickname == "" {
		nickname = req.Username
	}

	user := &entity.User{
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
		Nickname: nickname,
		Role:     constant.RoleUser,
		Status:   constant.UserStatusActive,
	}
	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	return &vo.UserRegisterVO{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Nickname: user.Nickname,
		Role:     user.Role,
	}, nil
}

// Login 登录：用户不存在 / 密码错误统一返回 401（防用户枚举），被禁用返回 403
func (s *AuthService) Login(ctx context.Context, req dto.LoginRequest) (*entity.User, error) {
	user, err := s.repo.GetByUsername(ctx, req.Username)
	if err != nil {
		return nil, err
	}
	if user == nil || !utils.CheckPassword(user.Password, req.Password) {
		return nil, utils.Unauthorized("invalid username or password")
	}
	if user.Status == constant.UserStatusDisabled {
		return nil, utils.Forbidden("account has been disabled")
	}
	return user, nil
}

// GetProfile 查询当前登录用户资料
func (s *AuthService) GetProfile(ctx context.Context, userID uint) (*vo.UserProfileVO, error) {
	user, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, utils.NotFound("user not found")
	}
	return toUserProfileVO(user), nil
}

// UpdateProfile 修改本人资料（邮箱改后仍需唯一）
func (s *AuthService) UpdateProfile(ctx context.Context, userID uint, req dto.UpdateProfileRequest) (*vo.UserProfileVO, error) {
	user, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, utils.NotFound("user not found")
	}

	if req.Email != "" && req.Email != user.Email {
		exists, err := s.repo.ExistsEmail(ctx, req.Email)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, utils.Conflict("email already exists")
		}
		user.Email = req.Email
	}
	if req.Nickname != "" {
		user.Nickname = req.Nickname
	}
	if req.Avatar != "" {
		user.Avatar = req.Avatar
	}
	if req.Bio != "" {
		user.Bio = req.Bio
	}

	if err := s.repo.Update(ctx, user); err != nil {
		return nil, err
	}
	return toUserProfileVO(user), nil
}

// ChangePassword 修改密码：校验旧密码正确后写入新密文
func (s *AuthService) ChangePassword(ctx context.Context, userID uint, req dto.ChangePasswordRequest) error {
	user, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return err
	}
	if user == nil {
		return utils.NotFound("user not found")
	}
	if !utils.CheckPassword(user.Password, req.OldPassword) {
		return utils.BadRequest("old password is incorrect")
	}

	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		return err
	}
	user.Password = hashedPassword
	return s.repo.Update(ctx, user)
}

// toUserProfileVO 实体转当前用户资料 VO
func toUserProfileVO(user *entity.User) *vo.UserProfileVO {
	return &vo.UserProfileVO{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Nickname:  user.Nickname,
		Avatar:    user.Avatar,
		Bio:       user.Bio,
		Role:      user.Role,
		Status:    user.Status,
		CreatedAt: user.CreatedAt,
	}
}
