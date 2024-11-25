package usecase

import (
	"go-api/model"
	"go-api/repository"
)

type BookUseCase interface {
	CreateBook(title, synopsis string, authorId int, genreIds []int) (*model.Book, error)
	GetBooks(title, author string, genres []string) (*[]model.Book, error)
	GetBookById(id int) (*model.Book, error)
	UpdateBook(id int, title, synopsis string, authorId int) error
	DeleteBook(id int) error
	AddStock(code, bookId int) (*model.BookStock, error)
	GetStock(code *int, bookId int) (*[]model.BookStock, error)
	UpdateStockStatus(id int, status string) error
	RemoveStock(id int, bookId *int) error
	CountAvailableBookStockById(bookId int) (int, error)
}

type bookUseCase struct {
	repository            repository.BookRepository
	reservationRepository repository.ReservationRepositoryInterface
}

func NewBookUseCase(repository repository.BookRepository, reservationRepo repository.ReservationRepositoryInterface) BookUseCase {
	return &bookUseCase{repository: repository, reservationRepository: reservationRepo}
}

func (uc *bookUseCase) CreateBook(title, synopsis string, authorId int, genreIds []int) (*model.Book, error) {
	return uc.repository.CreateBook(title, synopsis, authorId, genreIds)
}

func (uc *bookUseCase) GetBooks(title, author string, genres []string) (*[]model.Book, error) {
	books, err := uc.repository.GetBooks(title, author, genres)
	if err != nil {
		return nil, err
	}

	for i := range *books {
		book := &(*books)[i] // Get a pointer to the actual book in the slice

		amount, err := uc.CountAvailableBookStockById(book.Id)
		if err != nil {
			return nil, err
		}
		book.Amount = amount
	}
	return books, nil
}

func (uc *bookUseCase) GetBookById(id int) (*model.Book, error) {
	book, err := uc.repository.GetBookById(id)
	if err != nil {
		return nil, err
	}

	amount, err := uc.CountAvailableBookStockById(book.Id)
	if err != nil {
		return nil, err
	}
	book.Amount = amount
	return book, nil
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

func (uc *bookUseCase) RemoveStock(id int, bookId *int) error {
	return uc.repository.RemoveStock(id, bookId)
}

func (uc *bookUseCase) CountAvailableBookStockById(bookId int) (int, error) {
	availableBookStockCount, err := uc.CountBookStock(bookId, "available")
	if err != nil {
		return 0, err
	}

	pendingReservationsCount, err := uc.CountPendingReservationsByBookId(bookId)
	if err != nil {
		return 0, err
	}

	return availableBookStockCount - pendingReservationsCount, nil
}

func (uc *bookUseCase) CountBookStock(bookId int, status string) (int, error) {
	bookStockList, err := uc.GetStock(nil, bookId)
	if err != nil {
		return 0, err
	}

	var count int

	for _, bookStock := range *bookStockList {
		if status != "" {
			if bookStock.Status == status {
				count++
			}
		} else {
			count++
		}
	}
	return count, nil
}

func (uc *bookUseCase) CountPendingReservationsByBookId(bookId int) (int, error) {
	reservations, err := uc.reservationRepository.GetReservationsByBookId(bookId, "pending")
	if err != nil {
		return 0, err
	}
	return len(*reservations), nil
}
