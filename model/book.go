package model

type Book struct {
	ID        string  `json:"id_book"`
	Title     string  `json:"title"`
	Synopsis  string  `json:"synopsis"`
	Price     float64 `json:"price"`
	Amount    int     `json:"amount"`
	Author_id int     `json:"Author_id"`
}