package repository

import (
	"database/sql"
	"go-api/model"
	"strings"
	"fmt"
)

type ReservationRepositoryInterface interface {
	GetReservationsByFilters(userName, status, reservedAt string) ([]model.Reservation, error)
	GetReservationByID(reservationID int) (*model.Reservation, error)
	CreateReservation(reservationRequest *model.ReservationRequest) (*model.Reservation, error)
	UpdateReservationStatus(reservationID int, status string, adminID int) error
}

type ReservationRepository struct {
	db *sql.DB
}

func NewReservationRepository(db *sql.DB) ReservationRepositoryInterface {
	return &ReservationRepository{db}
}

// Modificando o método para aceitar filtros opcionais.
func (rr *ReservationRepository) GetReservationsByFilters(userName, status, reservedAt string) ([]model.Reservation, error) {
	// Usando strings.Builder para construir a query
	var sb strings.Builder
	sb.WriteString("SELECT * FROM reservation WHERE 1=1")

	// Armazenará os parâmetros de consulta
	var args []interface{}
	paramIndex := 1

	// Adicionando condições dinâmicas para os filtros passados
	if userName != "" {
		sb.WriteString(" AND user_name ILIKE $" + fmt.Sprint(paramIndex))
		args = append(args, "%"+userName+"%")
		paramIndex++
	}

	if status != "" {
		sb.WriteString(" AND status = $" + fmt.Sprint(paramIndex))
		args = append(args, status)
		paramIndex++
	}

	if reservedAt != "" {
		sb.WriteString(" AND reserved_at::date = $" + fmt.Sprint(paramIndex))
		args = append(args, reservedAt)
		paramIndex++
	}

	// Executando a consulta com os filtros aplicados
	query := sb.String()
	rows, err := rr.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("error executing query: %w", err)
	}
	defer rows.Close()

	// Processando os resultados
	var reservations []model.Reservation
	for rows.Next() {
		var res model.Reservation
		if err := rows.Scan(&res.ID, &res.ReservedAt, &res.ExpiresAt, &res.BorrowedDays, &res.Status, &res.UserID, &res.AdminID, &res.BookID); err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
		reservations = append(reservations, res)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return reservations, nil
}

func (rr *ReservationRepository) GetReservationByID(reservationID int) (*model.Reservation, error) {

	query := `
		SELECT id, reserved_at, expires_at, borrowed_days, status, fk_user_id, fk_admin_id, fk_book_id
		FROM reservation
		WHERE id = $1
	`

	reservation := &model.Reservation{}

	row := rr.db.QueryRow(query, reservationID)
	err := row.Scan(
		&reservation.ID,
		&reservation.ReservedAt,
		&reservation.ExpiresAt,
		&reservation.BorrowedDays,
		&reservation.Status,
		&reservation.UserID,
		&reservation.AdminID,
		&reservation.BookID,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("reservation with ID %d not found", reservationID)
		}
		return nil, fmt.Errorf("error fetching reservation: %w", err)
	}

	return reservation, nil
}

func (rr *ReservationRepository) CreateReservation(reservationRequest *model.ReservationRequest) (*model.Reservation, error) {
	
	query := `
		INSERT INTO reservation (borrowed_days, fk_user_id, fk_book_id)
		VALUES ($1, $2, $3)
		RETURNING id, reserved_at, expires_at, borrowed_days, status, fk_user_id, fk_book_id`
	
	var reservation model.Reservation
	err := rr.db.QueryRow(query, reservationRequest.BorrowedDays, reservationRequest.UserID, reservationRequest.BookID).
		Scan(&reservation.ID, &reservation.ReservedAt, &reservation.ExpiresAt, &reservation.BorrowedDays, &reservation.Status, &reservation.UserID, &reservation.BookID)
	if err != nil {
		return nil, fmt.Errorf("error inserting reservation: %w", err)
	}

	return &reservation, nil
}

func (rr *ReservationRepository) UpdateReservationStatus(reservationID int, status string, adminID int) error {
	query := `UPDATE reservation SET status = $1, fk_admin_id = $2 WHERE id = $3`
	_, err := rr.db.Exec(query, status, adminID, reservationID)
	if err != nil {
		return fmt.Errorf("failed to update reservation status: %w", err)
	}
	return nil
}