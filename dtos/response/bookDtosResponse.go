package response

import (
	"libraryOnline/models"
	"time"
)

type AuthorResponse struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	LastName string `json:"last_name"`
}

type EditorialResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type BookResponse struct {
	ID                uint              `json:"id"`
	Title             string            `json:"title"`
	AvailableQuantity int               `json:"available_quantity"`
	TotalQuantity     int               `json:"total_quantity"`
	Image             string            `json:"image"`
	Editorial         EditorialResponse `json:"editorial"`
	Authors           []AuthorResponse  `json:"authors"`
	CreatedAt         time.Time         `json:"created_at"`
}

func ToBookResponse(b models.Book) BookResponse {
	authors := make([]AuthorResponse, len(b.Authors))
	for i, a := range b.Authors {
		authors[i] = AuthorResponse{ID: a.ID, Name: a.Name, LastName: a.LastName}
	}
	return BookResponse{
		ID:                b.ID,
		Title:             b.Title,
		AvailableQuantity: b.AvailableQuantity,
		TotalQuantity:     b.TotalQuantity,
		Image:             b.Image,
		Editorial:         EditorialResponse{ID: b.Editorial.ID, Name: b.Editorial.Name},
		Authors:           authors,
		CreatedAt:         b.CreatedAt,
	}
}
