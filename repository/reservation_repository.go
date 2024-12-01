package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"go-api/model"
	"go-api/model/user"
	"strconv"
)

type ReservationRepository interface {
	CreateReservation(borrowedDays, userId, bookId int) (*model.Reservation, error)
	GetReservationsByFilters(userName string, status model.ReservationStatus, reservedAt string) (*[]model.Reservation, error)
	GetReservationsByBookId(id int, status string) (*[]model.Reservation, error)
	GetReservationById(id int) (*model.Reservation, error)
	UpdateReservationStatus(reservationID int, status string, adminID int) error
	CancelReservation(id, adminId int) error
}

type reservationRepository struct {
	db *sql.DB
}

func NewReservationRepository(db *sql.DB) ReservationRepository {
	return &reservationRepository{db}
}

func (rr *reservationRepository) CreateReservation(borrowedDays, userId, bookId int) (*model.Reservation, error) {
	query := `
		INSERT INTO reservation (borrowed_days, fk_user_id, fk_book_id)
		VALUES ($1, $2, $3)
		RETURNING id, reserved_at, expires_at`

	var res model.Reservation
	err := rr.db.QueryRow(query, borrowedDays, userId, bookId).Scan(&res.Id, &res.ReservedAt, &res.ExpiresAt)
	if err != nil {
		return nil, err
	}
	res.Status = model.ReservationPending
	res.BorrowedDays = borrowedDays
	return &res, nil
}

func (rr *reservationRepository) GetReservationsByFilters(userName string, status model.ReservationStatus, reservedAt string) (*[]model.Reservation, error) {
	query := `
	SELECT r.id            AS reservation_id  ,
	       r.reserved_at,
	       r.expires_at,
	       r.borrowed_days,
	       r.status        AS reservation_status,
	       r.fk_user_id    AS user_id,
	       usr.name        AS user_name,
	       r.fk_admin_id   AS admin_id,
	       adm.name        AS admin_name,
	       r.fk_book_id    AS book_id,
	       b.title         AS book_title,
		   (CURRENT_TIMESTAMP > r.expires_at) as is_expired
	FROM 
	       reservation r
	JOIN 
	       user_account usr ON r.fk_user_id = usr.id
	LEFT JOIN 
	       user_account adm ON r.fk_admin_id = adm.id
	JOIN 
	       book b ON r.fk_book_id = b.id
	WHERE 
		    1=1 -- Permite adicionar condições "AND"
   `
	var args []interface{}

	if userName != "" {
		query += ` AND usr.name ILIKE $` + strconv.Itoa(len(args)+1)
		args = append(args, "%"+userName+"%")
	}

	if status != "" {
		query += ` AND r.status = $` + strconv.Itoa(len(args)+1)
		args = append(args, string(status))
	}

	if reservedAt != "" {
		query += ` AND r.reserved_at::date = $` + strconv.Itoa(len(args)+1)
		args = append(args, reservedAt)
	}
	rows, err := rr.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("error fetching reservations: %w", err)
	}
	defer rows.Close()

	reservations := make([]model.Reservation, 0)
	for rows.Next() {
		var res model.Reservation
		res.UserAccount = &user.Account{}
		res.AdminAccount = &user.Account{}
		var adminId *int
		var adminName *string
		var isExpired bool

		if err := rows.Scan(
			&res.Id,
			&res.ReservedAt,
			&res.ExpiresAt,
			&res.BorrowedDays,
			&res.Status,
			&res.UserAccount.Id,
			&res.UserAccount.Name,
			&adminId,
			&adminName,
			&res.Book.Id,
			&res.Book.Title,
			&isExpired,
		); err != nil {
			return nil, err
		}

		if adminId != nil {
			res.AdminAccount.Id = *adminId
			res.AdminAccount.Name = *adminName
		} else {
			res.AdminAccount = nil
		}

		if isExpired {
			res.Status = model.ReservationExpired
		}

		reservations = append(reservations, res)
	}
	return &reservations, nil
}

func (rr *reservationRepository) GetReservationById(id int) (*model.Reservation, error) {
	query := `
	SELECT r.id            AS reservation_id  ,
	       r.reserved_at,
	       r.expires_at,
	       r.borrowed_days,
	       r.status        AS reservation_status,
	       r.fk_user_id    AS user_id,
	       usr.name        AS user_name,
	       r.fk_admin_id   AS admin_id,
	       adm.name        AS admin_name,
	       r.fk_book_id    AS book_id,
	       b.title         AS book_title,
		   (CURRENT_TIMESTAMP > r.expires_at) as is_expired
	FROM 
	       reservation r
	JOIN 
	       user_account usr ON r.fk_user_id = usr.id
	LEFT JOIN 
	       user_account adm ON r.fk_admin_id = adm.id
	JOIN 
	       book b ON r.fk_book_id = b.id
	WHERE
	       r.id = $1
   `

	res := model.Reservation{}
	res.UserAccount = &user.Account{}
	res.AdminAccount = &user.Account{}
	var adminId *int
	var adminName *string
	var isExpired bool

	err := rr.db.QueryRow(query, id).Scan(
		&res.Id,
		&res.ReservedAt,
		&res.ExpiresAt,
		&res.BorrowedDays,
		&res.Status,
		&res.UserAccount.Id,
		&res.UserAccount.Name,
		&adminId,
		&adminName,
		&res.Book.Id,
		&res.Book.Title,
		&isExpired,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("reservation with id %d not found", id)
		}
		return nil, err
	}

	if adminId != nil {
		res.AdminAccount.Id = *adminId
		res.AdminAccount.Name = *adminName
	} else {
		res.AdminAccount = nil
	}

	if isExpired {
		res.Status = model.ReservationExpired
	}

	return &res, nil
}

func (rr *reservationRepository) GetReservationsByBookId(id int, status string) (*[]model.Reservation, error) {
	query := `
	SELECT r.id            AS reservation_id  ,
	       r.reserved_at,
	       r.expires_at,
	       r.borrowed_days,
	       r.status        AS reservation_status,
	       r.fk_user_id    AS user_id,
	       usr.name        AS user_name,
	       r.fk_admin_id   AS admin_id,
	       adm.name        AS admin_name,
	       r.fk_book_id    AS book_id,
	       b.title         AS book_title,
		   (CURRENT_TIMESTAMP > r.expires_at) as is_expired
	FROM 
	       reservation r
	JOIN 
	       user_account usr ON r.fk_user_id = usr.id
	LEFT JOIN 
	       user_account adm ON r.fk_admin_id = adm.id
	JOIN 
	       book b ON r.fk_book_id = b.id
	WHERE 
		    r.fk_book_id = $1
   `

	var args []interface{}
	args = append(args, id)

	if status != "" {
		query += "AND status = $2"
		args = append(args, status)
	}

	rows, err := rr.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("error fetching reservations: %w", err)
	}
	defer rows.Close()

	reservations := make([]model.Reservation, 0)
	for rows.Next() {
		var res model.Reservation
		res.UserAccount = &user.Account{}
		res.AdminAccount = &user.Account{}
		var adminId *int
		var adminName *string
		var isExpired bool

		if err := rows.Scan(
			&res.Id,
			&res.ReservedAt,
			&res.ExpiresAt,
			&res.BorrowedDays,
			&res.Status,
			&res.UserAccount.Id,
			&res.UserAccount.Name,
			&adminId,
			&adminName,
			&res.Book.Id,
			&res.Book.Title,
			&isExpired,
		); err != nil {
			return nil, err
		}

		if adminId != nil {
			res.AdminAccount.Id = *adminId
			res.AdminAccount.Name = *adminName
		} else {
			res.AdminAccount = nil
		}

		if isExpired {
			res.Status = model.ReservationExpired
		}
		reservations = append(reservations, res)
	}
	return &reservations, nil
}

func (rr *reservationRepository) UpdateReservationStatus(reservationID int, status string, adminID int) error {
	query := `UPDATE reservation SET status = $1, fk_admin_id = $2 WHERE id = $3`
	_, err := rr.db.Exec(query, status, adminID, reservationID)
	if err != nil {
		return fmt.Errorf("failed to update reservation status: %w", err)
	}
	return nil
}

func (rr *reservationRepository) CancelReservation(id, adminId int) error {
	query := `UPDATE reservation SET status = 'cancelled', fk_admin_id = $1 WHERE id = $2`
	_, err := rr.db.Exec(query, id, adminId)
	if err != nil {
		return err
	}
	return nil
}
