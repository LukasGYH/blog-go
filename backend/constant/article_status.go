package constant

// ArticleStatus 文章状态
type ArticleStatus string

const (
	ArticleStatusDraft     ArticleStatus = "draft"     // 草稿
	ArticleStatusPublished ArticleStatus = "published" // 已发布
)
