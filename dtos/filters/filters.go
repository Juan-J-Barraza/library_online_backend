package filters

type FiltersUser struct {
	Name     string `json:"name"`
	LastName string `json:"last_name"`
	Role     string `json:"role"`
}

type FiltersBook struct {
	Title       string `query:"title"`
	EditorialId uint   `query:"editorial_id"`
	AuthorId    uint   `query:"author_id"`
}
