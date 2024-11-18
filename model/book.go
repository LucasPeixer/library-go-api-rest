package model

type Book struct {
	Id       int          `json:"id"`
	Title    string       `json:"title"`
	Synopsis string       `json:"synopsis"`
	Amount   int          `json:"amount"`
	Stock    *[]BookStock `json:"stock"` // Estoque pode ser omitido com null
	Author   *Author      `json:"author"`
	Genres   []Genre      `json:"genres"`
}

// NewBook cria uma nova instância de Book.
func NewBook(id int, title string, synopsis string, stock *[]BookStock, author *Author, genres []Genre) *Book {
	book := new(Book)
	book.Id = id
	book.Title = title
	book.Synopsis = synopsis

	if stock != nil {
		book.Amount = len(*stock)
	}

	book.Stock = stock
	book.Author = author
	book.Genres = genres
	return book
}
