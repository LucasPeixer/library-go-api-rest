package model

type Book struct {
	Id        string  `json:"id_book" binding:"required"`
	Title     string  `json:"title" binding:"required"`
	Synopsis  string  `json:"synopsis" binding:"required"`
	Price     float64 `json:"price" binding:"required"`
	Amount    int     `json:"amount" binding:"required"`
	Author_id string  `json:"author_id" binding:"required"`
}