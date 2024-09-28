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
			&bookObj.Amount)
		if(err != nil){
			fmt.Println(err)
		return []model.Book{}, err
		}

		bookList = append(bookList, bookObj)
	}

	rows.Close()

	return bookList, nil
}