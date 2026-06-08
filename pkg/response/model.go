package response

type JSONResponse struct {
	Status   string `json:"status"`
	Message  string `json:"message"`
	Data     any    `json:"data,omitempty"`
	Metadata any    `json:"metadata,omitempty"`
	Errors   any    `json:"errors,omitempty"`
}

type Metadata struct {
	Total    int     `json:"total" example:"150"`
	Page     int     `json:"page" example:"1"`
	Limit    int     `json:"limit" example:"10"`
	OrderBy  *string `json:"order_by,omitempty" example:"name"`
	OrderDir *string `json:"order_ir,omitempty" example:"asc"`
	Search   *string `json:"string,omitempty" example:"Elephant"`
}
