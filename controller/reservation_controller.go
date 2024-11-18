package controller

import (
	"net/http"
	"go-api/usecase"
	"github.com/gin-gonic/gin"
)

type ReservationController struct {
	UseCase usecase.ReservationUseCaseInterface
}

func NewReservationController(useCase usecase.ReservationUseCaseInterface) *ReservationController {
	return &ReservationController{UseCase: useCase}
}

func (rc *ReservationController) GetReservationsByFilters(c *gin.Context) {
	// Pegando os parâmetros da query string
	userName := c.DefaultQuery("user_name", "")
	status := c.DefaultQuery("status", "")
	reservedAt := c.DefaultQuery("reserved_at", "")

	reservations, err := rc.UseCase.GetReservationsByFilters(userName, status, reservedAt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, reservations)
}
