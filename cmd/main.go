package main

import (
	"go-api/controller"
	"go-api/db"
	"go-api/repository"
	"go-api/usecase"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {

	dbConnection, err := db.ConnectDB()

	if(err != nil){
		panic(err)
	}
	//Camada Repository
	bookRepository := repository.NewBookRepository(dbConnection)
	//Camada UseCase
	bookUseCase := usecase.NewBookUseCase(bookRepository)
	//Camada de Controllers
	bookController := controller.NewBookController(bookUseCase)
	
	server := gin.Default()

	server.GET("/book", bookController.GetBooks)

	server.Run(":8000")
}