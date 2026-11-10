package handlers

import (
	"github.com/gin-gonic/gin"

	"blog-api/models/dto"
	"blog-api/services"
	"blog-api/utils"
)

// CategoryHandler 分类接口
type CategoryHandler struct {
	categoryService *services.CategoryService
}

// NewCategoryHandler 创建分类处理器
func NewCategoryHandler(categoryService *services.CategoryService) *CategoryHandler {
	return &CategoryHandler{categoryService: categoryService}
}

// List 分类列表（含 article_count）
func (h *CategoryHandler) List(c *gin.Context) {
	categories, err := h.categoryService.List(c.Request.Context())
	if err != nil {
		utils.HandleError(c, err)
		return
	}
	utils.Success(c, categories)
}

// Create 新增分类
func (h *CategoryHandler) Create(c *gin.Context) {
	var req dto.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, utils.ParseBindingError(err))
		return
	}

	category, err := h.categoryService.Create(c.Request.Context(), req)
	if err != nil {
		utils.HandleError(c, err)
		return
	}
	utils.Created(c, category)
}

// Update 修改分类
func (h *CategoryHandler) Update(c *gin.Context) {
	id, ok := parseIDParam(c)
	if !ok {
		return
	}

	var req dto.UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, utils.ParseBindingError(err))
		return
	}

	category, err := h.categoryService.Update(c.Request.Context(), id, req)
	if err != nil {
		utils.HandleError(c, err)
		return
	}
	utils.Success(c, category)
}

// Delete 删除分类
func (h *CategoryHandler) Delete(c *gin.Context) {
	id, ok := parseIDParam(c)
	if !ok {
		return
	}

	if err := h.categoryService.Delete(c.Request.Context(), id); err != nil {
		utils.HandleError(c, err)
		return
	}
	utils.Success(c, gin.H{"id": id})
}
