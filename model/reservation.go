package model

import (
	"go-api/model/user"
	"time"
)

type ReservationStatus string

const (
	ReservationPending   ReservationStatus = "pending"
	ReservationCancelled ReservationStatus = "cancelled"
	ReservationExpired   ReservationStatus = "expired"
	ReservationCollected ReservationStatus = "collected"
	ReservationFinished  ReservationStatus = "finished"
)

type Reservation struct {
	Id           int               `json:"id" db:"id"`
	ReservedAt   time.Time         `json:"reserved_at" db:"reserved_at"`
	ExpiresAt    time.Time         `json:"expires_at" db:"expires_at"`
	BorrowedDays int               `json:"borrowed_days" db:"borrowed_days"`
	Status       ReservationStatus `json:"status" db:"status"`
	UserAccount  *user.Account     `json:"user_account,omitempty" db:"user_account"`
	AdminAccount *user.Account     `json:"admin_account,omitempty" db:"user_account"`
	Book         Book              `json:"book" db:"book"`
}
