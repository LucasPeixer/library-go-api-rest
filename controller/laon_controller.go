package controller

import (
	"net/http"
	"go-api/usecase"
	"go-api/model"
	"github.com/gin-gonic/gin"
	"time"
)

type LoanControllerInterface interface {
	CreateLoan(c *gin.Context)
}

type LoanController struct {
	loanUsecase       usecase.LoanUseCaseInterface
	reservationUsecase usecase.ReservationUseCaseInterface
}

// Corrigindo o nome do campo para 'loanUsecase'
func NewLoanController(loanUsecase usecase.LoanUseCaseInterface, reservationUsecase usecase.ReservationUseCaseInterface) *LoanController {
	return &LoanController{
		loanUsecase:       loanUsecase,
		reservationUsecase: reservationUsecase,
	}
}

func (lc *LoanController) CreateLoan(c *gin.Context) {
	var loanRequest model.LoanRequest

	// Bind JSON e validação do corpo da requisição
	if err := c.ShouldBindJSON(&loanRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	// Obtenção da reserva usando o 'reservationUsecase'
	reservation, err := lc.reservationUsecase.GetReservationByID(loanRequest.ReservationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Cálculo da data de devolução
	returnBy := time.Now().AddDate(0, 0, reservation.BorrowedDays)

	// Criação do empréstimo
	loan := &model.LoanRequest{
		ReturnBy:      returnBy,
		BookStockID:   loanRequest.BookStockID,
		ReservationID: loanRequest.ReservationID,
	}

	// Criação do empréstimo e atualização da reserva
	createdLoan, err := lc.loanUsecase.CreateLoanAndUpdateReservation(loan)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Retorno do empréstimo criado
	c.JSON(http.StatusCreated, createdLoan)
}
