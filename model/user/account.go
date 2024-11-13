package user

type Account struct {
	Id           int         `json:"id"`
	Name         string      `json:"name"`
	Phone        string      `json:"phone"`
	Email        string      `json:"email"`
	PasswordHash string      `json:"password_hash"`
	AccountRole  AccountRole `json:"account_role"`
	IsActive     bool        `json:"is_active"`
}
