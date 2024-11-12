package model

type Book struct {
	Id       int     `json:"id"`
	Title    string  `json:"title"`
	Synopsis string  `json:"synopsis"`
	Amount   int     `json:"amount"`
	Author   *Author `json:"author"` // Author can be null
	Genres   []Genre `json:"genres"` // Book can have zero genres
}

// NewBook is a constructor function to create a new Book instance.
func NewBook(id int, title string, synopsis string, amount int, author *Author, genres []Genre) *Book {
	return &Book{
		Id:       id,
		Title:    title,
		Synopsis: synopsis,
		Amount:   amount,
		Author:   author,
		Genres:   genres,
	}
}
