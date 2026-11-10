package repositories

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"blog-api/models/entity"
)

// CategoryRepository 分类数据访问（标准 CRUD + 聚合参考样板）
type CategoryRepository struct {
	db *gorm.DB
}

// NewCategoryRepository 创建分类仓储
func NewCategoryRepository(db *gorm.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

// List 返回全部分类，Preload 关联文章并裁剪字段以统计 article_count（2 次查询，无 N+1）
func (r *CategoryRepository) List(ctx context.Context) ([]entity.Category, error) {
	categories := make([]entity.Category, 0)
	err := r.db.WithContext(ctx).
		Preload("Articles", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "category_id")
		}).
		Order("sort_order asc, id asc").
		Find(&categories).Error
	if err != nil {
		return nil, err
	}
	return categories, nil
}

// GetByID 按主键查询分类，不存在返回 (nil, nil)
func (r *CategoryRepository) GetByID(ctx context.Context, id uint) (*entity.Category, error) {
	var category entity.Category
	if err := r.db.WithContext(ctx).First(&category, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &category, nil
}

// GetByName 按名称查询分类，不存在返回 (nil, nil)
func (r *CategoryRepository) GetByName(ctx context.Context, name string) (*entity.Category, error) {
	var category entity.Category
	if err := r.db.WithContext(ctx).Where("name = ?", name).First(&category).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &category, nil
}

// GetBySlug 按 slug 查询分类，不存在返回 (nil, nil)
func (r *CategoryRepository) GetBySlug(ctx context.Context, slug string) (*entity.Category, error) {
	var category entity.Category
	if err := r.db.WithContext(ctx).Where("slug = ?", slug).First(&category).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &category, nil
}

// Create 创建分类
func (r *CategoryRepository) Create(ctx context.Context, category *entity.Category) error {
	return r.db.WithContext(ctx).Create(category).Error
}

// Update 保存分类变更
func (r *CategoryRepository) Update(ctx context.Context, category *entity.Category) error {
	return r.db.WithContext(ctx).Save(category).Error
}

// Delete 硬删除分类
//   - 分类属「引用数据」，与 Tag 保持一致改为硬删除：软删除会留下行并长期占用
//     name / slug 的唯一索引，导致同名分类删除后无法再次创建（INSERT 撞唯一索引）
func (r *CategoryRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Unscoped().Delete(&entity.Category{}, id).Error
}

// CountArticles 统计分类下的文章数（已排除软删除）
func (r *CategoryRepository) CountArticles(ctx context.Context, categoryID uint) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&entity.Article{}).
		Where("category_id = ?", categoryID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
