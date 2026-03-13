package models

// PaginationInfo represents pagination metadata from API responses
type PaginationInfo struct {
	Cursor       string `json:"cursor,omitempty"`
	AfterCursor  string `json:"after_cursor,omitempty"`
	BeforeCursor string `json:"before_cursor,omitempty"`
	TotalPages   int    `json:"total_pages,omitempty"`
	CurrentPage  int    `json:"current_page,omitempty"`
	TotalCount   int    `json:"total_count,omitempty"`
}

// PagedResponse represents a paginated response with both data and pagination information
type PagedResponse struct {
	Data       interface{}    `json:"data"`
	Pagination PaginationInfo `json:"pagination"`
}