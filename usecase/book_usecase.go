package usecase

import (
	"go-api/model"
	"go-api/repository"
)

type BookUseCase interface {
	CreateBook(title, synopsis string, amount, authorId int, genreIds []int) (*model.Book, error)
	GetBooks(title, author string, genres []string) ([]model.Book, error)
	UpdateBook(id int, title, synopsis string, amount, authorId int) error
	DeleteBook(id int) error
}

type bookUseCase struct {
	repository repository.BookRepository
}

func NewBookUseCase(repository repository.BookRepository) BookUseCase {
	return &bookUseCase{repository: repository}
}

func (uc *bookUseCase) CreateBook(title, synopsis string, amount, authorId int, genreIds []int) (*model.Book, error) {
	return uc.repository.CreateBook(title, synopsis, amount, authorId, genreIds)
}

func (uc *bookUseCase) GetBooks(title, author string, genres []string) ([]model.Book, error) {
	return uc.repository.GetBooks(title, author, genres)
}

func (uc *bookUseCase) UpdateBook(id int, title, synopsis string, amount, authorId int) error {
	return uc.repository.UpdateBook(id, title, synopsis, amount, authorId)
}

func (uc *bookUseCase) DeleteBook(id int) error {
	return uc.repository.DeleteBook(id)
}
