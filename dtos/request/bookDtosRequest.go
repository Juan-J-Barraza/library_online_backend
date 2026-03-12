package request

type CreateOrUpdateBookRequest struct {
	Title             string `json:"title"              validate:"required,min=2,max=100"`
	AvailableQuantity int    `json:"available_quantity" validate:"required,min=0"`
	TotalQuantity     int    `json:"total_quantity"     validate:"required,min=1"`
	Image             string `json:"image"`
	EditorialId       uint   `json:"editorial_id"       validate:"required"`
	AuthorIds         []uint `json:"author_ids"         validate:"required,min=1"`
}
