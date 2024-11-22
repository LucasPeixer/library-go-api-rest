package usecase

import (
	"go-api/model"
	"go-api/repository"
)

type LoanUseCaseInterface interface {
	GetUserLoans(userID int) ([]model.Loan, error)
}

type LoanUseCase struct {
	LoanRepo repository.LoanRepositoryInterface
}

func NewLoanUseCase(loanRepo repository.LoanRepositoryInterface) *LoanUseCase {
	return &LoanUseCase{LoanRepo: loanRepo}
}

func (lu *LoanUseCase) GetUserLoans(userID int) ([]model.Loan, error) {
	return lu.LoanRepo.GetUserLoans(userID)
}
