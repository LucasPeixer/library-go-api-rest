package routes

import (
	"go-api/controller"
	"go-api/initializers"
	"go-api/middleware"
	"go-api/repository"
	"go-api/usecase"

	"github.com/gin-gonic/gin"
)

func LoanRoutes(rg *gin.RouterGroup) {
	loanRepository := repository.NewLoanRepository(initializers.DB)
	reservationRepository := repository.NewReservationRepository(initializers.DB)
	bookStockRepository := repository.NewBookStockRepository(initializers.DB)

	loanUseCase := usecase.NewLoanUseCase(loanRepository, reservationRepository, bookStockRepository)
	loanController := controller.NewLoanController(loanUseCase)

	loan := rg.Group("/loans", middleware.JWTAuthMiddleware)
	{
		loan.POST("/create", loanController.CreateLoanAndUpdateReservation)
	}
}
