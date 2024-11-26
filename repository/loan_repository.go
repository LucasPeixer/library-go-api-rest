package repository

import (
	"database/sql"
	"fmt"
	"go-api/model"
	"strings"
)

type LoanRepositoryInterface interface {
	CreateLoan(loan *model.LoanRequest) (*model.Loan, error)
	GetLoanByID(id int) (*model.Loan, error)
	UpdateLoan(loan *model.Loan) error
	GetLoansByFilters(loanedAt, returnBy, returnedAt, status string, bookStockID, reservationID int) ([]model.Loan, error)
}

type loanRepository struct {
	db *sql.DB
}

func NewLoanRepository(db *sql.DB) LoanRepositoryInterface {
	return &loanRepository{db: db}
}

func (lr *loanRepository) GetLoansByFilters(loanedAt, returnBy, returnedAt, status string, bookStockID, reservationID int) ([]model.Loan, error) {

	var sb strings.Builder
	sb.WriteString("SELECT * FROM loan WHERE 1=1")

	var args []interface{}
	paramIndex := 1

	if loanedAt != "" {
			sb.WriteString(" AND loaned_at = $" + fmt.Sprint(paramIndex))
			args = append(args, loanedAt)
			paramIndex++
	}

	if returnBy != "" {
			sb.WriteString(" AND return_by = $" + fmt.Sprint(paramIndex))
			args = append(args, returnBy)
			paramIndex++
	}

	if returnedAt != "" {
			sb.WriteString(" AND returned_at = $" + fmt.Sprint(paramIndex))
			args = append(args, returnedAt)
			paramIndex++
	}

	if status != "" {
			sb.WriteString(" AND status = $" + fmt.Sprint(paramIndex))
			args = append(args, status)
			paramIndex++
	}

	if bookStockID > 0 {
			sb.WriteString(" AND fk_book_stock_id = $" + fmt.Sprint(paramIndex))
			args = append(args, bookStockID)
			paramIndex++
	}

	if reservationID > 0 {
			sb.WriteString(" AND fk_reservation_id = $" + fmt.Sprint(paramIndex))
			args = append(args, reservationID)
			paramIndex++
	}

	// Executando a consulta com os filtros aplicados
	query := sb.String()
	rows, err := lr.db.Query(query, args...)
	if err != nil {
			return nil, fmt.Errorf("error executing query: %w", err)
	}
	defer rows.Close()

	// Processando os resultados
	loans := make([]model.Loan, 0)
	for rows.Next() {
		var loan model.Loan
		if err := rows.Scan(&loan.ID, &loan.LoanedAt, &loan.ReturnBy, &loan.ReturnedAt, &loan.Status,&loan.AdminID, &loan.BookStockID,&loan.ReservationID); err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
		loans = append(loans, loan)
	}

	if err := rows.Err(); err != nil {
			return nil, fmt.Errorf("row iteration error: %w", err)
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