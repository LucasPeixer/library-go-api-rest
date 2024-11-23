package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"go-api/model"
	"go-api/model/user"
	"strconv"
)

type UserRepository interface {
	CreateUser(name, phone, email, passwordHash string, fkAccountRole int) error
	GetUserByEmail(email string) (*user.Account, error)
	GetUsersByFilters(name, email string) (*[]user.Account, error)
	GetUserById(id int) (*user.Account, error)
	GetUserLoans(userID int) ([]model.Loan, error)
	GetUserReservation(userID int) ([]*model.Reservation, error)
	ActivateUser(id int) error
	DeactivateUser(id int) error
	DeleteUser(id int) error
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db}
}

func (ur *userRepository) CreateUser(name, phone, email, passwordHash string, fkAccountRole int) error {
	query := `
        INSERT INTO user_account (name, phone, email, password_hash, fk_account_role)
        VALUES ($1, $2, $3, $4, $5)
    `

	_, err := ur.db.Exec(query, name, phone, email, passwordHash, fkAccountRole)
	if err != nil {
		return fmt.Errorf("error creating user: %v", err)
	}
	return nil
}

func (ur *userRepository) GetUserByEmail(email string) (*user.Account, error) {
	query := `
	SELECT 
	       ua.id             AS user_id,
	       ua.name           AS user_name,
	       ua.phone          AS user_phone,
	       ua.email          AS user_email,
	       ua.password_hash  AS user_password_hash,
	       ua.is_active      AS user_is_active,
	       ar.id             AS account_role_id,
	       ar.name           AS account_role_name
	FROM 
	       user_account ua
	JOIN
	       account_role ar ON ua.fk_account_role = ar.id
	WHERE email = $1`

	var userAccount user.Account

	err := ur.db.QueryRow(query, email).Scan(
		&userAccount.Id,
		&userAccount.Name,
		&userAccount.Phone,
		&userAccount.Email,
		&userAccount.PasswordHash,
		&userAccount.IsActive,
		&userAccount.AccountRole.Id,
		&userAccount.AccountRole.Name,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user with email %s not found", email)
		}
		return nil, err // Return other errors if they occur
	}

	return &userAccount, nil
}

func (ur *userRepository) GetUsersByFilters(name, email string) (*[]user.Account, error) {
	query := `
	SELECT 
	       ua.id             AS user_id,
	       ua.name           AS user_name,
	       ua.phone          AS user_phone,
	       ua.email          AS user_email,
	       ua.is_active      AS user_is_active,
	       ar.id             AS account_role_id,
	       ar.name           AS account_role_name
	FROM 
	       user_account ua
	JOIN
	       account_role ar ON ua.fk_account_role = ar.id
	WHERE 1=1`

	var args []interface{}

	if name != "" {
		query += " AND ua.name ILIKE $" + strconv.Itoa(len(args)+1)
		args = append(args, "%"+name+"%")
	}

	if email != "" {
		query += " AND ua.email ILIKE $" + strconv.Itoa(len(args)+1)
		args = append(args, "%"+email+"%")
	}

	rows, err := ur.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	userAccountList := make([]user.Account, 0)
	for rows.Next() {
		var userAccount user.Account
		err := rows.Scan(
			&userAccount.Id,
			&userAccount.Name,
			&userAccount.Phone,
			&userAccount.Email,
			&userAccount.IsActive,
			&userAccount.AccountRole.Id,
			&userAccount.AccountRole.Name,
		)
		if err != nil {
			return nil, err
		}
		userAccountList = append(userAccountList, userAccount)
	}
	return &userAccountList, nil
}

func (ur *userRepository) GetUserById(id int) (*user.Account, error) {
	query := `
	SELECT 
	       ua.id             AS user_id,
	       ua.name           AS user_name,
	       ua.phone          AS user_phone,
	       ua.email          AS user_email,
	       ua.is_active      AS user_is_active,
	       ar.id             AS account_role_id,
	       ar.name           AS account_role_name
	FROM 
	       user_account ua
	JOIN
	       account_role ar ON ua.fk_account_role = ar.id
	WHERE ua.id = $1`

	var userAccount user.Account

	err := ur.db.QueryRow(query, id).Scan(
		&userAccount.Id,
		&userAccount.Name,
		&userAccount.Phone,
		&userAccount.Email,
		&userAccount.IsActive,
		&userAccount.AccountRole.Id,
		&userAccount.AccountRole.Name,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user with id %d not found", id)
		}
		return nil, err // Return other errors if they occur
	}

	return &userAccount, nil
}

func (ur *userRepository) GetUserLoans(userID int) ([]model.Loan, error){
	query := `
			SELECT l.id, l.loaned_at, l.return_by, l.returned_at, l.status, 
       l.fk_admin_id, l.fk_book_stock_id, l.fk_reservation_id
			FROM loan l
			JOIN reservation r ON l.fk_reservation_id = r.id 
			WHERE r.fk_user_id = $1;`

	rows, err := ur.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var loans []model.Loan
	for rows.Next() {
		var loan model.Loan
		err := rows.Scan(
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
		loans = append(loans, loan)
	}

	return loans, nil
}

func (ur *userRepository) GetUserReservation(userID int) ([]*model.Reservation, error) {
	query := `
		SELECT 
			r.id, 
			r.fk_user_id, 
			r.fk_book_id, 
			r.borrowed_days, 
			r.status, 
			r.reserved_at, 
			r.expires_at
		FROM reservation r
		WHERE r.fk_user_id = $1
	`
	rows, err := ur.DB.Query(query, userID)
	if err != nil {
		log.Printf("Error executing query: %v", err)
		return nil, fmt.Errorf("error fetching reservations: %w", err)
	}
	defer rows.Close()

	var reservations []*model.Reservation

	for rows.Next() {
		var reservation model.Reservation
		if err := rows.Scan(&reservation.ID, &reservation.FkUserId, &reservation.FkBookId, &reservation.BorrowedDays, &reservation.Status, &reservation.ReservedAt, &reservation.ExpiresAt); err != nil {
			log.Printf("Error reading reservation data: %v", err)
			return nil, fmt.Errorf("error reading reservation data: %w", err)
		}

		reservations = append(reservations, &reservation)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error iterating over reservations: %v", err)
		return nil, fmt.Errorf("error iterating over reservations: %w", err)
	}

	return reservations, nil
}

func (ur *userRepository) ActivateUser(id int) error {
	err := ur.toggleUser(id, true)
	if err != nil {
		return fmt.Errorf("error activating user: %v", err)
	}
	return nil
}

func (ur *userRepository) DeactivateUser(id int) error {
	err := ur.toggleUser(id, false)
	if err != nil {
		return fmt.Errorf("error deactivating user: %v", err)
	}
	return nil
}

func (ur *userRepository) toggleUser(id int, status bool) error {
	query := `
	UPDATE user_account 
	SET is_active = $1 
	WHERE id = $2
	RETURNING id
	`
	var userId int
	err := ur.db.QueryRow(query, status, id).Scan(&userId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("user with id %d not found", id)
		}
		return err
	}
	return nil
}

func (ur *userRepository) DeleteUser(id int) error {
	query := `
		DELETE FROM user_account
		WHERE id = $1
		RETURNING id
	`
	var userId int
	err := ur.db.QueryRow(query, id).Scan(&userId)
	if err != nil {
		return fmt.Errorf("error deleting user: %v", err)
	}
	return nil
}
