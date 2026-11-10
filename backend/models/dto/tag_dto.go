package dto

// CreateTagRequest 新增标签请求
type CreateTagRequest struct {
	Name string `json:"name" binding:"required,max=50"`
}

// UpdateTagRequest 修改标签请求
type UpdateTagRequest struct {
	Name string `json:"name" binding:"required,max=50"`
}

// TagQuery 标签列表查询参数（query 绑定）
type TagQuery struct {
	Keyword string `form:"keyword"`
}
