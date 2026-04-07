package request

import "slices"

type Pagination struct {
	Page     int    `query:"page" validate:"omitempty,min=1"`
	PerPage  int    `query:"per_page" validate:"omitempty,min=1,max=100"`
	OrderBy  string `query:"order_by"`
	OrderDir string `query:"order_dir" validate:"omitempty,oneof=asc desc"`
}

func (p *Pagination) SetDefaults(col string, dir string, wl []string) {
	if p.Page <= 0 {
		p.Page = 1
	}
	if p.PerPage <= 0 {
		p.PerPage = 10
	}
	if p.OrderBy == "" || !slices.Contains(wl, p.OrderBy) {
		p.OrderBy = col
	}
	if p.OrderDir == "" {
		p.OrderDir = dir
	}
}
