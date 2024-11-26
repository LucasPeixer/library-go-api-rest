package repository

import (
	"database/sql"
	"fmt"
	"go-api/model"
)

type LoanRepositoryInterface interface {
	CreateLoan(loan *model.LoanRequest) (*model.Loan, error)
	GetLoanByID(id int) (*model.Loan, error)
	UpdateLoan(loan *model.Loan) error
	GetLoansByFilters(filters map[string]interface{}) ([]model.Loan, error)
}

type loanRepository struct {
	db *sql.DB
}

func NewLoanRepository(db *sql.DB) LoanRepositoryInterface {
	return &loanRepository{db: db}
}

func (lr *loanRepository) GetLoansByFilters(filters map[string]interface{}) ([]model.Loan, error) {
    query := "SELECT * FROM loans WHERE 1=1"
    args := []interface{}{}

    if loanedAt, ok := filters["loaned_at"].(string); ok && loanedAt != "" {
        query += " AND loaned_at = ?"
        args = append(args, loanedAt)
    }
    if returnBy, ok := filters["return_by"].(string); ok && returnBy != "" {
        query += " AND return_by = ?"
        args = append(args, returnBy)
    }
    if returnedAt, ok := filters["returned_at"].(string); ok && returnedAt != "" {
        query += " AND returned_at = ?"
        args = append(args, returnedAt)
    }
    if status, ok := filters["status"].(string); ok && status != "" {
        query += " AND status = ?"
        args = append(args, status)
    }
    if bookStockID, ok := filters["book_stock_id"].(int); ok && bookStockID > 0 {
        query += " AND fk_book_stock_id = ?"
        args = append(args, bookStockID)
    }
    if reservationID, ok := filters["reservation_id"].(int); ok && reservationID > 0 {
        query += " AND fk_reservation_id = ?"
        args = append(args, reservationID)
    }

    rows, err := lr.db.Query(query, args...)
    if err != nil {
        return nil, fmt.Errorf("error executing query: %w", err)
    }
    defer rows.Close()

    var loans []model.Loan
    for rows.Next() {
        var loan model.Loan
        if err := rows.Scan(&loan.ID, &loan.LoanedAt, &loan.ReturnBy, &loan.ReturnedAt, &loan.Status, &loan.BookStockID, &loan.ReservationID); err != nil {
            return nil, fmt.Errorf("error scanning row: %w", err)
        }
        loans = append(loans, loan)
    }

    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("error iterating rows: %w", err)
    }

    return loans, nil
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

func (lr *loanRepository) GetLoanByID(id int) (*model.Loan, error) {
	var loan model.Loan
	query := `SELECT * FROM loan WHERE id = $1`
	err := lr.db.QueryRow(query, id).Scan(
		&loan.ID,
		&loan.LoanedAt,
		&loan.ReturnBy,
		&loan.ReturnedAt,
		&loan.Status,
		&loan.AdminID,
		&loan.BookStockID,
		&loan.ReservationID,
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
	_, err := lr.db.Exec(query, loan.ReturnedAt, loan.AdminID, loan.Status, loan.ID)
	if err != nil {
		return fmt.Errorf("failed to update loan: %w", err)
	}
	return nil
}