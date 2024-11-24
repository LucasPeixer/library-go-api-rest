package controller

import (
	"net/http"
	"go-api/usecase"
	"go-api/model"
	"github.com/gin-gonic/gin"
	"database/sql" 
)

type LoanController struct {
	UseCase usecase.LoanUseCaseInterface
}


func NewLoanController(useCase usecase.LoanUseCaseInterface) *LoanController {
	return &LoanController{UseCase: useCase}
}


func (lc *LoanController) CreateLoan(c *gin.Context) {

	var loanRequest model.LoanRequest
	if err := c.ShouldBindJSON(&loanRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	loan, err := lc.UseCase.CreateLoanAndUpdateReservation(&loanRequest) // Correção: Usando 'UseCase' no lugar de 'LoanUseCase'
	if err != nil {
		if err == sql.ErrNoRows { 
			c.JSON(http.StatusNotFound, gin.H{"error": "reservation not found or invalid"})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusCreated, loan)
}
