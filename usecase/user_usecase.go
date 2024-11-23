package usecase

import (
	"errors"
	"go-api/model"
	"go-api/model/user"
	"go-api/repository"
	"go-api/utils"
	"fmt"
)

type UserUseCase interface {
	Login(email, password string) (string, error)
	Register(name, phone, email, passwordHash string, fkAccountRole int) error
	GetUsersByFilters(name, email string) (*[]user.Account, error)
	GetUserById(id int) (*user.Account, error)
	GetUserLoans(userID int) ([]model.Loan, error)
	ActivateUser(id int) error
	DeactivateUser(id int) error
	DeleteUser(id int) error
}

type userUseCase struct {
	userRepo repository.UserRepository
}

func NewUserUseCase(repo repository.UserRepository) UserUseCase {
	return &userUseCase{userRepo: repo}
}

// Login busca o usuário pelo email, verifica a senha hashed e retorna um token JWT.
func (uu *userUseCase) Login(email, password string) (string, error) {
	// Busca o usuário pelo seu email
	userAccount, err := uu.userRepo.GetUserByEmail(email)
	if err != nil {
		return "", err
	}
	if !utils.CheckPasswordHash(password, userAccount.PasswordHash) {
		return "", errors.New("invalid credentials")
	}
	// Gera um token JWT com ID e Role se atendido todas as condições
	return utils.GenerateJWT(userAccount.Id, userAccount.AccountRole.Name)
}

// Register cria um novo usuário no banco de dados.
func (uu *userUseCase) Register(name, phone, email, passwordHash string, fkAccountRole int) error {
	return uu.userRepo.CreateUser(name, phone, email, passwordHash, fkAccountRole)
}

func (uu *userUseCase) GetUsersByFilters(name, email string) (*[]user.Account, error) {
	return uu.userRepo.GetUsersByFilters(name, email)
}

func (uu *userUseCase) GetUserById(id int) (*user.Account, error) {
	return uu.userRepo.GetUserById(id)
}

func (uu *userUseCase) GetUserLoans(userID int) ([]model.Loan, error) {
	return uu.userRepo.GetUserLoans(userID)
}

func (uu *userUseCase) GetUserReservations(userID int) ([]*model.Reservation, error) {
	// Chama o repositório para pegar as reservas do usuário
	reservations, err := uu.userRepo.GetUserReservation(userID)
	if err != nil {
		return nil, fmt.Errorf("error fetching user reservations: %w", err)
	}
	return reservations, nil
}

func (uu *userUseCase) ActivateUser(id int) error {
	return uu.userRepo.ActivateUser(id)
}

func (uu *userUseCase) DeactivateUser(id int) error {
	return uu.userRepo.DeactivateUser(id)
}

func (uu *userUseCase) DeleteUser(id int) error {
	return uu.userRepo.DeleteUser(id)
}
