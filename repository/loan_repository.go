package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"go-api/model"
)

type LoanRepository interface {

}

type loanRepository struct {
	db *sql.DB
}

func NewLoanRepository(db *sql.DB) LoanRepository {
	return &loanRepository{db: db}
}

func (lr *LoanRepository) CreateLoan(loanRequest *model.LoanRequest) (*model.Loan, error) {
	query := `
			INSERT INTO loan (return_by, fk_book_stock_id, fk_reservation_id)
			VALUES ($1, $2, $3)
			RETURNING id, loaned_at, return_by, returned_at, status, fk_admin_id, fk_book_stock_id, fk_reservation_id`
	
	var loan model.Loan
	err := lr.db.QueryRow(query, loanRequest.ReturnBy, loanRequest.AdminID, loanRequest.BookStockID, loanRequest.ReservationID).
			Scan(&loan.ID, &loan.LoanedAt, &loan.ReturnBy, &loan.ReturnedAt, &loan.Status, &loan.AdminID, &loan.BookStockID, &loan.ReservationID)
	if err != nil {
			return nil, fmt.Errorf("error inserting loan: %w", err)
	}

	return &loan, nil
}