package dto

// CreateCategoryRequest 新增分类请求
type CreateCategoryRequest struct {
	Name        string `json:"name" binding:"required,max=50"`
	Slug        string `json:"slug" binding:"omitempty,max=50"`
	Description string `json:"description" binding:"omitempty,max=255"`
	SortOrder   int    `json:"sort_order"`
}

// UpdateCategoryRequest 修改分类请求
type UpdateCategoryRequest struct {
	Name        string `json:"name" binding:"required,max=50"`
	Slug        string `json:"slug" binding:"omitempty,max=50"`
	Description string `json:"description" binding:"omitempty,max=255"`
	SortOrder   int    `json:"sort_order"`
}
