package controller

import (
	"net/http"
	"go-api/usecase"
	"github.com/gin-gonic/gin"
)

type ReservationController struct {
	UseCase usecase.ReservationUseCaseInterface
}

// NewReservationController cria e retorna um novo controlador de reserva
func NewReservationController(useCase usecase.ReservationUseCaseInterface) *ReservationController {
	return &ReservationController{UseCase: useCase}
}

func (c *ReservationController) GetReservationsByFilters(ctx *gin.Context) {
	// Pegando os parâmetros da query string
	userName := ctx.DefaultQuery("user_name", "")
	status := ctx.DefaultQuery("status", "")
	reservedAt := ctx.DefaultQuery("reserved_at", "")

	// Chamando o usecase com os filtros passados
	reservations, err := c.UseCase.GetReservationsByFilters(userName, status, reservedAt) // Corrigido para UseCase
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, reservations)
}
