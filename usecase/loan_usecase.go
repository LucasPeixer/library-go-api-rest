package usecase

import (
	"fmt"
	"go-api/repository"
	"go-api/model"
	"time"
)

type LoanUseCaseInterface interface {
	CreateLoan(request *model.LoanRequest) (*model.Loan, error)
}

type LoanUseCase struct {
	LoanRepo repository.LoanRepositoryInterface
}

func NewLoanUseCase(loanRepo repository.LoanRepositoryInterface, reservationRepo repository.ReservationRepositoryInterface, bookStockRepo repository.BookStockRepositoryInterface){
	return &LoanUseCase{
		LoanRepo: loanRepo,
		bookStockRepo: bookStockRepo,
		reservationRepo: reservationRepo
	}
}

func (lu *LoanUseCase) CreateLoan(request *model.LoanRequest) (*model.Loan, error) {

	/*reservation, err := lu.reservationRepo.GetReservationByID(request.ReservationID)
	if err != nil {
			return nil, fmt.Errorf("error fetching reservation: %w", err)
	}

	if reservation.Status != "collected" {
			return nil, fmt.Errorf("reservation is not collected")
	}*/

	/*bookStock, err := lu.bookStockRepo.GetBookStockByID(request.BookStockID)
	if err != nil {
			return nil, fmt.Errorf("error fetching book stock: %w", err)
	}

	// Ensure book stock status is "available"
	if bookStock.Status != "available" {
			return nil, fmt.Errorf("book stock is not available")
	}*/

	// Calculate return_by using reservation's BorrowedDays and current timestamp
	/*returnBy := time.Now().AddDate(0, 0, reservation.BorrowedDays)

	loan := &model.Loan{
			ReturnBy:      returnBy,
			BookStockID:   request.BookStockID,
			ReservationID: request.ReservationID,
	}*/

	// Create loan
	createdLoan, err := lu.loanRepo.CreateLoan(loan)
	if err != nil {
			return nil, fmt.Errorf("error creating loan: %w", err)
	}

	return createdLoan, nil
}
