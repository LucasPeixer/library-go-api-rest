package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"go-api/model/user"
)

type UserRepository interface {
	CreateUser(name, phone, email, passwordHash string, fkAccountRole int) error
	GetUserByEmail(email string) (*user.Account, error)
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
