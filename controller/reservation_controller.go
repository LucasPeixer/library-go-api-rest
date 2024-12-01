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
	userName := c.DefaultQuery("user_name", "")
	status := c.DefaultQuery("status", "")
	reservedAt := c.DefaultQuery("reserved_at", "")

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

	userID, err := strconv.Atoi(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Pegando os dados da reserva do corpo da requisição
	var request model.ReservationRequest
	request.UserID = userID

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid data"})
		return
	}

	// Criando a reserva
	reservation, err := rc.useCase.CreateReservation(&request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("error when creating reservation: %v", err)})
		return
	}

	// Retornando o sucesso com a reserva criada
	c.JSON(http.StatusCreated, gin.H{
		"id":            reservation.ID,
		"reserved_at":   reservation.ReservedAt,
		"expires_at":    reservation.ExpiresAt,
		"status":        reservation.Status,
		"borrowed_days": reservation.BorrowedDays,
		"fk_user_id":    reservation.UserID,
		"fk_book_id":    reservation.BookID,
	})
}
