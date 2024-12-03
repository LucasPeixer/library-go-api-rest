package usecase

import (
	"fmt"
	"go-api/model"
	"go-api/repository"
)

type LoanUseCase interface {
	CreateLoanAndUpdateReservation(reservationId, bookStockId, adminId int) (*model.Loan, error)
	GetLoansByFilters(userName, status model.LoanStatus, loanedAt string) (*[]model.Loan, error)
	GetLoanById(id int) (*model.Loan, error)
	GetLoanByReservationId(reservationId int) (*model.Loan, error)
	FinishLoan(loanId, adminId int) error
}

type loanUseCase struct {
	loanRepo        repository.LoanRepository
	bookRepo        repository.BookRepository
	reservationRepo repository.ReservationRepository
}

func NewLoanUseCase(
	loanRepo repository.LoanRepository,
	reservationRepo repository.ReservationRepository,
	bookStockRepo repository.BookRepository) LoanUseCase {
	return &loanUseCase{
		loanRepo:        loanRepo,
		bookRepo:        bookStockRepo,
		reservationRepo: reservationRepo,
	}
}

func (lu *loanUseCase) CreateLoanAndUpdateReservation(reservationId, bookStockId, adminId int) (*model.Loan, error) {
	reservation, err := lu.reservationRepo.GetReservationById(reservationId)
	if err != nil {
		return nil, fmt.Errorf("error fetching reservation: %w", err)
	}

	if reservation.Status == model.ReservationExpired {
		return nil, fmt.Errorf("reservation has expired")
	}

	if reservation.Status != model.ReservationPending {
		return nil, fmt.Errorf("reservation is not pending")
	}

	bookStock, err := lu.bookRepo.GetStockById(bookStockId)
	if err != nil {
		return nil, fmt.Errorf("error fetching book stock: %w", err)
	}

	if bookStock.Status != model.BookStockAvailable {
		return nil, fmt.Errorf("book stock is not available")
	}

	err = lu.reservationRepo.UpdateReservationStatus(reservationId, "collected", adminId)
	if err != nil {
		return nil, fmt.Errorf("failed to update reservation status: %w", err)
	}

	err = lu.bookRepo.UpdateStockStatus(bookStockId, "borrowed")
	if err != nil {
		return nil, fmt.Errorf("failed to update book stock status: %w", err)
	}

	createdLoan, err := lu.loanRepo.CreateLoan(reservationId, bookStockId, reservation.BorrowedDays)
	if err != nil {
		return nil, fmt.Errorf("error creating loan: %w", err)
	}

	return createdLoan, nil
}

func (lu *loanUseCase) GetLoansByFilters(userName, status model.LoanStatus, loanedAt string) (*[]model.Loan, error) {
	//TODO implement me
	panic("implement me")
}

func (lu *loanUseCase) GetLoanById(id int) (*model.Loan, error) {
	//TODO implement me
	panic("implement me")
}

func (lu *loanUseCase) GetLoanByReservationId(reservationId int) (*model.Loan, error) {
	//TODO implement me
	panic("implement me")
}

func (lu *loanUseCase) FinishLoan(loanId, adminId int) error {
	loan, err := lu.loanRepo.GetLoanById(loanId)
	if err != nil {
		return err
	}

	if loan.Status != model.LoanBorrowed {
		return fmt.Errorf("loan is not borrowed")
	}

	err = lu.loanRepo.FinishLoan(loanId, adminId)
	if err != nil {
		return err
	}

	return nil
}
