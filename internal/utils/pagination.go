package utils

// PaginationRequest struct
type PaginationRequest struct {
	Page         int32 `json:"page,omitempty" validate:"number,gte=0"`
	ItemsPerPage int32 `json:"items_per_page,omitempty" validate:"number,gte=0"`
}

// Limit of items
func (r PaginationRequest) Limit() int32 {
	if r.ItemsPerPage > 0 {
		return r.ItemsPerPage
	}
	return 20
}

// Offset items
func (r PaginationRequest) Offset() int32 {
	if r.Page > 1 {
		return (r.Page - 1) * r.Limit()
	}
	return 0
}
