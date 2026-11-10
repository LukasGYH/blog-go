package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"blog-api/middleware"
	"blog-api/models/dto"
	"blog-api/services"
	"blog-api/utils"
)

// defaultHotLimit 热门文章默认条数
const defaultHotLimit = 5

// ArticleHandler 文章接口
type ArticleHandler struct {
	articleService *services.ArticleService
}

// NewArticleHandler 创建文章处理器
func NewArticleHandler(articleService *services.ArticleService) *ArticleHandler {
	return &ArticleHandler{articleService: articleService}
}

// List 文章列表（分页 + 筛选 + 排序）
func (h *ArticleHandler) List(c *gin.Context) {
	var query dto.ArticleQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		utils.ValidationError(c, utils.ParseBindingError(err))
		return
	}
	query.Normalize()
	list, total, err := h.articleService.List(c.Request.Context(), query, middleware.CurrentUserID(c), middleware.CurrentUserRole(c))
	if err != nil {
		utils.HandleError(c, err)
		return
	}
	utils.SuccessPage(c, list, total, query.Page, query.PageSize)
}

// GetByID 文章详情
func (h *ArticleHandler) GetByID(c *gin.Context) {
	id, ok := parseIDParam(c)
	if !ok {
		return
	}
	article, err := h.articleService.GetByID(c.Request.Context(), id, middleware.CurrentUserID(c), middleware.CurrentUserRole(c))
	if err != nil {
		utils.HandleError(c, err)
		return
	}
	utils.Success(c, article)
}

// Create 新建文章
func (h *ArticleHandler) Create(c *gin.Context) {
	var req dto.CreateArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, utils.ParseBindingError(err))
		return
	}
	article, err := h.articleService.Create(c.Request.Context(), req, middleware.CurrentUserID(c))
	if err != nil {
		utils.HandleError(c, err)
		return
	}
	utils.Created(c, article)
}

// Update 编辑文章
func (h *ArticleHandler) Update(c *gin.Context) {
	id, ok := parseIDParam(c)
	if !ok {
		return
	}
	var req dto.UpdateArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, utils.ParseBindingError(err))
		return
	}
	article, err := h.articleService.Update(c.Request.Context(), id, req, middleware.CurrentUserID(c), middleware.CurrentUserRole(c))
	if err != nil {
		utils.HandleError(c, err)
		return
	}
	utils.Success(c, article)
}

// ChangeStatus 切换文章状态（草稿 ↔ 发布）
func (h *ArticleHandler) ChangeStatus(c *gin.Context) {
	id, ok := parseIDParam(c)
	if !ok {
		return
	}
	var req dto.UpdateArticleStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, utils.ParseBindingError(err))
		return
	}
	detailVO, err := h.articleService.ChangeStatus(c.Request.Context(), id, req.Status, middleware.CurrentUserID(c), middleware.CurrentUserRole(c))
	if err != nil {
		utils.HandleError(c, err)
		return
	}
	utils.Success(c, detailVO)
}

// Delete 删除文章（软删除）
func (h *ArticleHandler) Delete(c *gin.Context) {
	id, ok := parseIDParam(c)
	if !ok {
		return
	}
	if err := h.articleService.Delete(c.Request.Context(), id, middleware.CurrentUserID(c), middleware.CurrentUserRole(c)); err != nil {
		utils.HandleError(c, err)
		return
	}
	utils.Success(c, gin.H{"id": id})
}

// Archive 归档统计
func (h *ArticleHandler) Archive(c *gin.Context) {
	items, err := h.articleService.Archive(c.Request.Context())
	if err != nil {
		utils.HandleError(c, err)
		return
	}
	utils.Success(c, items)
}

// Hot 热门文章
func (h *ArticleHandler) Hot(c *gin.Context) {
	limit := defaultHotLimit
	if v, ok := c.GetQuery("limit"); ok {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			limit = n
		}
	}
	list, err := h.articleService.Hot(c.Request.Context(), limit)
	if err != nil {
		utils.HandleError(c, err)
		return
	}
	utils.Success(c, list)
}
