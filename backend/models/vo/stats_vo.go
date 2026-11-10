package vo

// StatsOverviewVO 后台概览统计出参
type StatsOverviewVO struct {
	ArticleCount   int64 `json:"article_count"`
	PublishedCount int64 `json:"published_count"`
	DraftCount     int64 `json:"draft_count"`
	UserCount      int64 `json:"user_count"`
	CommentCount   int64 `json:"comment_count"`
	CategoryCount  int64 `json:"category_count"`
	TagCount       int64 `json:"tag_count"`
}

// ArchiveItemVO 归档统计项（按年月聚合）
type ArchiveItemVO struct {
	Year  int   `json:"year"`
	Month int   `json:"month"`
	Count int64 `json:"count"`
}

// CategoryStatsVO 分类维度统计项
type CategoryStatsVO struct {
	CategoryID   uint   `json:"category_id"`
	Name         string `json:"name"`
	ArticleCount int64  `json:"article_count"`
}
