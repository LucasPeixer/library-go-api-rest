package usecase

import (
	"go-api/model"
	"go-api/repository"
)

type BookUsecase struct {
	repository repository.BookRepository
}

func NewBookUseCase(repo repository.BookRepository) BookUsecase {
	return BookUsecase{
		repository: repo,
	}
}

func (bu *BookUsecase) GetBooks() ([]model.Book,error){
	return bu.repository.GetBooks()
}

func (bu *BookUsecase) CreateBook(book model.Book) (string,error){
	lastInsertID, err := bu.repository.CreateBook(book)
	if(err != nil){
		return "", err
	}

	return lastInsertID, nil
}
