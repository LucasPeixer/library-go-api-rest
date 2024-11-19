package usecase

import (
	"errors"
	"go-api/model/user"
	"go-api/repository"
	"go-api/utils"
)

type UserUseCase interface {
	Login(email, password string) (string, error)
	Register(name, phone, email, passwordHash string, fkAccountRole int) error
	GetUsersByFilters(name, email string) (*[]user.Account, error)
	GetUserById(id int) (*user.Account, error)
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
	user, err := uu.userRepo.GetUserByEmail(email)
	if err != nil {
		return "", err
	}
	if !utils.CheckPasswordHash(password, user.PasswordHash) {
		return "", errors.New("invalid credentials")
	}
	// Gera um token JWT com ID e Role se atendido todas as condições
	return utils.GenerateJWT(user.Id, user.AccountRole.Name)
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
