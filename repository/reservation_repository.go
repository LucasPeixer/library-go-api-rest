package repository

import (
	"database/sql"
	"go-api/model"
	"strings"
	"fmt"
)

type ReservationRepositoryInterface interface {
	GetReservationsByFilters(userName, status, reservedAt string) ([]model.Reservation, error)
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
