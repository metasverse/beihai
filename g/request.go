package g

import "github.com/eatmoreapple/regia"

type Pagination interface {
	Page() int
	PageSize() int
	Limit() int
	Offset() int
}

type QueryPagination struct {
	Ctx             *regia.Context
	MaxPageSize     int
	PageParam       string
	PageSizeParam   string
	DefaultPageSize int
}

func NewQueryPagination(ctx *regia.Context) Pagination {
	return &QueryPagination{
		Ctx:             ctx,
		MaxPageSize:     200,
		PageParam:       "page",
		PageSizeParam:   "page_size",
		DefaultPageSize: 10,
	}
}

func (q QueryPagination) Page() int {
	page, _ := q.Ctx.QueryValue(q.PageParam).Int(1)
	if page <= 0 {
		page = 1
	}
	return page
}

func (q QueryPagination) PageSize() int {
	pageSize, _ := q.Ctx.QueryValue(q.PageSizeParam).Int(q.DefaultPageSize)
	if pageSize > q.MaxPageSize {
		pageSize = q.MaxPageSize
	}
	if pageSize <= 0 {
		pageSize = q.DefaultPageSize
	}
	return pageSize
}

func (q QueryPagination) Limit() int {
	return q.PageSize()
}

func (q QueryPagination) Offset() int {
	return (q.Page() - 1) * q.PageSize()
}
