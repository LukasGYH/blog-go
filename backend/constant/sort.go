package constant

// SortType 列表排序方式
type SortType string

const (
	SortLatest SortType = "latest" // 按发布时间倒序
	SortHot    SortType = "hot"    // 按浏览量倒序
)
