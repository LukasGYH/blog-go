package vo

import "time"

// CommentVO 评论出参（两级：顶层评论内嵌 replies）
type CommentVO struct {
	ID        uint        `json:"id"`
	Content   string      `json:"content"`
	ParentID  uint        `json:"parent_id"`
	CreatedAt time.Time   `json:"created_at"`
	User      UserBriefVO `json:"user"`
	Replies   []CommentVO `json:"replies"`
}
