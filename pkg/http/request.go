package http

type PaginationQuery struct {
	PageSize int `form:"page_size"`
	PageNum  int `form:"page_num"`
}

func (p *PaginationQuery) GetOffset() int {
	return (p.PageNum - 1) * p.PageSize
}

type SearchQuery struct {
	Search   string   `form:"search"`
	SearchBy []string `form:"search_by"`
}

type DateRange struct {
	StartDate int64 `form:"start_date"`
	EndDate   int64 `form:"end_date"`
}

type OrderQuery struct {
	Order   string `form:"order" validate:"oneof='asc' 'desc'"`
	OrderBy string `form:"order_by"`
}

type Query struct {
	SearchQuery
	OrderQuery
	PaginationQuery
	DateRange
}

func (q Query) Get() (Query, map[string]interface{}) {
	return q, nil
}

func DefaultQuery() Query {
	return Query{
		SearchQuery: SearchQuery{
			Search:   "",
			SearchBy: []string{},
		},
		OrderQuery: OrderQuery{
			Order:   "desc",
			OrderBy: "created_at",
		},
		PaginationQuery: PaginationQuery{
			PageSize: 10,
			PageNum:  1,
		},
		DateRange: DateRange{
			StartDate: 0,
			EndDate:   0,
		},
	}
}
