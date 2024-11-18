package usecase

import (
	"go-api/repository"
	"go-api/model"
)

type ReservationUseCaseInterface interface {
	GetReservationsByFilters(userName, status, reservedAt string) ([]model.Reservation, error)
}

type ReservationUseCase struct {
	Repo repository.ReservationRepositoryInterface
}

// NewReservationUseCase cria e retorna uma nova instância de ReservationUseCase
func NewReservationUseCase(repo repository.ReservationRepositoryInterface) ReservationUseCaseInterface {
	return &ReservationUseCase{
		Repo: repo,
	}
}

func (u *ReservationUseCase) GetReservationsByFilters(userName, status, reservedAt string) ([]model.Reservation, error) {
	return u.Repo.GetReservationsByFilters(userName, status, reservedAt)
}
