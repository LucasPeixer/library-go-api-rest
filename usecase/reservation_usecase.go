package usecase

import (
	"fmt"
	"go-api/model"
	"go-api/repository"
	"time"
)

type ReservationUseCase interface {
	CreateReservation(borrowedDays, userId, bookId int) (*model.Reservation, error)
	GetReservationsByFilters(userName string, status model.ReservationStatus, reservedAt string) (*[]model.Reservation, error)
	GetReservationById(id int) (*model.Reservation, error)
}

type reservationUseCase struct {
	reservationRepo repository.ReservationRepository
	userRepo        repository.UserRepository
	bookRepo        repository.BookRepository
}

// NewReservationUseCase cria e retorna uma nova instância de ReservationUseCase
func NewReservationUseCase(reservationRepo repository.ReservationRepository,
	userRepo repository.UserRepository,
	bookRepo repository.BookRepository) ReservationUseCase {
	return &reservationUseCase{
		reservationRepo: reservationRepo,
		userRepo:        userRepo,
		bookRepo:        bookRepo}
}

func (ru *reservationUseCase) GetReservationsByFilters(userName string, status model.ReservationStatus, reservedAt string) (*[]model.Reservation, error) {
	return ru.reservationRepo.GetReservationsByFilters(userName, status, reservedAt)
}

func (ru *reservationUseCase) CreateReservation(borrowedDays, userId, bookId int) (*model.Reservation, error) {
	user, err := ru.userRepo.GetUserById(userId)
	if err != nil {
		return nil, fmt.Errorf("error when searching for user: %w", err)
	}
	if user.IsActive != true {
		return nil, fmt.Errorf("user is not active")
	}

	if borrowedDays != 30 && borrowedDays != 60 && borrowedDays != 90 {
		return nil, fmt.Errorf("borrowed days must be 30, 60, or 90")
	}

	activeLoans, err := ru.userRepo.GetUserLoans(userId)
	if err != nil {
		return nil, fmt.Errorf("error when searching for user loans: %w", err)
	}

	borrowedLoansCount := 0
	for _, loan := range *activeLoans {
		if loan.Status == "borrowed" {
			borrowedLoansCount++
			// Verificar se o empréstimo está em atraso
			if loan.ReturnBy.Before(time.Now()) {
				return nil, fmt.Errorf("user has overdue loans")
			}
		}
	}

	activeReservations, err := ru.userRepo.GetUserReservations(userId)
	if err != nil {
		return nil, fmt.Errorf("error when searching for user reservations: %w", err)
	}

	pendingReservationsCount := 0
	for _, res := range *activeReservations {
		if res.Status == model.ReservationPending {
			pendingReservationsCount++
		}
	}

	totalActive := pendingReservationsCount + borrowedLoansCount
	if totalActive >= 5 {
		return nil, fmt.Errorf("user already has 5 or more active reservations/loans")
	}

	bUseCase := NewBookUseCase(ru.bookRepo, ru.reservationRepo)
	amount, err := bUseCase.CountAvailableBookStockById(bookId)
	if err != nil {
		return nil, fmt.Errorf("error when getting book stock amount: %w", err)
	}
	if amount <= 0 {
		return nil, fmt.Errorf("book out of stock")
	}

	return ru.reservationRepo.CreateReservation(borrowedDays, userId, bookId)
}

func (ru *reservationUseCase) GetReservationById(id int) (*model.Reservation, error) {
	return ru.reservationRepo.GetReservationById(id)
}
