package response

type JSONResponse struct {
	Status   string `json:"status"`
	Message  string `json:"message"`
	Data     any    `json:"data,omitempty"`
	Metadata any    `json:"metadata,omitempty"`
	Errors   any    `json:"errors,omitempty"`
}

type Metadata struct {
	Total     int    `json:"total" example:"150"`
	Filtered  int    `json:"filtered" example:"67"`
	Page      int    `json:"page" example:"1"`
	PerPage   int    `json:"per_page" example:"10"`
	TotalPage int    `json:"total_page" example:"7"`
	OrderBy   string `json:"order_by" example:"name"`
	OrderDir  string `json:"order_dir" example:"asc"`
}
