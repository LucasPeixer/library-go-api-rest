package model

import "time"

type ReservationStatus string

const (
	Pending  ReservationStatus = "pending"
	Cancelled ReservationStatus = "cancelled"
	Expired   ReservationStatus = "expired"
	Collected ReservationStatus = "collected"
	Finished ReservationStatus = "finished"
)

type Reservation struct {
	ID           int    `json:"id" db:"id"`
	ReservedAt   time.Time `json:"reserved_at" db:"reserved_at"`
	ExpiresAt    time.Time `json:"expires_at" db:"expires_at"`
	BorrowedDays int       `json:"borrowed_days" db:"borrowed_days"`
	Status       ReservationStatus    `json:"status" db:"status"`
	UserID       int    `json:"fk_user_id" db:"fk_user_id"`
	AdminID      int    `json:"fk_admin_id" db:"fk_admin_id"`
	BookID       int    `json:"fk_book_id" db:"fk_book_id"`
}