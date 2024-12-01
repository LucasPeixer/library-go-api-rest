package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-api/model"
	"go-api/usecase"
	"net/http"
	"strconv"
)

type LoanController interface {
	CreateLoan(c *gin.Context)
	UpdateLoan(c *gin.Context)
}

type loanController struct {
	loanUseCase        usecase.LoanUseCase
	reservationUseCase usecase.ReservationUseCase
}

func NewLoanController(loanUseCase usecase.LoanUseCase, reservationUseCase usecase.ReservationUseCase) LoanController {
	return &loanController{
		loanUseCase:        loanUseCase,
		reservationUseCase: reservationUseCase,
	}
}

func (lc *loanController) CreateLoan(c *gin.Context) {
	adminIdStr, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized user"})
		return
	}

	adminId, err := strconv.Atoi(adminIdStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid admin ID"})
		return
	}

	var i struct {
		ReservationId int `json:"reservation_id" binding:"required"`
		BookStockId   int `json:"book_stock_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&i); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	loan, err := lc.loanUseCase.CreateLoanAndUpdateReservation(i.ReservationId, i.BookStockId, adminId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, loan)
}

func (lc *loanController) UpdateLoan(c *gin.Context) {
	loanIDStr := c.Param("id")
	loanId, err := strconv.Atoi(loanIDStr)
	fmt.Printf("Loan ID: %d\n", loanId)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid loan ID"})
		return
	}

	var request model.LoanUpdateRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	// Obtém o ID do administrador do JWT
	adminIdStr, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	adminId, err := strconv.Atoi(adminIdStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid admin ID"})
		return
	}

	err = lc.loanUseCase.UpdateLoan(request, adminId, loanId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "loan updated successfully"})
}
