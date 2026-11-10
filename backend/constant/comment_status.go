package constant

// CommentStatus 评论状态
type CommentStatus string

const (
	CommentStatusVisible CommentStatus = "visible" // 可见
	CommentStatusHidden  CommentStatus = "hidden"  // 已隐藏
)
