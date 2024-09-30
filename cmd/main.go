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
	genreRepository := repository.NewGenreRepository(dbConnection)
	//Camada UseCase
	bookUseCase := usecase.NewBookUseCase(bookRepository)
	genreUseCase := usecase.NewGenreUseCase(genreRepository)
	//Camada de Controllers
	bookController := controller.NewBookController(bookUseCase)
	genreController := controller.NewGenreController(genreUseCase)
	
	server := gin.Default()
	
	server.GET("/livros", bookController.GetBooks)
	server.GET("/generos", genreController.GetGenres)
	server.POST("/livros/registro", bookController.CreateBook)

	server.Run(":8000")
}