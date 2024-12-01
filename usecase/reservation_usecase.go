package usecase

import (
	"fmt"
	"go-api/model"
	"go-api/repository"
	"time"
)

type ReservationUseCase interface {
	GetReservationsByFilters(userName, status, reservedAt string) ([]model.Reservation, error)
	CreateReservation(reservation *model.ReservationRequest) (*model.Reservation, error)
	GetReservationByID(reservationID int) (*model.Reservation, error)
}

type reservationUseCase struct {
	reservationRepo repository.ReservationRepository
	userRepo        repository.UserRepository
	bookRepo        repository.BookRepository
}

// NewReservationUseCase cria e retorna uma nova instância de ReservationUseCase
func NewReservationUseCase(reservationRepo repository.ReservationRepository,
	userRepo repository.UserRepository, bookRepo repository.BookRepository) ReservationUseCase {
	return &reservationUseCase{
		reservationRepo: reservationRepo,
		userRepo:        userRepo,
		bookRepo:        bookRepo,
	}
}

func (ru *reservationUseCase) GetReservationsByFilters(userName, status, reservedAt string) ([]model.Reservation, error) {
	return ru.reservationRepo.GetReservationsByFilters(userName, status, reservedAt)
}

func (ru *reservationUseCase) CreateReservation(reservation *model.ReservationRequest) (*model.Reservation, error) {

	user, err := ru.userRepo.GetUserById(reservation.UserID)
	if err != nil {
		return nil, fmt.Errorf("error when searching for user: %w", err)
	}
	if user.IsActive != true {
		return nil, fmt.Errorf("user is not active")
	}

	if reservation.BorrowedDays != 30 && reservation.BorrowedDays != 60 && reservation.BorrowedDays != 90 {
		return nil, fmt.Errorf("borrowed days must be 30, 60, or 90")
	}

	activeLoans, err := ru.userRepo.GetUserLoans(reservation.UserID)
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

	userReservations, err := ru.userRepo.GetUserReservations(reservation.UserID)
	if err != nil {
		return nil, fmt.Errorf("error when searching for user reservations: %w", err)
	}

	pendingReservationsCount := 0
	for _, res := range *userReservations {
		if res.Status == "pending" {
			pendingReservationsCount++
		}
	}

	totalActive := pendingReservationsCount + borrowedLoansCount
	if totalActive >= 5 {
		return nil, fmt.Errorf("user already has 5 or more active reservations/loans")
	}

	bUseCase := NewBookUseCase(ru.bookRepo, ru.reservationRepo)
	amount, err := bUseCase.CountAvailableBookStockById(reservation.BookID)
	if err != nil {
		return nil, fmt.Errorf("error when searching for book: %w", err)
	}
	if amount <= 0 {
		return nil, fmt.Errorf("book out of stock")
	}

	newReservation, err := ru.reservationRepo.CreateReservation(reservation)
	if err != nil {
		return nil, fmt.Errorf("error when creating reservation: %w", err)
	}

	return newReservation, nil
}

func (ru *reservationUseCase) GetReservationByID(reservationID int) (*model.Reservation, error) {
	return ru.reservationRepo.GetReservationByID(reservationID)
}
