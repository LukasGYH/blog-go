package repositories

import (
	"context"

	"gorm.io/gorm"

	"blog-api/constant"
	"blog-api/models/entity"
	"blog-api/models/vo"
)

// StatsRepository 统计数据访问
type StatsRepository struct {
	db *gorm.DB
}

// NewStatsRepository 创建统计仓储
func NewStatsRepository(db *gorm.DB) *StatsRepository {
	return &StatsRepository{db: db}
}

// Overview 概览统计（文章总数 / 已发布 / 草稿 / 用户数 / 评论数 / 分类数 / 标签数）
func (r *StatsRepository) Overview(ctx context.Context) (*vo.StatsOverviewVO, error) {
	overview := &vo.StatsOverviewVO{}

	// 文章总数（含草稿）
	if err := r.db.WithContext(ctx).Model(&entity.Article{}).
		Count(&overview.ArticleCount).Error; err != nil {
		return nil, err
	}
	// 已发布
	if err := r.db.WithContext(ctx).Model(&entity.Article{}).
		Where("status = ?", constant.ArticleStatusPublished).
		Count(&overview.PublishedCount).Error; err != nil {
		return nil, err
	}
	// 草稿
	if err := r.db.WithContext(ctx).Model(&entity.Article{}).
		Where("status = ?", constant.ArticleStatusDraft).
		Count(&overview.DraftCount).Error; err != nil {
		return nil, err
	}
	// 用户数
	if err := r.db.WithContext(ctx).Model(&entity.User{}).
		Count(&overview.UserCount).Error; err != nil {
		return nil, err
	}
	// 评论数
	if err := r.db.WithContext(ctx).Model(&entity.Comment{}).
		Count(&overview.CommentCount).Error; err != nil {
		return nil, err
	}
	// 分类数
	if err := r.db.WithContext(ctx).Model(&entity.Category{}).
		Count(&overview.CategoryCount).Error; err != nil {
		return nil, err
	}
	// 标签数
	if err := r.db.WithContext(ctx).Model(&entity.Tag{}).
		Count(&overview.TagCount).Error; err != nil {
		return nil, err
	}
	return overview, nil
}

// CountByCategory 各分类文章数统计（LEFT JOIN 保证无文章的分类也返回 0）
func (r *StatsRepository) CountByCategory(ctx context.Context) ([]vo.CategoryStatsVO, error) {
	stats := make([]vo.CategoryStatsVO, 0)
	if err := r.db.WithContext(ctx).
		Table("categories c").
		Select("c.id as category_id, c.name as name, count(a.id) as article_count").
		Joins("LEFT JOIN articles a ON a.category_id = c.id AND a.deleted_at IS NULL").
		Where("c.deleted_at IS NULL").
		Group("c.id, c.name").
		Order("c.id asc").
		Scan(&stats).Error; err != nil {
		return nil, err
	}
	return stats, nil
}
