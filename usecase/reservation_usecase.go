package usecase

import (
	"go-api/repository"
	"go-api/model"
)

type ReservationUseCaseInterface interface {
	GetReservationsByFilters(userName, status, reservedAt string) ([]model.Reservation, error)
	CreateReservation(reservation *model.Reservation) (*model.Reservation, error)
}

type ReservationUseCase struct {
	ReservationRepo repository.ReservationRepositoryInterface
	userRepo        repository.UserRepositoryInterface
	bookRepo        repository.BookRepositoryInterface
}

// NewReservationUseCase cria e retorna uma nova instância de ReservationUseCase
func NewReservationUseCase(repo repository.ReservationRepositoryInterface) ReservationUseCaseInterface {
	return &ReservationUseCase{
		ReservationRepo: repo,
		userRepo:        userRepo,
		bookRepo:        bookRepo,
	}
}

func (ru *ReservationUseCase) GetReservationsByFilters(userName, status, reservedAt string) ([]model.Reservation, error) {
	return ru.ReservationRepo.GetReservationsByFilters(userName, status, reservedAt)
}

func (ru *ReservationUseCase) CreateReservation(reservation *model.Reservation) (*model.Reservation, error) {

	//Verificar se a conta do usuário está ativa 
	user, err := ru.userRepo.GetUserByID(reservation.UserID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar usuário: %w", err)
	}
	if user.Status != "active" {
		return nil, fmt.Errorf("usuário não está ativo")
	}

	//Verificar se o usuário não tem empréstimos em atraso
	/*activeLoans, err := ru.userRepo.GetUserLoans(reservation.UserID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar empréstimos do usuário: %w", err)
	}
	for _, loan := range activeLoans {
	// Verificar se o empréstimo está em atraso
	if loan.ReturnBy.Before(time.Now()) {
			return nil, fmt.Errorf("usuário tem empréstimos em atraso")
		}
	}*/


	//Verificar se o usuário não excedeu o limite de 5 reservas/empréstimos ativos 
	/*activeReservations, err := ru.reservationRepo.GetReservationsByFilters(fmt.Sprintf("%d", reservation.UserID), "pending", "")
	if err != nil {
		return nil, fmt.Errorf("error when searching for active user reservations: %w", err)
	}

	activeLoans, err := ru.userRepo.GetUserLoans(reservation.UserID)
	if err != nil {
		return nil, fmt.Errorf("error when searching for active user loans: %w", err)
	}

	totalActive := len(activeReservations) + len(activeLoans)

	if totalActive >= 5 {
		return nil, fmt.Errorf("user already has 5 or more active reservations/loans")
	}*/

	book, err := ru.bookRepo.GetBookByID(reservation.BookID)
	if err != nil {
		return nil, fmt.Errorf("error when searching for book: %w", err)
	}

	if book.Amount <= 0 {
		return nil, fmt.Errorf("book out of stock")
	}

	newReservation, err := ru.reservationRepo.CreateReservation(reservation)
	if err != nil {
		return nil, fmt.Errorf("error when creating reservation: %w", err)
	}

	return newReservation, nil
}
