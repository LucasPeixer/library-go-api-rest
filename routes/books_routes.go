package routes

import (
	"go-api/controller"
	"go-api/db"
	"go-api/repository"
	"go-api/usecase"

	"github.com/gin-gonic/gin"
)


func BooksRouter(rg *gin.RouterGroup) {
	// Inicializando toda a estrutura de serviços dos Livros
	dbConnection, err := db.ConnectDB()
	if(err != nil){
		panic(err)
	}

	bookRepository := repository.NewBookRepository(dbConnection)
	bookUseCase := usecase.NewBookUseCase(bookRepository)
	bookController := controller.NewBookController(bookUseCase)

		//Separando as rotas dos livros
		books := rg.Group("/livro")
		{
			books.GET("/", bookController.GetAllBooks)
			books.GET("/:{id}", bookController.GetBooks)
			books.POST("/registro", bookController.CreateBook)
		}

}