package utils

const (
	defaultPageSize = 10 // 每页默认条数
	maxPageSize     = 50 // 每页最大条数
)

// PageQuery 通用分页参数
type PageQuery struct {
	Page     int `form:"page" json:"page"`
	PageSize int `form:"page_size" json:"page_size"`
}

// Normalize 规整分页参数：page 最小为 1；page_size 默认 10、上限 50
func (p *PageQuery) Normalize() {
	if p.Page < 1 {
		p.Page = 1
	}
	if p.PageSize < 1 {
		p.PageSize = defaultPageSize
	}
	if p.PageSize > maxPageSize {
		p.PageSize = maxPageSize
	}
}

// Offset 计算 SQL 偏移量
func (p PageQuery) Offset() int {
	return (p.Page - 1) * p.PageSize
}
