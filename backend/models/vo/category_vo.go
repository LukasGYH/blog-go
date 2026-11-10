package vo

// CategoryVO 分类列表出参（含文章数）
type CategoryVO struct {
	ID           uint   `json:"id"`
	Name         string `json:"name"`
	Slug         string `json:"slug"`
	Description  string `json:"description"`
	SortOrder    int    `json:"sort_order"`
	ArticleCount int64  `json:"article_count"`
}

// CategoryBriefVO 嵌套在文章中的分类精简信息
type CategoryBriefVO struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}
