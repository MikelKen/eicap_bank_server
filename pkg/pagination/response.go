package pagination

import "math"

type Response[T any] struct {
	Items      []T   `json:"items"`
	Total      int64 `json:"total"`
	Page       int   `json:"page"`
	PerPage    int   `json:"per_page"`
	TotalPages int   `json:"total_pages"`
}

func NewResponse[T any](items []T, total int64, params Params) Response[T] {
	totalPages := int(math.Ceil(float64(total) / float64(params.PerPage)))

	if items == nil {
		items = []T{}
	}

	return Response[T]{
		Items:      items,
		Total:      total,
		Page:       params.Page,
		PerPage:    params.PerPage,
		TotalPages: totalPages,
	}
}
