package model

import "time"

type LoanStatus string

const (
	LoanBorrowed LoanStatus = "borrowed"
	LoanReturned LoanStatus = "returned"
)

type Loan struct {
	Id            int        `json:"id" db:"id"`
	LoanedAt      time.Time  `json:"loaned_at" db:"loaned_at"`
	ReturnBy      time.Time  `json:"return_by" db:"return_by"`
	ReturnedAt    *time.Time `json:"returned_at,omitempty" db:"returned_at"`
	Status        LoanStatus `json:"status" db:"status"`
	AdminId       *int       `json:"admin_id,omitempty" db:"fk_admin_id"`
	BookStockId   int        `json:"book_stock_id" db:"fk_book_stock_id"`
	ReservationId int        `json:"reservation_id" db:"fk_reservation_id"`
}

type LoanRequest struct {
	ReturnBy      time.Time `json:"return_by"`
	BookStockID   int       `json:"book_stock_id" binding:"required"`
	ReservationID int       `json:"reservation_id" binding:"required"`
	AdminID       *int      `json:"admin_id"`
}

type LoanUpdateRequest struct {
	ID         int        `json:"id"`
	AdminId    int        `json:"admin_id"`
	Status     LoanStatus `json:"status" binding:"required" db:"status"`
	ReturnedAt *time.Time `json:"returned_at"`
}
