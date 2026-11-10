package services

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"blog-api/models/dto"
	"blog-api/models/entity"
	"blog-api/models/vo"
	"blog-api/repositories"
	"blog-api/utils"
)

// CategoryService 分类业务逻辑（标准 CRUD + 删除前关联校验参考样板）
type CategoryService struct {
	repo *repositories.CategoryRepository
}

// NewCategoryService 创建分类服务
func NewCategoryService(repo *repositories.CategoryRepository) *CategoryService {
	return &CategoryService{repo: repo}
}

// List 分类列表（含 article_count）
func (s *CategoryService) List(ctx context.Context) ([]vo.CategoryVO, error) {
	categories, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]vo.CategoryVO, 0, len(categories))
	for _, category := range categories {
		result = append(result, vo.CategoryVO{
			ID:           category.ID,
			Name:         category.Name,
			Slug:         category.Slug,
			Description:  category.Description,
			SortOrder:    category.SortOrder,
			ArticleCount: int64(len(category.Articles)),
		})
	}
	return result, nil
}

// Create 新增分类：名称唯一，slug 未传则由名称生成
func (s *CategoryService) Create(ctx context.Context, req dto.CreateCategoryRequest) (*vo.CategoryVO, error) {
	existing, err := s.repo.GetByName(ctx, req.Name)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, utils.Conflict("category name already exists")
	}

	slug := req.Slug
	if slug == "" {
		slug = utils.GenerateSlug(req.Name)
	}
	if err := s.ensureSlugAvailable(ctx, slug, 0); err != nil {
		return nil, err
	}

	category := &entity.Category{
		Name:        req.Name,
		Slug:        slug,
		Description: req.Description,
		SortOrder:   req.SortOrder,
	}
	if err := s.repo.Create(ctx, category); err != nil {
		// 并发下仍可能撞 name / slug 唯一索引，兜底识别为 409
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, utils.Conflict("category name already exists")
		}
		return nil, err
	}

	return &vo.CategoryVO{
		ID:           category.ID,
		Name:         category.Name,
		Slug:         category.Slug,
		Description:  category.Description,
		SortOrder:    category.SortOrder,
		ArticleCount: 0,
	}, nil
}

// Update 修改分类：名称 / slug 改后仍需唯一
func (s *CategoryService) Update(ctx context.Context, id uint, req dto.UpdateCategoryRequest) (*vo.CategoryVO, error) {
	category, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if category == nil {
		return nil, utils.NotFound("category not found")
	}

	if req.Name != category.Name {
		existing, err := s.repo.GetByName(ctx, req.Name)
		if err != nil {
			return nil, err
		}
		if existing != nil && existing.ID != category.ID {
			return nil, utils.Conflict("category name already exists")
		}
	}

	slug := req.Slug
	if slug == "" {
		slug = utils.GenerateSlug(req.Name)
	}
	if err := s.ensureSlugAvailable(ctx, slug, category.ID); err != nil {
		return nil, err
	}

	category.Name = req.Name
	category.Slug = slug
	category.Description = req.Description
	category.SortOrder = req.SortOrder
	if err := s.repo.Update(ctx, category); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, utils.Conflict("category name already exists")
		}
		return nil, err
	}

	articleCount, err := s.repo.CountArticles(ctx, category.ID)
	if err != nil {
		return nil, err
	}
	return &vo.CategoryVO{
		ID:           category.ID,
		Name:         category.Name,
		Slug:         category.Slug,
		Description:  category.Description,
		SortOrder:    category.SortOrder,
		ArticleCount: articleCount,
	}, nil
}

// Delete 删除分类：该分类下仍有文章时禁止删除（409）
func (s *CategoryService) Delete(ctx context.Context, id uint) error {
	category, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if category == nil {
		return utils.NotFound("category not found")
	}

	articleCount, err := s.repo.CountArticles(ctx, id)
	if err != nil {
		return err
	}
	if articleCount > 0 {
		return utils.Conflict("category has articles, cannot delete")
	}

	return s.repo.Delete(ctx, id)
}

// ensureSlugAvailable 校验 slug 未被其它分类占用（excludeID 为当前分类，创建时传 0）
func (s *CategoryService) ensureSlugAvailable(ctx context.Context, slug string, excludeID uint) error {
	existing, err := s.repo.GetBySlug(ctx, slug)
	if err != nil {
		return err
	}
	if existing != nil && existing.ID != excludeID {
		return utils.Conflict("category slug already exists")
	}
	return nil
}
