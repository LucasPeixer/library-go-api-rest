package model

import "time"

type Reservation struct {
	ID           string    `json:"id" db:"id"`
	ReservedAt   time.Time `json:"reserved_at" db:"reserved_at"`
	ExpiresAt    time.Time `json:"expires_at" db:"expires_at"`
	BorrowedDays int       `json:"borrowed_days" db:"borrowed_days"`
	Status       string    `json:"status" db:"status"`
	UserID       string    `json:"fk_user_id" db:"fk_user_id"`
	AdminID      string    `json:"fk_admin_id" db:"fk_admin_id"`
	BookID       string    `json:"fk_book_id" db:"fk_book_id"`
}