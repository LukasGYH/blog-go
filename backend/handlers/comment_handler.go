package handlers

import (
	"github.com/gin-gonic/gin"

	"blog-api/middleware"
	"blog-api/models/dto"
	"blog-api/services"
	"blog-api/utils"
)

// CommentHandler 评论接口
type CommentHandler struct {
	commentService *services.CommentService
}

// NewCommentHandler 创建评论处理器
func NewCommentHandler(commentService *services.CommentService) *CommentHandler {
	return &CommentHandler{commentService: commentService}
}

// ListByArticle 评论列表（顶层 + 内嵌 replies）
func (h *CommentHandler) ListByArticle(c *gin.Context) {
	articleID, ok := parseIDParam(c)
	if !ok {
		return
	}
	var page utils.PageQuery
	if err := c.ShouldBindQuery(&page); err != nil {
		utils.ValidationError(c, utils.ParseBindingError(err))
		return
	}
	page.Normalize()
	list, total, err := h.commentService.ListByArticle(c.Request.Context(), articleID, page)
	if err != nil {
		utils.HandleError(c, err)
		return
	}
	utils.SuccessPage(c, list, total, page.Page, page.PageSize)
}

// Create 发表评论 / 回复
func (h *CommentHandler) Create(c *gin.Context) {
	articleID, ok := parseIDParam(c)
	if !ok {
		return
	}
	var req dto.CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, utils.ParseBindingError(err))
		return
	}
	comment, err := h.commentService.Create(c.Request.Context(), articleID, req, middleware.CurrentUserID(c))
	if err != nil {
		utils.HandleError(c, err)
		return
	}
	utils.Created(c, comment)
}

// Delete 删除评论
func (h *CommentHandler) Delete(c *gin.Context) {
	id, ok := parseIDParam(c)
	if !ok {
		return
	}
	if err := h.commentService.Delete(c.Request.Context(), id, middleware.CurrentUserID(c), middleware.CurrentUserRole(c)); err != nil {
		utils.HandleError(c, err)
		return
	}
	utils.Success(c, gin.H{"id": id})
}
