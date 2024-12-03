package repository

import (
	"database/sql"
	"fmt"
	"go-api/model"
	"go-api/model/user"
	"strconv"
)

type LoanRepository interface {
	CreateLoan(reservationId, bookStockId, borrowedDays int) (*model.Loan, error)
	GetLoansByFilters(userName string, status model.LoanStatus, loanedAt string) (*[]model.Loan, error)
	GetLoanById(id int) (*model.Loan, error)
	FinishLoan(id, adminId int) error
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

	loan.Status = model.LoanBorrowed
	loan.ReservationId = reservationId

	return &loan, nil
}

func (lr *loanRepository) GetLoansByFilters(userName string, status model.LoanStatus, loanedAt string) (*[]model.Loan, error) {
	query := `
	SELECT 
	    l.id                AS loan_id,
	    l.loaned_at,
	    l.return_by,
	    l.returned_at,
	    l.status            AS loan_status,
	    u.id                AS user_account_id,
	    u.name              AS user_account_name,
	    a.id                AS admin_account_id,
	    a.name              AS admin_account_name,
	    bs.id               AS book_stock_id,
	    bs.code             AS book_stock_code,
	    l.fk_reservation_id AS reservation_id
	FROM 
	    loan l
	LEFT JOIN
	    reservation r ON l.fk_reservation_id = r.id
	LEFT JOIN
	    user_account u ON r.fk_user_id = u.id
	LEFT JOIN
	    user_account a ON l.fk_admin_id = a.id
	JOIN
		 book_stock bs ON l.fk_book_stock_id = bs.id
	WHERE 
		 1=1`

	var args []interface{}

	if userName != "" {
		query += ` AND u.name ILIKE $` + strconv.Itoa(len(args)+1)
		args = append(args, "%"+userName+"%")
	}

	if status != "" {
		query += ` AND l.status = $` + strconv.Itoa(len(args)+1)
		args = append(args, string(status))
	}

	if loanedAt != "" {
		query += ` AND l.loaned_at::date = $` + strconv.Itoa(len(args)+1)
		args = append(args, loanedAt)
	}
	rows, err := lr.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	loans := make([]model.Loan, 0)
	for rows.Next() {
		var loan model.Loan
		loan.UserAccount = &user.Account{}
		loan.BookStock = &model.BookStock{}
		var adminId *int
		var adminName *string

		if err := rows.Scan(
			&loan.Id,
			&loan.LoanedAt,
			&loan.ReturnBy,
			&loan.ReturnedAt,
			&loan.Status,
			&loan.UserAccount.Id,
			&loan.UserAccount.Name,
			&adminId,
			&adminName,
			&loan.BookStock.Id,
			&loan.BookStock.Code,
			&loan.ReservationId,
		); err != nil {
			return nil, err
		}

		if adminId != nil {
			loan.AdminAccount = &user.Account{Id: *adminId, Name: *adminName}
		}

		loans = append(loans, loan)
	}
	return &loans, nil
}

func (lr *loanRepository) GetLoanById(id int) (*model.Loan, error) {
	query := `
	SELECT 
	    l.id                AS loan_id,
	    l.loaned_at,
	    l.return_by,
	    l.returned_at,
	    l.status            AS loan_status,
	    u.id                AS user_account_id,
	    u.name              AS user_account_name,
	    a.id                AS admin_account_id,
	    a.name              AS admin_account_name,
	    bs.id               AS book_stock_id,
	    bs.code             AS book_stock_code,
	    l.fk_reservation_id AS reservation_id
	FROM 
	    loan l
	LEFT JOIN
	    reservation r ON l.fk_reservation_id = r.id
	LEFT JOIN
	    user_account u ON r.fk_user_id = u.id
	LEFT JOIN
	    user_account a ON l.fk_admin_id = a.id
	JOIN
		 book_stock bs ON l.fk_book_stock_id = bs.id
	WHERE 
		 l.id = $1;`

	var loan model.Loan
	loan.UserAccount = &user.Account{}
	loan.BookStock = &model.BookStock{}
	var adminId *int
	var adminName *string

	err := lr.db.QueryRow(query, id).Scan(
		&loan.Id,
		&loan.LoanedAt,
		&loan.ReturnBy,
		&loan.ReturnedAt,
		&loan.Status,
		&loan.UserAccount.Id,
		&loan.UserAccount.Name,
		&adminId,
		&adminName,
		&loan.BookStock.Id,
		&loan.BookStock.Code,
		&loan.ReservationId,
	)
	if err != nil {
		return nil, err
	}

	if adminId != nil {
		loan.AdminAccount = &user.Account{Id: *adminId, Name: *adminName}
	}

	return &loan, nil
}

func (lr *loanRepository) FinishLoan(id, adminId int) error {
	query := `UPDATE loan SET returned_at = CURRENT_TIMESTAMP, status = 'returned', fk_admin_id = $1 WHERE id = $2`
	_, err := lr.db.Exec(query, adminId, id)
	if err != nil {
		return fmt.Errorf("failed to finish loan: %w", err)
	}
	return nil
}
