package repositories

import (
	"context"

	"gorm.io/gorm"

	"blog-api/constant"
	"blog-api/models/dto"
	"blog-api/models/entity"
	"blog-api/models/vo"
)

// ArticleRepository 文章数据访问
type ArticleRepository struct {
	db *gorm.DB
}

// NewArticleRepository 创建文章仓储
func NewArticleRepository(db *gorm.DB) *ArticleRepository {
	return &ArticleRepository{db: db}
}

// List 多条件分页查询文章卡片
//   - 关键字匹配 title / summary（LIKE）、分类、标签（JOIN article_tags）、作者、状态筛选
//   - 排序：latest 按 published_at 倒序；hot 按 view_count 倒序
//   - 分页 + Count 总数；Preload(Author/Category/Tags)
//   - 列表字段裁剪（不返回 content）并批量统计 comment_count
func (r *ArticleRepository) List(ctx context.Context, query dto.ArticleQuery) ([]entity.Article, int64, error) {
	// base() 只承载过滤条件，绝不包含分页；
	// 否则 Count 会连带 OFFSET/LIMIT，导致第 2 页起的 total 被算成 0（GORM 的坑）。
	base := func() *gorm.DB {
		tx := r.db.WithContext(ctx).Model(&entity.Article{})
		if query.Keyword != "" {
			tx = tx.Where("title LIKE ? OR summary LIKE ?", "%"+query.Keyword+"%", "%"+query.Keyword+"%")
		}
		if query.CategoryID != 0 {
			tx = tx.Where("category_id = ?", query.CategoryID)
		}
		if query.TagID != 0 {
			tx = tx.Joins("JOIN article_tags ON article_tags.article_id = articles.id").
				Where("article_tags.tag_id = ?", query.TagID)
		}
		if query.AuthorID != 0 {
			tx = tx.Where("author_id = ?", query.AuthorID)
		}
		if query.Status != "" {
			tx = tx.Where("status = ?", query.Status)
		}
		return tx
	}

	var count int64
	if err := base().Count(&count).Error; err != nil {
		return nil, 0, err
	}

	page := query.PageQuery()
	articles := make([]entity.Article, 0)
	// 字段裁剪：列表只需要卡片字段，不把 content 这类大字段拖出来。
	tx := base().
		Select("articles.id", "articles.title", "articles.slug", "articles.summary",
			"articles.cover_image", "articles.status", "articles.view_count",
			"articles.author_id", "articles.category_id", "articles.published_at", "articles.created_at").
		Preload("Author").Preload("Category").Preload("Tags").
		Offset(page.Offset()).Limit(page.PageSize)
	if query.Sort == constant.SortHot {
		tx = tx.Order("view_count DESC")
	} else {
		tx = tx.Order("published_at DESC")
	}
	if err := tx.Find(&articles).Error; err != nil {
		return nil, 0, err
	}

	// 一次 IN 查询批量统计评论数，避免逐条查询造成 N+1。
	counts, err := r.countComments(ctx, articleIDs(articles))
	if err != nil {
		return nil, 0, err
	}
	for i := range articles {
		articles[i].CommentCount = counts[articles[i].ID]
	}
	return articles, count, nil
}

// GetBySlug 按 slug 查询文章
//   - 不存在时返回 gorm.ErrRecordNotFound，由 Service 层翻译为 404
func (r *ArticleRepository) GetBySlug(ctx context.Context, slug string) (*entity.Article, error) {
	var article entity.Article
	if err := r.db.WithContext(ctx).
		Preload("Author").Preload("Category").Preload("Tags").
		Where("slug = ?", slug).
		First(&article).Error; err != nil {
		return nil, err
	}
	return &article, nil
}

// GetByID 按主键查询文章（预加载 Author / Category / Tags 并附带评论数）
//   - 不存在时返回 gorm.ErrRecordNotFound，由 Service 层翻译为 404
func (r *ArticleRepository) GetByID(ctx context.Context, id uint) (*entity.Article, error) {
	var article entity.Article
	if err := r.db.WithContext(ctx).
		Preload("Author").Preload("Category").Preload("Tags").
		First(&article, id).Error; err != nil {
		return nil, err
	}
	counts, err := r.countComments(ctx, []uint{article.ID})
	if err != nil {
		return nil, err
	}
	article.CommentCount = counts[article.ID]
	return &article, nil
}

// SlugExists 判断 slug 是否已被占用
//   - 使用 Unscoped：软删除的记录仍占用唯一索引，必须一并纳入判断，否则重新创建同标题会撞唯一约束
func (r *ArticleRepository) SlugExists(ctx context.Context, slug string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Unscoped().Model(&entity.Article{}).
		Where("slug = ?", slug).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// Create 创建文章并在事务内重设标签关联
func (r *ArticleRepository) Create(ctx context.Context, article *entity.Article, tagIDs []uint) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Omit("Tags").Create(article).Error; err != nil {
			return err
		}
		if err := tx.Model(article).Association("Tags").Replace(tagRefs(tagIDs)); err != nil {
			return err
		}
		return nil
	})
}

// Update 保存文章变更并在事务内重设标签关联（先清旧关联再写新关联）
func (r *ArticleRepository) Update(ctx context.Context, article *entity.Article, tagIDs []uint) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&entity.Article{}).
			Where("id = ?", article.ID).
			Select("title", "summary", "category_id", "content", "cover_image", "status", "published_at").
			Updates(article).Error; err != nil {
			return err
		}
		if err := tx.Model(article).Association("Tags").Replace(tagRefs(tagIDs)); err != nil {
			return err
		}
		return nil
	})
}

// Delete 软删除文章（gorm.DeletedAt）
func (r *ArticleRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&entity.Article{}, id).Error
}

// UpdateStatus 更新文章状态（是否写入 published_at 由 Service 编排）
func (r *ArticleRepository) UpdateStatus(ctx context.Context, article *entity.Article) error {
	updates := map[string]any{"status": article.Status}
	if article.PublishedAt != nil {
		updates["published_at"] = article.PublishedAt
	}
	return r.db.WithContext(ctx).
		Model(&entity.Article{}).
		Where("id = ?", article.ID).
		Updates(updates).Error
}

// IncrViewCount 使用原子表达式自增浏览量，避免「读-改-写」
func (r *ArticleRepository) IncrViewCount(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).
		Model(&entity.Article{}).
		Where("id = ?", id).
		Update("view_count", gorm.Expr("view_count + 1")).Error
}

// Archive 按年月聚合已发布文章数量
func (r *ArticleRepository) Archive(ctx context.Context) ([]vo.ArchiveItemVO, error) {
	archives := make([]vo.ArchiveItemVO, 0)
	if err := r.db.WithContext(ctx).Model(&entity.Article{}).
		Select("strftime('%Y', published_at) as year, strftime('%m', published_at) as month, count(*) as count").
		Where("status = ? and published_at is not null", constant.ArticleStatusPublished).
		Group("strftime('%Y-%m', published_at)").
		Order("year DESC, month DESC").
		Scan(&archives).Error; err != nil {
		return nil, err
	}
	return archives, nil
}

// Hot 按浏览量倒序取 N 篇已发布文章
func (r *ArticleRepository) Hot(ctx context.Context, limit int) ([]entity.Article, error) {
	articles := make([]entity.Article, 0)
	if err := r.db.WithContext(ctx).Model(&entity.Article{}).
		Preload("Author").Preload("Category").Preload("Tags").
		Where("status = ?", constant.ArticleStatusPublished).
		Order("view_count DESC").
		Limit(limit).
		Find(&articles).Error; err != nil {
		return nil, err
	}
	counts, err := r.countComments(ctx, articleIDs(articles))
	if err != nil {
		return nil, err
	}
	for i := range articles {
		articles[i].CommentCount = counts[articles[i].ID]
	}
	return articles, nil
}

// CountByCategory 统计分类下文章数（供删除分类前校验）
func (r *ArticleRepository) CountByCategory(ctx context.Context, categoryID uint) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&entity.Article{}).
		Where("status = ? and category_id = ?", constant.ArticleStatusPublished, categoryID).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// countComments 批量统计文章评论数（一次 IN + GROUP BY），返回 articleID -> 评论数
func (r *ArticleRepository) countComments(ctx context.Context, articleIDs []uint) (map[uint]int64, error) {
	result := make(map[uint]int64, len(articleIDs))
	if len(articleIDs) == 0 {
		return result, nil
	}
	var rows []struct {
		ArticleID uint
		Total     int64
	}
	if err := r.db.WithContext(ctx).Model(&entity.Comment{}).
		Select("article_id, count(*) as total").
		Where("article_id IN ?", articleIDs).
		Group("article_id").
		Scan(&rows).Error; err != nil {
		return nil, err
	}
	for _, row := range rows {
		result[row.ArticleID] = row.Total
	}
	return result, nil
}

// articleIDs 抽取文章主键切片，供批量聚合使用
func articleIDs(articles []entity.Article) []uint {
	ids := make([]uint, 0, len(articles))
	for _, article := range articles {
		ids = append(ids, article.ID)
	}
	return ids
}

// tagRefs 把标签主键切片转成只带 ID 的最小实体，供多对多关联 Replace 使用
//   - GORM 的 Association.Replace 不接受 []uint，必须是关联类型（或其实体切片）
func tagRefs(tagIDs []uint) []entity.Tag {
	tags := make([]entity.Tag, 0, len(tagIDs))
	for _, id := range tagIDs {
		tags = append(tags, entity.Tag{ID: id})
	}
	return tags
}
