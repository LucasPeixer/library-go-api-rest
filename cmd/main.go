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
	bookRepository := repository.NewBookRepository(dbConnection)
	//Camada UseCase
	bookUseCase := usecase.NewBookUseCase(bookRepository)
	//Camada de Controllers
	bookController := controller.NewBookController(bookUseCase)
	
	server := gin.Default()

	server.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
				"message": "pong",
		})
	})

	server.GET("/book", bookController.GetBooks)

	server.Run(":8000")
}