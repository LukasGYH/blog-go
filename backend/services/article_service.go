package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	"blog-api/constant"
	"blog-api/models/dto"
	"blog-api/models/entity"
	"blog-api/models/vo"
	"blog-api/repositories"
	"blog-api/utils"
)

// summaryMaxRunes 摘要自动截取的最大字符数
const summaryMaxRunes = 120

// ArticleService 文章业务逻辑
type ArticleService struct {
	repo *repositories.ArticleRepository
}

// NewArticleService 创建文章服务
func NewArticleService(repo *repositories.ArticleRepository) *ArticleService {
	return &ArticleService{repo: repo}
}

// List 文章列表：非管理员强制只看已发布，DTO 转 VO
func (s *ArticleService) List(ctx context.Context, query dto.ArticleQuery, userID uint, role constant.Role) ([]vo.ArticleCardVO, int64, error) {
	// 非管理员（含游客）一律只能看已发布文章；管理员可按 status 参数筛选
	if role != constant.RoleAdmin {
		query.Status = constant.ArticleStatusPublished
	}
	list, total, err := s.repo.List(ctx, query)
	if err != nil {
		return nil, 0, fmt.Errorf("查询文章列表失败: %w", err)
	}
	res := make([]vo.ArticleCardVO, 0, len(list))
	for _, article := range list {
		res = append(res, s.toCardVO(article))
	}
	return res, total, nil
}

// GetByID 文章详情：越权校验、浏览量自增、DTO 转 VO
func (s *ArticleService) GetByID(ctx context.Context, id uint, userID uint, role constant.Role) (*vo.ArticleDetailVO, error) {
	article, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.NotFound("文章不存在")
		}
		return nil, fmt.Errorf("查询文章失败: %w", err)
	}
	// 草稿仅作者本人或管理员可见
	if article.Status != constant.ArticleStatusPublished && article.AuthorID != userID && role != constant.RoleAdmin {
		return nil, utils.Forbidden("您没有权限查看该文章")
	}
	// 浏览已发布文章时原子自增浏览量（A-03）
	if article.Status == constant.ArticleStatusPublished {
		if err := s.repo.IncrViewCount(ctx, article.ID); err != nil {
			return nil, fmt.Errorf("更新浏览量失败: %w", err)
		}
		article.ViewCount++
	}
	detail := s.toDetailVO(*article)
	return &detail, nil
}

// Create 新建文章：slug 唯一化、摘要补全、事务编排
func (s *ArticleService) Create(ctx context.Context, req dto.CreateArticleRequest, userID uint) (*vo.ArticleDetailVO, error) {
	// slug 由标题生成并做唯一化；必须接收返回值，否则冲突时不会追加 -1/-2
	var slugErr error
	slug := utils.EnsureUnique(utils.GenerateSlug(req.Title), func(candidate string) bool {
		exists, err := s.repo.SlugExists(ctx, candidate)
		if err != nil {
			slugErr = err
			// 出错时返回 false 让 EnsureUnique 结束循环，随后统一上报错误
			return false
		}
		return exists
	})
	if slugErr != nil {
		return nil, fmt.Errorf("校验 slug 唯一性失败: %w", slugErr)
	}

	summary := req.Summary
	if summary == "" {
		summary = truncateSummary(req.Content, summaryMaxRunes)
	}
	status := req.Status
	if status == "" {
		status = constant.ArticleStatusDraft
	}

	article := &entity.Article{
		Title:      req.Title,
		Slug:       slug,
		Summary:    summary,
		Content:    req.Content,
		Status:     status,
		CoverImage: req.CoverImage,
		AuthorID:   userID,
		CategoryID: req.CategoryID,
	}
	// 直接以已发布状态创建时，补写首次发布时间
	if status == constant.ArticleStatusPublished {
		now := time.Now()
		article.PublishedAt = &now
	}

	if err := s.repo.Create(ctx, article, req.TagIDs); err != nil {
		return nil, fmt.Errorf("创建文章失败: %w", err)
	}
	return s.detailByID(ctx, article.ID)
}

// Update 编辑文章：越权校验（作者或管理员）
func (s *ArticleService) Update(ctx context.Context, id uint, req dto.UpdateArticleRequest, userID uint, role constant.Role) (*vo.ArticleDetailVO, error) {
	article, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.NotFound("文章不存在")
		}
		return nil, fmt.Errorf("查询文章失败: %w", err)
	}
	if article.AuthorID != userID && role != constant.RoleAdmin {
		return nil, utils.Forbidden("您没有权限编辑该文章")
	}

	article.Title = req.Title
	article.Summary = req.Summary
	article.Content = req.Content
	article.CoverImage = req.CoverImage
	article.CategoryID = req.CategoryID
	if article.Summary == "" {
		article.Summary = truncateSummary(req.Content, summaryMaxRunes)
	}
	// 请求里带了 status 才更新，并遵循「首次发布写 published_at」
	if req.Status != "" {
		article.Status = req.Status
		if req.Status == constant.ArticleStatusPublished && article.PublishedAt == nil {
			now := time.Now()
			article.PublishedAt = &now
		}
	}

	if err := s.repo.Update(ctx, article, req.TagIDs); err != nil {
		return nil, fmt.Errorf("更新文章失败: %w", err)
	}
	return s.detailByID(ctx, id)
}

// Delete 删除文章（软删除）：越权校验（作者或管理员）
func (s *ArticleService) Delete(ctx context.Context, id uint, userID uint, role constant.Role) error {
	article, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.NotFound("文章不存在")
		}
		return fmt.Errorf("查询文章失败: %w", err)
	}
	if article.AuthorID != userID && role != constant.RoleAdmin {
		return utils.Forbidden("您没有权限删除该文章")
	}
	if err := s.repo.Delete(ctx, article.ID); err != nil {
		return fmt.Errorf("删除文章失败: %w", err)
	}
	return nil
}

// ChangeStatus 切换状态：越权校验，首次发布写 published_at
func (s *ArticleService) ChangeStatus(ctx context.Context, id uint, status constant.ArticleStatus, userID uint, role constant.Role) (*vo.ArticleDetailVO, error) {
	article, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.NotFound("文章不存在")
		}
		return nil, fmt.Errorf("查询文章失败: %w", err)
	}
	if article.AuthorID != userID && role != constant.RoleAdmin {
		return nil, utils.Forbidden("您没有权限切换该文章状态")
	}
	if article.Status == status {
		detail := s.toDetailVO(*article)
		return &detail, nil
	}

	article.Status = status
	if status == constant.ArticleStatusPublished && article.PublishedAt == nil {
		now := time.Now()
		article.PublishedAt = &now
	}
	if err := s.repo.UpdateStatus(ctx, article); err != nil {
		return nil, fmt.Errorf("更新文章状态失败: %w", err)
	}
	detail := s.toDetailVO(*article)
	return &detail, nil
}

// Archive 归档统计
func (s *ArticleService) Archive(ctx context.Context) ([]vo.ArchiveItemVO, error) {
	archives, err := s.repo.Archive(ctx)
	if err != nil {
		return nil, fmt.Errorf("查询归档统计失败: %w", err)
	}
	return archives, nil
}

// Hot 热门文章：limit 兜底与上限保护
func (s *ArticleService) Hot(ctx context.Context, limit int) ([]vo.ArticleCardVO, error) {
	if limit <= 0 {
		limit = 5
	}
	if limit > 50 {
		limit = 50
	}
	articles, err := s.repo.Hot(ctx, limit)
	if err != nil {
		return nil, fmt.Errorf("查询热门文章失败: %w", err)
	}
	res := make([]vo.ArticleCardVO, 0, len(articles))
	for _, article := range articles {
		res = append(res, s.toCardVO(article))
	}
	return res, nil
}

// detailByID 按主键重新读取文章并组装为详情 VO
func (s *ArticleService) detailByID(ctx context.Context, id uint) (*vo.ArticleDetailVO, error) {
	article, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("读取文章详情失败: %w", err)
	}
	detail := s.toDetailVO(*article)
	return &detail, nil
}

// toCardVO 实体转文章卡片 VO
func (s *ArticleService) toCardVO(article entity.Article) vo.ArticleCardVO {
	return vo.ArticleCardVO{
		ID:          article.ID,
		Title:       article.Title,
		Slug:        article.Slug,
		Summary:     article.Summary,
		CoverImage:  article.CoverImage,
		Status:      article.Status,
		ViewCount:   article.ViewCount,
		PublishedAt: article.PublishedAt,
		CreatedAt:   article.CreatedAt,
		Author: vo.UserBriefVO{
			ID:       article.Author.ID,
			Username: article.Author.Username,
			Nickname: article.Author.Nickname,
			Avatar:   article.Author.Avatar,
		},
		Category:     categoryBrief(article.Category),
		Tags:         tagBriefs(article.Tags),
		CommentCount: article.CommentCount,
	}
}

// toDetailVO 实体转文章详情 VO
func (s *ArticleService) toDetailVO(article entity.Article) vo.ArticleDetailVO {
	return vo.ArticleDetailVO{
		ArticleCardVO: s.toCardVO(article),
		Content:       article.Content,
		UpdatedAt:     article.UpdatedAt,
	}
}

// categoryBrief 分类转精简 VO，nil 安全
func categoryBrief(category *entity.Category) *vo.CategoryBriefVO {
	if category == nil || category.ID == 0 {
		return nil
	}
	return &vo.CategoryBriefVO{ID: category.ID, Name: category.Name, Slug: category.Slug}
}

// tagBriefs 标签切片转精简 VO（返回非 nil 空切片，避免 JSON 出现 null）
func tagBriefs(tags []entity.Tag) []vo.TagBriefVO {
	result := make([]vo.TagBriefVO, 0, len(tags))
	for _, tag := range tags {
		result = append(result, vo.TagBriefVO{ID: tag.ID, Name: tag.Name, Slug: tag.Slug})
	}
	return result
}

// truncateSummary 按字符（rune）截取正文作为摘要，超长时追加省略号
func truncateSummary(content string, maxRunes int) string {
	runes := []rune(content)
	if len(runes) <= maxRunes {
		return content
	}
	return string(runes[:maxRunes]) + "..."
}
