package routes

import (
	"go-api/controller"
	"go-api/initializers"
	"go-api/middleware"
	"go-api/repository"
	"go-api/usecase"

	"github.com/gin-gonic/gin"
)

// UserRoutes registra todas as rotas de usuário.
func UserRoutes(rg *gin.RouterGroup) {
	userRepository := repository.NewUserRepository(initializers.DB)
	userUseCase := usecase.NewUserUseCase(userRepository)
	userController := controller.NewUserController(userUseCase)

	rg.POST("/login", userController.Login)
	rg.POST("/register",
		middleware.JWTAuthMiddleware,
		middleware.RoleRequired("admin"),
		userController.Register,
	)

	users := rg.Group("/users", middleware.JWTAuthMiddleware, middleware.RoleRequired("admin"))
	{
		users.POST("/register", userController.Register)
		users.GET("/", userController.GetUsersByFilters)
		users.GET("/:id", userController.GetUserById)
		users.PUT("/activate/:id", userController.ToggleUser("activate"))
		users.PUT("/deactivate/:id", userController.ToggleUser("deactivate"))
		users.DELETE("/delete/:id", userController.DeleteUser)

		reservations := users.Group("/:id/reservations")
		{
			reservations.GET("/", userController.GetUserReservations)
			reservations.PUT("/cancel/:reservation-id", userController.CancelUserReservation)
		}
	}

	user := rg.Group("/user")
	{
		user.POST("/login", userController.Login)

		reservations := user.Group("/reservations", middleware.JWTAuthMiddleware)
		{
			reservations.GET("/", userController.GetLoggedUserReservations)
			reservations.PUT("/:id/cancel", userController.CancelLoggedUserReservation)
		}

		user.GET("/loans", middleware.JWTAuthMiddleware, userController.GetUserLoans)
	}
}
