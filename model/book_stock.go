package model

type BookStock struct {
	Id     int    `json:"id"`
	Status string `json:"status"`
	Code   int    `json:"code"`
	BookId int    `json:"book_id"`
}
