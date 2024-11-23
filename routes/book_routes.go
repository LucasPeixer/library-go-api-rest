package routes

import (
	"go-api/controller"
	"go-api/initializers"
	"go-api/middleware"
	"go-api/repository"
	"go-api/usecase"

	"github.com/gin-gonic/gin"
)

// BookRoutes registra todas as rotas de livro.
func BookRoutes(rg *gin.RouterGroup) {
	bookRepository := repository.NewBookRepository(initializers.DB)
	bookUseCase := usecase.NewBookUseCase(bookRepository)
	bookController := controller.NewBookController(bookUseCase)

	// Cria um grupo de rotas para '/books' que requerem autorização JWT, algumas com autorização 'admin'
	books := rg.Group("/books", middleware.JWTAuthMiddleware)
	{
		books.POST("/create", middleware.RoleRequired("admin"), bookController.CreateBook)
		books.GET("/", bookController.GetBooks)
		books.GET("/:id", bookController.GetBookById)
		books.PUT("/update/:id", middleware.RoleRequired("admin"), bookController.UpdateBook)
		books.DELETE("/delete/:id", middleware.RoleRequired("admin"), bookController.DeleteBook)

		stock := books.Group("/:id/stock", middleware.RoleRequired("admin"))
		{
			stock.POST("/add", bookController.AddStock)
			stock.GET("/", bookController.GetStock)
			stock.PUT("/update-status/:stock-id", bookController.UpdateStockStatus)
			stock.DELETE("/remove/:stock-id", bookController.RemoveStock)
		}
	}
}
