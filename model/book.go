package model

type Book struct {
	ID        string  `json:"id_book"`
	Title     string  `json:"title" binding:"required"`
	Synopsis  string  `json:"synopsis" binding:"required"`
	Price     float64 `json:"price" binding:"required"`
	Amount    int     `json:"amount" binding:"required"`
	Author_id int     `json:"author_id" binding:"required"`
}