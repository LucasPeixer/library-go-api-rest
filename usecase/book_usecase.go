package usecase

import (
	"go-api/model"
	"go-api/repository"
)

type BookUseCase interface {
	CreateBook(title, synopsis string, bookCodes []int, authorId int, genreIds []int) (*model.Book, error)
	GetBooks(title, author string, genres []string) (*[]model.Book, error)
	GetBookById(id int) (*model.Book, error)
	UpdateBook(id int, title, synopsis string, authorId int) error
	DeleteBook(id int) error
	AddStock(code, bookId int) (*model.BookStock, error)
	GetStock(code *int, bookId int) (*[]model.BookStock, error)
	UpdateStockStatus(id int, status string) error
	RemoveStock(id int) error
}

type bookUseCase struct {
	repository repository.BookRepository
}

func NewBookUseCase(repository repository.BookRepository) BookUseCase {
	return &bookUseCase{repository: repository}
}

func (uc *bookUseCase) CreateBook(title, synopsis string, bookCodes []int, authorId int, genreIds []int) (*model.Book, error) {
	return uc.repository.CreateBook(title, synopsis, bookCodes, authorId, genreIds)
}

func (uc *bookUseCase) GetBooks(title, author string, genres []string) (*[]model.Book, error) {
	return uc.repository.GetBooks(title, author, genres)
}

func (uc *bookUseCase) GetBookById(id int) (*model.Book, error) {
	return uc.repository.GetBookById(id)
}

func (uc *bookUseCase) UpdateBook(id int, title, synopsis string, authorId int) error {
	return uc.repository.UpdateBook(id, title, synopsis, authorId)
}

func (uc *bookUseCase) DeleteBook(id int) error {
	return uc.repository.DeleteBook(id)
}

func (uc *bookUseCase) AddStock(code, bookId int) (*model.BookStock, error) {
	return uc.repository.AddStock(code, bookId)
}

func (uc *bookUseCase) GetStock(code *int, bookId int) (*[]model.BookStock, error) {
	return uc.repository.GetStock(code, bookId)
}

func (uc *bookUseCase) UpdateStockStatus(id int, status string) error {
	//TODO implement me
	panic("implement me")
}

func (uc *bookUseCase) RemoveStock(id int) error {
	return uc.repository.RemoveStock(id)
}
