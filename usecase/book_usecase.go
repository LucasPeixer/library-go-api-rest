package usecase

import (
	"go-api/model"
	"go-api/repository"

	"github.com/gin-gonic/gin"
)

type BookUsecase struct {
	repository repository.BookRepository
}

func NewBookUseCase(repo repository.BookRepository) BookUsecase {
	return BookUsecase{
		repository: repo,
	}
}

func (bu *BookUsecase) GetAllBooks() ([]model.Book,error){
	return bu.repository.GetAllBooks()
}

func (bu *BookUsecase) GetBooks(ctx *gin.Context) ([]model.Book, error){
	books, err := bu.repository.GetBooks(ctx)

	if err != nil {
		return nil, err
	}	
	return books, nil
}

func (bu *BookUsecase) CreateBook(book model.Book) (string,error){
	lastInsertID, err := bu.repository.CreateBook(book)
	if(err != nil){
		return "", err
	}

	return lastInsertID, nil
}

func (bu *BookUsecase) DeleteBook(ctx *gin.Context) (string,error){
	lastDeleteBook, err := bu.repository.DeleteBook(ctx)
	if(err != nil){
		return "", err
	}

	return lastDeleteBook, nil
}

func (bu* BookUsecase) UpdateBook(book model.Book) error{
	return bu.repository.UpdateBook(book)
}
