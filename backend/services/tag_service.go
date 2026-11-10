package services

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"blog-api/models/dto"
	"blog-api/models/entity"
	"blog-api/models/vo"
	"blog-api/repositories"
	"blog-api/utils"
)

// TagService 标签业务逻辑
type TagService struct {
	repo *repositories.TagRepository
}

// NewTagService 创建标签服务
func NewTagService(repo *repositories.TagRepository) *TagService {
	return &TagService{repo: repo}
}

// List 标签列表（含文章数，不分页）
func (s *TagService) List(ctx context.Context, keyword string) ([]vo.TagVO, error) {
	list, err := s.repo.List(ctx, keyword)
	if err != nil {
		return nil, fmt.Errorf("查询标签列表失败: %w", err)
	}
	return list, nil
}

// Create 新增标签：名称唯一，slug 由名称生成并唯一化
func (s *TagService) Create(ctx context.Context, req dto.CreateTagRequest) (*vo.TagVO, error) {
	// 名称唯一性预检（并发下仍可能撞唯一索引，故 Create 处再兜底识别 gorm.ErrDuplicatedKey）
	if err := s.ensureNameAvailable(ctx, req.Name, 0); err != nil {
		return nil, err
	}

	var slugErr error
	slug := utils.EnsureUnique(utils.GenerateSlug(req.Name), func(candidate string) bool {
		existing, err := s.repo.GetBySlug(ctx, candidate)
		if err != nil {
			slugErr = err
			// 出错时返回 false 结束循环，随后统一上报错误
			return false
		}
		return existing != nil
	})
	if slugErr != nil {
		return nil, fmt.Errorf("校验 slug 唯一性失败: %w", slugErr)
	}

	tag := &entity.Tag{Name: req.Name, Slug: slug}
	if err := s.repo.Create(ctx, tag); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, utils.Conflict("标签名称已存在")
		}
		return nil, fmt.Errorf("创建标签失败: %w", err)
	}
	return &vo.TagVO{ID: tag.ID, Name: tag.Name, Slug: tag.Slug}, nil
}

// Update 修改标签：名称唯一性预检（排除自身）
func (s *TagService) Update(ctx context.Context, id uint, req dto.UpdateTagRequest) (*vo.TagVO, error) {
	tag, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.NotFound("标签不存在")
		}
		return nil, fmt.Errorf("查询标签失败: %w", err)
	}

	if tag.Name != req.Name {
		if err := s.ensureNameAvailable(ctx, req.Name, id); err != nil {
			return nil, err
		}
		tag.Name = req.Name
	}

	if err := s.repo.Update(ctx, tag); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, utils.Conflict("标签名称已存在")
		}
		return nil, fmt.Errorf("更新标签失败: %w", err)
	}
	return &vo.TagVO{
		ID:           tag.ID,
		Name:         tag.Name,
		Slug:         tag.Slug,
		ArticleCount: int64(len(tag.Articles)),
	}, nil
}

// Delete 删除标签（硬删除并清理中间表关联）
func (s *TagService) Delete(ctx context.Context, id uint) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("删除标签失败: %w", err)
	}
	return nil
}

// ensureNameAvailable 校验标签名称是否已被其他标签占用
//   - excludeID 为当前正在编辑的标签 ID，创建场景传 0
func (s *TagService) ensureNameAvailable(ctx context.Context, name string, excludeID uint) error {
	existing, err := s.repo.GetByName(ctx, name)
	if err != nil {
		return fmt.Errorf("校验标签名称失败: %w", err)
	}
	if existing != nil && existing.ID != excludeID {
		return utils.Conflict("标签名称已存在")
	}
	return nil
}
