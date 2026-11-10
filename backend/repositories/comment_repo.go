package repositories

import (
	"context"

	"gorm.io/gorm"

	"blog-api/models/entity"
	"blog-api/utils"
)

// CommentRepository 评论数据访问
type CommentRepository struct {
	db *gorm.DB
}

// NewCommentRepository 创建评论仓储
func NewCommentRepository(db *gorm.DB) *CommentRepository {
	return &CommentRepository{db: db}
}

// ListTopLevel 分页查询顶层评论（parent_id = 0），Preload 评论人 User
//   - Count 与列表复用同一套过滤条件；分页只作用于列表，避免污染 Count
func (r *CommentRepository) ListTopLevel(ctx context.Context, articleID uint, page utils.PageQuery) ([]entity.Comment, int64, error) {
	base := func() *gorm.DB {
		return r.db.WithContext(ctx).
			Model(&entity.Comment{}).
			Where("article_id = ? and parent_id = 0", articleID)
	}

	var count int64
	if err := base().Count(&count).Error; err != nil {
		return nil, 0, err
	}

	comments := make([]entity.Comment, 0)
	if err := base().
		Preload("User").
		Order("created_at asc, id asc").
		Offset(page.Offset()).
		Limit(page.PageSize).
		Find(&comments).Error; err != nil {
		return nil, 0, err
	}
	return comments, count, nil
}

// ListReplies 按 parentIDs 一次性查询全部回复（Preload User），为树形组装准备数据
func (r *CommentRepository) ListReplies(ctx context.Context, parentIDs []uint) ([]entity.Comment, error) {
	replies := make([]entity.Comment, 0)
	if len(parentIDs) == 0 {
		// 空 IN 会生成 `IN (NULL)`，直接短路返回
		return replies, nil
	}
	if err := r.db.WithContext(ctx).
		Preload("User").
		Where("parent_id IN ?", parentIDs).
		Order("created_at asc, id asc").
		Find(&replies).Error; err != nil {
		return nil, err
	}
	return replies, nil
}

// GetByID 按主键查询评论
//   - 不存在时返回 gorm.ErrRecordNotFound，由 Service 层翻译为 404
func (r *CommentRepository) GetByID(ctx context.Context, id uint) (*entity.Comment, error) {
	var comment entity.Comment
	if err := r.db.WithContext(ctx).
		Preload("User").
		First(&comment, id).Error; err != nil {
		return nil, err
	}
	return &comment, nil
}

// Create 创建评论
func (r *CommentRepository) Create(ctx context.Context, comment *entity.Comment) error {
	return r.db.WithContext(ctx).Create(comment).Error
}

// Delete 软删除评论
func (r *CommentRepository) Delete(ctx context.Context, id uint) error {
	return r.deleteWithDB(r.db.WithContext(ctx), id)
}

// DeleteReplies 软删除某顶层评论下的全部回复
func (r *CommentRepository) DeleteReplies(ctx context.Context, parentID uint) error {
	return r.deleteRepliesWithDB(r.db.WithContext(ctx), parentID)
}

// DeleteWithReplies 在一个事务内软删除评论及其全部子回复
//   - 关键点：事务内所有操作必须基于闭包传入的 tx，用 r.db 发的 SQL 会跑到事务外
func (r *CommentRepository) DeleteWithReplies(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := r.deleteWithDB(tx, id); err != nil {
			return err
		}
		if err := r.deleteRepliesWithDB(tx, id); err != nil {
			return err
		}
		return nil
	})
}

// deleteWithDB 接收 *gorm.DB，便于复用同一事务句柄
func (r *CommentRepository) deleteWithDB(db *gorm.DB, id uint) error {
	return db.Delete(&entity.Comment{}, id).Error
}

// deleteRepliesWithDB 接收 *gorm.DB，便于复用同一事务句柄
func (r *CommentRepository) deleteRepliesWithDB(db *gorm.DB, parentID uint) error {
	return db.Where("parent_id = ?", parentID).Delete(&entity.Comment{}).Error
}
