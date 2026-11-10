package handlers

import (
	"github.com/gin-gonic/gin"

	"blog-api/models/dto"
	"blog-api/services"
	"blog-api/utils"
)

// TagHandler 标签接口
type TagHandler struct {
	tagService *services.TagService
}

// NewTagHandler 创建标签处理器
func NewTagHandler(tagService *services.TagService) *TagHandler {
	return &TagHandler{tagService: tagService}
}

// List 标签列表（支持关键字搜索，含文章数，返回裸数组，与 /categories 一致）
func (h *TagHandler) List(c *gin.Context) {
	var query dto.TagQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		utils.ValidationError(c, utils.ParseBindingError(err))
		return
	}
	list, err := h.tagService.List(c.Request.Context(), query.Keyword)
	if err != nil {
		utils.HandleError(c, err)
		return
	}
	utils.Success(c, list)
}

// Create 新增标签
func (h *TagHandler) Create(c *gin.Context) {
	var req dto.CreateTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, utils.ParseBindingError(err))
		return
	}
	tag, err := h.tagService.Create(c.Request.Context(), req)
	if err != nil {
		utils.HandleError(c, err)
		return
	}
	utils.Created(c, tag)
}

// Update 修改标签
func (h *TagHandler) Update(c *gin.Context) {
	id, ok := parseIDParam(c)
	if !ok {
		return
	}
	var req dto.UpdateTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, utils.ParseBindingError(err))
		return
	}
	tag, err := h.tagService.Update(c.Request.Context(), id, req)
	if err != nil {
		utils.HandleError(c, err)
		return
	}
	utils.Success(c, tag)
}

// Delete 删除标签
func (h *TagHandler) Delete(c *gin.Context) {
	id, ok := parseIDParam(c)
	if !ok {
		return
	}
	if err := h.tagService.Delete(c.Request.Context(), id); err != nil {
		utils.HandleError(c, err)
		return
	}
	utils.Success(c, gin.H{"id": id})
}
