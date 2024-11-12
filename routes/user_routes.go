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
}
