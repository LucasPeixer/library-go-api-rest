package usecase

import (
	"fmt"
	"go-api/model"
	"go-api/repository"
	"time"
)

type LoanUseCaseInterface interface {
	CreateLoanAndUpdateReservation(request *model.LoanRequest) (*model.Loan, error)
	UpdateLoan(request model.LoanUpdateRequest, adminID int, loanId int) error
}

type LoanUseCase struct {
	loanRepo        repository.LoanRepositoryInterface
	bookStockRepo   repository.BookRepository
	reservationRepo repository.ReservationRepositoryInterface
}


func NewLoanUseCase(
	loanRepo repository.LoanRepositoryInterface,
	reservationRepo repository.ReservationRepositoryInterface,
	bookStockRepo repository.BookRepository,
) *LoanUseCase {
	return &LoanUseCase{
			loanRepo:        loanRepo,
			bookStockRepo:   bookStockRepo,
			reservationRepo: reservationRepo,
	}
}


func (lu *LoanUseCase) CreateLoanAndUpdateReservation(request *model.LoanRequest) (*model.Loan, error) {

	reservation, err := lu.reservationRepo.GetReservationByID(request.ReservationID)
	if err != nil {
			return nil, fmt.Errorf("error fetching reservation: %w", err)
	}
	
	if reservation.Status != "pending" {
			return nil, fmt.Errorf("reservation is not pending")
	}

	expiryBuffer := time.Now().Add(-30 * time.Minute)
	if reservation.ExpiresAt.Before(expiryBuffer) {
    	return nil, fmt.Errorf("reservation has expired")
	}
	
	bookStock, err := lu.bookStockRepo.GetStockById(request.BookStockID)
	if err != nil {
			return nil, fmt.Errorf("error fetching book stock: %w", err)
	}

	if bookStock.Status != "available" {
			return nil, fmt.Errorf("book stock is not available")
	}

	err = lu.reservationRepo.UpdateReservationStatus(request.ReservationID, "collected",*request.AdminID)
	if err != nil {
		return nil, fmt.Errorf("failed to update reservation status: %w", err)
	}

	err = lu.bookStockRepo.UpdateStockStatus(request.BookStockID, "borrowed")
	if err != nil {
		return nil, fmt.Errorf("failed to update book stock status: %w", err)
	}
	
	createdLoan, err := lu.loanRepo.CreateLoan(request)
	if err != nil {
			return nil, fmt.Errorf("error creating loan: %w", err)
	}

	return createdLoan, nil
}

func (lu *LoanUseCase) UpdateLoan(request model.LoanUpdateRequest, adminID int, loanId int) error {

	loan, err := lu.loanRepo.GetLoanByID(loanId)
	if err != nil {
		return fmt.Errorf("loan not found: %w", err)
	}

	// Verifica o status
	if loan.Status != "borrowed" {
		return fmt.Errorf("loan is not in borrowed status")
	}

	// Atualiza os dados
	now := time.Now()
	loan.ReturnedAt = &now
	loan.AdminID = &adminID
	loan.Status = "returned"

	// Atualiza no repositório
	err = lu.loanRepo.UpdateLoan(loan)
	if err != nil {
		return fmt.Errorf("failed to update loan: %w", err)
	}

	return nil
}