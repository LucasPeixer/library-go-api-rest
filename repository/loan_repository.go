package repository

import (
	"database/sql"
	"fmt"
	"go-api/model"
)

type LoanRepositoryInterface interface {
	CreateLoan(loan *model.LoanRequest) (*model.Loan, error)
}

type loanRepository struct {
	db *sql.DB
}

func NewLoanRepository(db *sql.DB) LoanRepositoryInterface {
	return &loanRepository{db: db}
}

func (lr *loanRepository) CreateLoan(loan *model.LoanRequest) (*model.Loan, error) {
	query := `
			INSERT INTO loan (return_by, fk_book_stock_id, fk_reservation_id)
			VALUES ($1, $2, $3)
			RETURNING id, loaned_at, return_by, returned_at, status, fk_admin_id, fk_book_stock_id, fk_reservation_id`

	var createdLoan model.Loan
	err := lr.db.QueryRow(
			query,
			loan.ReturnBy,
			loan.BookStockID,
			loan.ReservationID,
	).Scan(
			&createdLoan.ID,
			&createdLoan.LoanedAt,
			&createdLoan.ReturnBy,
			&createdLoan.ReturnedAt,
			&createdLoan.Status,
			&createdLoan.AdminID,
			&createdLoan.BookStockID,
			&createdLoan.ReservationID,
	)
	if err != nil {
			return nil, fmt.Errorf("error inserting loan: %w", err)
	}

	return &createdLoan, nil
}