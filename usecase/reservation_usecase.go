package usecase

import (
	"go-api/repository"
	"go-api/model"
)

type ReservationUseCaseInterface interface {
	GetReservationsByFilters(userName, status, reservedAt string) ([]model.Reservation, error)
}

type ReservationUseCase struct {
	ReservationRepo repository.ReservationRepositoryInterface
}

// NewReservationUseCase cria e retorna uma nova instância de ReservationUseCase
func NewReservationUseCase(repo repository.ReservationRepositoryInterface) ReservationUseCaseInterface {
	return &ReservationUseCase{
		ReservationRepo: repo,
	}
}

func (ru *ReservationUseCase) GetReservationsByFilters(userName, status, reservedAt string) ([]model.Reservation, error) {
	return ru.ReservationRepo.GetReservationsByFilters(userName, status, reservedAt)
}
