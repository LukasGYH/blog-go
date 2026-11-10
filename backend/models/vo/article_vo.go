package vo

import (
	"time"

	"blog-api/constant"
)

// ArticleCardVO 文章卡片（列表用，不含正文）
type ArticleCardVO struct {
	ID           uint                   `json:"id"`
	Title        string                 `json:"title"`
	Slug         string                 `json:"slug"`
	Summary      string                 `json:"summary"`
	CoverImage   string                 `json:"cover_image"`
	Status       constant.ArticleStatus `json:"status"`
	ViewCount    int                    `json:"view_count"`
	PublishedAt  *time.Time             `json:"published_at"`
	CreatedAt    time.Time              `json:"created_at"`
	Author       UserBriefVO            `json:"author"`
	Category     *CategoryBriefVO       `json:"category"`
	Tags         []TagBriefVO           `json:"tags"`
	CommentCount int64                  `json:"comment_count"`
}

// ArticleDetailVO 文章详情（在卡片字段基础上增加正文与更新时间）
type ArticleDetailVO struct {
	ArticleCardVO
	Content   string    `json:"content"`
	UpdatedAt time.Time `json:"updated_at"`
}
