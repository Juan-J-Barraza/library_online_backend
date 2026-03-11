package utils

type Pagination struct {
	Page       int `query:"page" json:"page"`
	PageSize   int `query:"page_size" json:"page_size"`
	TotalItems int `json:"total_items"`
	TotalPages int `json:"total_pages"`
	Data       any `json:"data"`
}

func (p *Pagination) Calculate() {
	if p.Page <= 0 {
		p.Page = 1
	}
	if p.PageSize <= 0 {
		p.PageSize = 10
	}
	p.TotalPages = (p.TotalItems + p.PageSize - 1) / p.PageSize

}
