package pagination

type Params struct {
	Page    int    `form:"page" json:"page" query:"page"`
	PerPage int    `form:"per_page" json:"per_page" query:"per_page"`
	Sort    string `form:"sort" json:"sort" query:"sort"`
	Order   string `form:"order" json:"order" query:"order"`
}

func (p *Params) SetDefaults() {
	if p.Page < 1 {
		p.Page = 1
	}

	if p.PerPage < 1 {
		p.PerPage = 15
	}

	if p.Order == "" {
		p.Order = "asc"
	}
}
