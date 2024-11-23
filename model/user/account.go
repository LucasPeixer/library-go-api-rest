package user

type Account struct {
	Id           int         `json:"id"`
	Name         string      `json:"name"`
	Cpf          string      `json:"cpf"`
	Phone        string      `json:"phone"`
	Email        string      `json:"email"`
	PasswordHash string      `json:"password_hash"`
	AccountRole  AccountRole `json:"account_role"`
	IsActive     bool        `json:"is_active"`
}
