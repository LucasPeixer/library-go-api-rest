package repository

import (
	"database/sql"
	"fmt"
	"go-api/model"
)

type BookRepository struct {
	connection *sql.DB
}

func NewBookRepository (connection *sql.DB) BookRepository {
	return BookRepository{
		connection: connection,
	}
}

func (pr *BookRepository) GetBooks() ([]model.Book, error){
	query := "SELECT * FROM books"
	rows, err := pr.connection.Query(query)
	if(err != nil){
		fmt.Println(err)
		return []model.Book{}, err
	}

	var bookList []model.Book
	var bookObj model.Book

	for rows.Next(){
		err = rows.Scan(
			&bookObj.ID,
			&bookObj.Title,
			&bookObj.Synopsis,
			&bookObj.Price,
			&bookObj.Amount,
			&bookObj.Author_id)

		if(err != nil){
			fmt.Println(err)
		return []model.Book{}, err
		}

		bookList = append(bookList, bookObj)
	}

	rows.Close()

	return bookList, nil
}

func (br *BookRepository) CreateBook(book model.Book) (string, error){
	query := `
		INSERT INTO books (title, synopsis, price, amount, author_id) 
		VALUES ($1, $2, $3, $4, $5) RETURNING id`

		var lastInsertID string

		err := br.connection.QueryRow(query, book.Title, book.Synopsis, book.Price, book.Amount, book.Author_id).Scan(&lastInsertID)
		if err != nil {
			fmt.Println(err)
			return "", err
		}

		return lastInsertID, nil
}