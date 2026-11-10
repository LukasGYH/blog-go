package repositories

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"blog-api/models/entity"
	"blog-api/models/vo"
)

// TagRepository 标签数据访问
type TagRepository struct {
	db *gorm.DB
}

// NewTagRepository 创建标签仓储
func NewTagRepository(db *gorm.DB) *TagRepository {
	return &TagRepository{db: db}
}

// List 查询标签列表（含各标签文章数；支持名称关键字模糊匹配）
//   - 前端标签云 / 标签多选 / 后台标签管理均按「一次性返回全部」消费，故不分页，与 /categories 保持一致
func (r *TagRepository) List(ctx context.Context, keyword string) ([]vo.TagVO, error) {
	query := r.db.WithContext(ctx).Model(&entity.Tag{})
	if keyword != "" {
		query = query.Where("name LIKE ?", "%"+keyword+"%")
	}

	// Preload 关联文章只为计数，只选 id 避免拖出 content 等大字段
	tags := make([]entity.Tag, 0)
	if err := query.
		Preload("Articles", func(db *gorm.DB) *gorm.DB {
			return db.Select("articles.id")
		}).
		Order("id asc").
		Find(&tags).Error; err != nil {
		return nil, err
	}

	// entity -> vo。文章数取关联切片长度，无 N+1。
	result := make([]vo.TagVO, 0, len(tags))
	for _, tag := range tags {
		result = append(result, vo.TagVO{
			ID:           tag.ID,
			Name:         tag.Name,
			Slug:         tag.Slug,
			ArticleCount: int64(len(tag.Articles)),
		})
	}
	return result, nil
}

// GetByID 按主键查询标签（含文章数）
//   - 不存在时返回 gorm.ErrRecordNotFound，由 Service 层翻译为 404
func (r *TagRepository) GetByID(ctx context.Context, id uint) (*entity.Tag, error) {
	var tag entity.Tag
	if err := r.db.WithContext(ctx).Preload("Articles").First(&tag, id).Error; err != nil {
		return nil, err
	}
	return &tag, nil
}

// GetByName 按名称查询标签，不存在返回 (nil, nil)
func (r *TagRepository) GetByName(ctx context.Context, name string) (*entity.Tag, error) {
	var tag entity.Tag
	if err := r.db.WithContext(ctx).Where("name = ?", name).First(&tag).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &tag, nil
}

// GetBySlug 按 slug 查询标签，不存在返回 (nil, nil)
func (r *TagRepository) GetBySlug(ctx context.Context, slug string) (*entity.Tag, error) {
	var tag entity.Tag
	if err := r.db.WithContext(ctx).Where("slug = ?", slug).First(&tag).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &tag, nil
}

// Create 创建标签（唯一约束冲突时返回 gorm.ErrDuplicatedKey，由 Service 翻译为 409）
func (r *TagRepository) Create(ctx context.Context, tag *entity.Tag) error {
	return r.db.WithContext(ctx).Create(tag).Error
}

// Update 保存标签变更（唯一约束冲突时返回 gorm.ErrDuplicatedKey）
//   - 只更新 name / slug 两列：用 Save 会连带 upsert 预加载出来的 Articles 关联，造成误写文章表
func (r *TagRepository) Update(ctx context.Context, tag *entity.Tag) error {
	return r.db.WithContext(ctx).
		Model(&entity.Tag{}).
		Where("id = ?", tag.ID).
		Select("name", "slug").
		Updates(tag).Error
}

// Delete 删除标签：事务内先清理 article_tags 中间表关联，再硬删除标签本身
func (r *TagRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 先清理中间表，避免遗留悬空关联（article_tags 无软删除，走原生 DELETE）
		if err := tx.Exec("DELETE FROM article_tags WHERE tag_id = ?", id).Error; err != nil {
			return err
		}
		return tx.Unscoped().Delete(&entity.Tag{}, id).Error
	})
}
