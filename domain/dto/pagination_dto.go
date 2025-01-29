package dto

type PaginationResponse struct {
	TotalData int64 `json:"total_data"`
	TotalPage int   `json:"total_page"`
	Page      int   `json:"page"`
	Limit     int   `json:"limit"`
}
