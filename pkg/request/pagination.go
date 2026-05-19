package request

type Pagination struct {
	Page     int    `query:"page" validate:"omitempty,min=1"`
	Limit    int    `query:"limit" validate:"omitempty,min=1,max=100"`
	Search   string `query:"search"`
	OrderBy  string `query:"order_by"`
	OrderDir string `query:"order_dir" validate:"omitempty,oneof=asc desc ASC DESC"`
}

func (p *Pagination) SetDefaults(col string, dir string) {
	if p.Page <= 0 {
		p.Page = 1
	}
	if p.Limit <= 0 {
		p.Limit = 10
	}
	if p.OrderBy == "" {
		p.OrderBy = col
	}
	if p.OrderDir == "" {
		p.OrderDir = dir
	}
}
