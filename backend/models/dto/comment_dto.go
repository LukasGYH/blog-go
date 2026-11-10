package dto

// CreateCommentRequest 发表评论 / 回复请求
type CreateCommentRequest struct {
	Content  string `json:"content" binding:"required,min=1,max=500"`
	ParentID uint   `json:"parent_id"` // 省略或 0 表示顶层评论
}
