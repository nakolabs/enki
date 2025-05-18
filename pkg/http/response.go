package http

import "math"

type Response struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message"`
	Errors  interface{} `json:"errors,omitempty"`
	Meta    interface{} `json:"meta,omitempty"`
}

func NewResponse() *Response {
	return &Response{}
}

func (r *Response) SetCode(code int) *Response {
	r.Code = code
	return r
}

func (r *Response) SetData(data interface{}) *Response {
	r.Data = data
	return r
}

func (r *Response) SetMessage(message string) *Response {
	r.Message = message
	return r
}

func (r *Response) SetErrors(e interface{}) *Response {
	r.Errors = e
	return r
}

func (r *Response) SetMeta(meta interface{}) *Response {
	r.Meta = meta
	return r
}

type PaginationResponse struct {
	TotalData int `json:"total_data"`
	TotalPage int `json:"total_page"`
	PageNum   int `json:"page_num"`
	PageSize  int `json:"page_size"`
}

type Meta struct {
	Pagination PaginationResponse     `json:"pagination"`
	Filter     map[string]interface{} `json:"filter"`
	Order      map[string]interface{} `json:"order"`
}

func NewMetaFromQuery(query Query, totalData int) *Meta {
	return &Meta{
		Pagination: PaginationResponse{
			PageNum:   query.PageNum,
			PageSize:  query.PageSize,
			TotalData: totalData,
			TotalPage: int(math.Ceil(float64(totalData) / float64(query.PageSize))),
		},
		Filter: map[string]interface{}{
			"search_by":  query.SearchBy,
			"search":     query.Search,
			"start_date": query.StartDate,
			"end_date":   query.EndDate,
		},
		Order: map[string]interface{}{
			"order":    query.Order,
			"order_by": query.OrderBy,
		},
	}
}
