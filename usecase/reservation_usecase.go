package usecase

import (
	"fmt"
	"go-api/repository"
	"go-api/model"
	"time"
)

type ReservationUseCaseInterface interface {
	GetReservationsByFilters(userName, status, reservedAt string) ([]model.Reservation, error)
	CreateReservation(reservation *model.ReservationRequest) (*model.Reservation, error)
	GetReservationByID(reservationID int) (*model.Reservation, error)
}

type ReservationUseCase struct {
	ReservationRepo repository.ReservationRepositoryInterface
	userRepo repository.UserRepository
	bookRepo repository.BookRepository
}

// NewReservationUseCase cria e retorna uma nova instância de ReservationUseCase
func NewReservationUseCase(reservationRepo repository.ReservationRepositoryInterface, userRepo repository.UserRepository, bookRepo repository.BookRepository) ReservationUseCaseInterface {
	return &ReservationUseCase{
		ReservationRepo: reservationRepo,
		userRepo:        userRepo,
		bookRepo:        bookRepo,
	}
}

func (ru *ReservationUseCase) GetReservationsByFilters(userName, status, reservedAt string) ([]model.Reservation, error) {
	return ru.ReservationRepo.GetReservationsByFilters(userName, status, reservedAt)
}

func (ru *ReservationUseCase) CreateReservation(reservation *model.ReservationRequest) (*model.Reservation, error) {

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

	book, err := ru.bookRepo.GetBookById(reservation.BookID)
	if err != nil {
		return nil, fmt.Errorf("error when searching for book: %w", err)
	}
	if book.Amount <= 0 {
		return nil, fmt.Errorf("book out of stock")
	}

	newReservation, err := ru.ReservationRepo.CreateReservation(reservation)
	if err != nil {
		return nil, fmt.Errorf("error when creating reservation: %w", err)
	}

	return newReservation, nil
}

func (ru *ReservationUseCase) GetReservationByID(reservationID int) (*model.Reservation, error) {
	return ru.ReservationRepo.GetReservationByID(reservationID)
}
