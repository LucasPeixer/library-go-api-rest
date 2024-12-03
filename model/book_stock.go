package model

type BookStockStatus string

const (
	BookStockAvailable BookStockStatus = "available"
	BookStockBorrowed  BookStockStatus = "borrowed"
	BookStockMissing   BookStockStatus = "missing"
)

type BookStock struct {
	Id     int             `json:"id"`
	Status BookStockStatus `json:"status,omitempty"`
	Code   int             `json:"code"`
	BookId int             `json:"book_id"`
}
