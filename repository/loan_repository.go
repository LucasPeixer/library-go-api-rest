package repository

import (
	"database/sql"
	"fmt"
	"go-api/model"
)

type LoanRepository interface {
	CreateLoan(reservationId, bookStockId, borrowedDays int) (*model.Loan, error)
	GetLoanById(id int) (*model.Loan, error)
	UpdateLoan(loan *model.Loan) error
}

type loanRepository struct {
	db *sql.DB
}

func NewLoanRepository(db *sql.DB) LoanRepository {
	return &loanRepository{db: db}
}

func (lr *loanRepository) CreateLoan(reservationId, bookStockId, borrowedDays int) (*model.Loan, error) {
	query := `
	INSERT INTO loan (fk_reservation_id, fk_book_stock_id, return_by) 
	VALUES ($1, $2, CURRENT_TIMESTAMP + ($3 || ' days')::INTERVAL)
	RETURNING id, loaned_at, return_by
	`

	var loan model.Loan
	err := lr.db.QueryRow(query, reservationId, bookStockId, borrowedDays).Scan(&loan.Id, &loan.LoanedAt, &loan.ReturnBy)
	if err != nil {
		return nil, err
	}

	loan.BookStockId = bookStockId
	loan.Status = model.LoanBorrowed

	return &loan, nil
}

func (lr *loanRepository) GetLoanById(id int) (*model.Loan, error) {
	var loan model.Loan
	query := `SELECT * FROM loan WHERE id = $1`
	err := lr.db.QueryRow(query, id).Scan(
		&loan.Id,
		&loan.LoanedAt,
		&loan.ReturnBy,
		&loan.ReturnedAt,
		&loan.Status,
		&loan.AdminId,
		&loan.BookStockId,
		&loan.ReservationId,
	)
	if err != nil {
		return nil, err
	}
	return &loan, nil
}

func (lr *loanRepository) UpdateLoan(loan *model.Loan) error {
	query := `UPDATE loan 
              SET returned_at = $1, fk_admin_id = $2, status = $3 
              WHERE id = $4`
	_, err := lr.db.Exec(query, loan.ReturnedAt, loan.AdminId, loan.Status, loan.Id)
	if err != nil {
		return fmt.Errorf("failed to update loan: %w", err)
	}
	return nil
}
