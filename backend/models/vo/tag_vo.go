package vo

// TagVO 标签列表出参（含文章数）
type TagVO struct {
	ID           uint   `json:"id"`
	Name         string `json:"name"`
	Slug         string `json:"slug"`
	ArticleCount int64  `json:"article_count"`
}

// TagBriefVO 嵌套在文章中的标签精简信息
type TagBriefVO struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}
