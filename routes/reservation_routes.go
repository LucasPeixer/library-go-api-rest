package routes

import (
	"go-api/controller"
	"go-api/initializers"
	"go-api/middleware"
	"go-api/repository"
	"go-api/usecase"

	"github.com/gin-gonic/gin"
)

// ReservationRoutes registra todas as rotas de reserva.
func ReservationRoutes(rg *gin.RouterGroup) {
	// Criação das dependências do repositório, usecase e controller para reserva
	reservationRepository := repository.NewReservationRepository(initializers.DB)
	reservationUseCase := usecase.NewReservationUseCase(reservationRepository)
	reservationController := controller.NewReservationController(reservationUseCase)

	reservation := rg.Group("/reservations", middleware.JWTAuthMiddleware)
	{
		reservation.GET("/",middleware.RoleRequired("admin"), reservationController.GetReservationsByFilters)
	}
}
