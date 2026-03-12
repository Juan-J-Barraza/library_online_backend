package response

type DashboardResponse struct {
	TotalBooks     int `json:"total_books"`
	AvailableBooks int `json:"available_books"`
	BorrowedBooks  int `json:"borrowed_books"`
	ReservedBooks  int `json:"reserved_books"`
}
