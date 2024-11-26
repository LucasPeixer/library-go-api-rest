package controller

import (
	"github.com/gin-gonic/gin"
	"go-api/model"
	"go-api/usecase"
	"net/http"
	"strconv"
	"time"
	"fmt"
)

type LoanControllerInterface interface {
	CreateLoan(c *gin.Context)
	UpdateLoan(c *gin.Context)
	GetLoansByFilters(c *gin.Context)
}

type LoanController struct {
	loanUsecase        usecase.LoanUseCaseInterface
	reservationUsecase usecase.ReservationUseCaseInterface
}

// Corrigindo o nome do campo para 'loanUsecase'
func NewLoanController(loanUsecase usecase.LoanUseCaseInterface, reservationUsecase usecase.ReservationUseCaseInterface) *LoanController {
	return &LoanController{
		loanUsecase:        loanUsecase,
		reservationUsecase: reservationUsecase,
	}
}

func (lc *LoanController) GetLoansByFilters(c *gin.Context) {
	loanedAtStr := c.DefaultQuery("loaned_at", "")
	returnByStr := c.DefaultQuery("return_by", "")
	returnedAtStr := c.DefaultQuery("returned_at", "")
	status := c.DefaultQuery("status", "")
	bookStockIDStr := c.DefaultQuery("book_stock_id", "")
	reservationIDStr := c.DefaultQuery("reservation_id", "")

	filters := map[string]interface{}{
			"loaned_at":   loanedAtStr,
			"return_by":   returnByStr,
			"returned_at": returnedAtStr,
			"status":      status,
	}

	if bookStockIDStr != "" {
			bookStockID, err := strconv.Atoi(bookStockIDStr)
			if err == nil {
					filters["book_stock_id"] = bookStockID
			}
	}

	if reservationIDStr != "" {
			reservationID, err := strconv.Atoi(reservationIDStr)
			if err == nil {
					filters["reservation_id"] = reservationID
			}
	}

	loans, err := lc.loanUsecase.GetLoansByFilter(filters)
	if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
	}

	c.JSON(http.StatusOK, loans)
}

func (lc *LoanController) CreateLoan(c *gin.Context) {
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
		AdminID:       &adminId,
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

func (lc *LoanController) UpdateLoan(c *gin.Context) {
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

	err = lc.loanUsecase.UpdateLoan(request, adminId, loanId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "loan updated successfully"})
}