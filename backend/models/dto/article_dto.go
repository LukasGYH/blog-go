package dto

import (
	"blog-api/constant"
	"blog-api/utils"
)

// CreateArticleRequest 新建文章请求
type CreateArticleRequest struct {
	Title      string                 `json:"title" binding:"required,max=200"`
	Summary    string                 `json:"summary" binding:"omitempty,max=500"`
	Content    string                 `json:"content" binding:"required"`
	CoverImage string                 `json:"cover_image" binding:"omitempty,max=255"`
	CategoryID *uint                  `json:"category_id"`
	TagIDs     []uint                 `json:"tag_ids"`
	Status     constant.ArticleStatus `json:"status" binding:"omitempty,oneof=draft published"`
}

// UpdateArticleRequest 编辑文章请求
type UpdateArticleRequest struct {
	Title      string                 `json:"title" binding:"required,max=200"`
	Summary    string                 `json:"summary" binding:"omitempty,max=500"`
	Content    string                 `json:"content" binding:"required"`
	CoverImage string                 `json:"cover_image" binding:"omitempty,max=255"`
	CategoryID *uint                  `json:"category_id"`
	TagIDs     []uint                 `json:"tag_ids"`
	Status     constant.ArticleStatus `json:"status" binding:"omitempty,oneof=draft published"`
}

// UpdateArticleStatusRequest 切换文章状态请求
type UpdateArticleStatusRequest struct {
	Status constant.ArticleStatus `json:"status" binding:"required,oneof=draft published"`
}

// ArticleQuery 文章列表查询参数（query 绑定）
type ArticleQuery struct {
	Page       int                    `form:"page"`
	PageSize   int                    `form:"page_size"`
	Keyword    string                 `form:"keyword"`
	CategoryID uint                   `form:"category_id"`
	TagID      uint                   `form:"tag_id"`
	AuthorID   uint                   `form:"author_id"`
	Status     constant.ArticleStatus `form:"status"`
	Sort       constant.SortType      `form:"sort"`
}

// Normalize 规整查询参数：分页裁剪（page<1→1、page_size 默认 10 上限 50）并补默认排序
func (q *ArticleQuery) Normalize() {
	page := utils.PageQuery{Page: q.Page, PageSize: q.PageSize}
	page.Normalize()
	q.Page = page.Page
	q.PageSize = page.PageSize

	if q.Sort != constant.SortHot {
		q.Sort = constant.SortLatest
	}
}

// PageQuery 返回规整后的通用分页参数
func (q ArticleQuery) PageQuery() utils.PageQuery {
	return utils.PageQuery{Page: q.Page, PageSize: q.PageSize}
}
