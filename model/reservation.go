package model

import "time"

type ReservationStatus string

const (
	ReservationPending   ReservationStatus = "pending"
	ReservationCancelled ReservationStatus = "cancelled"
	ReservationExpired   ReservationStatus = "expired"
	ReservationCollected ReservationStatus = "collected"
	ReservationFinished  ReservationStatus = "finished"
)

type Reservation struct {
	ID           int               `json:"id" db:"id"`
	ReservedAt   time.Time         `json:"reserved_at" db:"reserved_at"`
	ExpiresAt    time.Time         `json:"expires_at" db:"expires_at"`
	BorrowedDays int               `json:"borrowed_days" db:"borrowed_days"`
	Status       ReservationStatus `json:"status" db:"status"`
	UserID       int               `json:"fk_user_id" db:"fk_user_id"`
	AdminID      *int              `json:"fk_admin_id" db:"fk_admin_id"`
	BookID       int               `json:"fk_book_id" db:"fk_book_id"`
}

type ReservationRequest struct {
	UserID       int `json:"user_id" binding:"required"`
	BookID       int `json:"book_id" binding:"required"`
	BorrowedDays int `json:"borrowed_days" binding:"required"`
}
