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
	bookStockRepository := repository.NewBookRepository(initializers.DB)
	userRepository := repository.NewUserRepository(initializers.DB)
	bookRepository := repository.NewBookRepository(initializers.DB)
	reservationUseCase := usecase.NewReservationUseCase(reservationRepository,userRepository,bookRepository)

	loanUseCase := usecase.NewLoanUseCase(loanRepository, reservationRepository, bookStockRepository)
	loanController := controller.NewLoanController(loanUseCase, reservationUseCase)

	loan := rg.Group("/loans", middleware.JWTAuthMiddleware)
	{
		loan.POST("/create",middleware.RoleRequired("admin"),loanController.CreateLoan)
		loan.PUT("/finish-loan/:id",middleware.RoleRequired("admin"),loanController.UpdateLoan)
	}
}
