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
	GetLoansByFilter(filters map[string]interface{}) ([]model.Loan, error)
}

type loanRepository struct {
	db *sql.DB
}

func NewLoanRepository(db *sql.DB) LoanRepositoryInterface {
	return &loanRepository{db: db}
}

func (lr *loanRepository) GetLoansByFilter(filters map[string]interface{}) ([]model.Loan, error) {
	// Base da query
	query := `SELECT id, loaned_at, return_by, returned_at, status, fk_book_stock_id, fk_reservation_id FROM loans`
	var conditions []string
	var args []interface{}
	argIndex := 1

	// Constrói as condições com base nos filtros
	for key, value := range filters {
		conditions = append(conditions, fmt.Sprintf("%s = $%d", key, argIndex))
		args = append(args, value)
		argIndex++
	}

	// Adiciona condições à query se existirem filtros
	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	// Executa a query
	rows, err := lr.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	// Mapeia os resultados
	var loans []model.Loan
	for rows.Next() {
		var loan model.Loan
		err := rows.Scan(
			&loan.ID,
			&loan.LoanedAt,
			&loan.ReturnBy,
			&loan.ReturnedAt,
			&loan.Status,
			&loan.BookStockID,
			&loan.ReservationID,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		loans = append(loans, loan)
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