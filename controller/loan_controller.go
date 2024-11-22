package controller

import (
	"go-api/usecase"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type LoanController interface {
	GetUserLoans(c *gin.Context)
}

type loanController struct {
	useCase usecase.LoanUseCaseInterface
}

func NewLoanController(useCase usecase.LoanUseCaseInterface) LoanController {
	return &loanController{useCase: useCase}
}

func (lc *loanController) GetUserLoans(c *gin.Context) {
	userIDParam := c.Param("id")
	userID, err := strconv.Atoi(userIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	loans, err := lc.useCase.GetUserLoans(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, loans)
}
