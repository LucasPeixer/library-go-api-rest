package controller

import (
	"github.com/gin-gonic/gin"
	"go-api/model"
	"go-api/usecase"
	"net/http"
	"strconv"
)

type ReservationController interface {
	GetReservationsByFilters(c *gin.Context)
	CreateReservation(c *gin.Context)
}

type reservationController struct {
	useCase usecase.ReservationUseCase
}

func NewReservationController(useCase usecase.ReservationUseCase) ReservationController {
	return &reservationController{useCase: useCase}
}

func (rc *reservationController) GetReservationsByFilters(c *gin.Context) {
	// Pegando os parâmetros da query string
	userName := c.Query("user_name")
	status := c.Query("status")
	reservedAt := c.Query("reserved_at")

	reservations, err := rc.useCase.GetReservationsByFilters(userName, model.ReservationStatus(status), reservedAt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, reservations)
}

func (rc *reservationController) CreateReservation(c *gin.Context) {
	userIDStr, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized user"})
		return
	}

	userId, err := strconv.Atoi(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var i struct {
		BorrowedDays int `json:"borrowed_days" binding:"required"`
		BookId       int `json:"book_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&i); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid reservation input"})
		return
	}

	reservation, err := rc.useCase.CreateReservation(i.BorrowedDays, userId, i.BookId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, reservation)

}
