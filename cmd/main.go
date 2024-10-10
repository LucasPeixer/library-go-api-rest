package main

import (
	"go-api/routes"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {	
//Camada Repository
//genreRepository := repository.NewGenreRepository(dbConnection)
//Camada UseCase
//genreUseCase := usecase.NewGenreUseCase(genreRepository)
//Camada de Controllers
//genreController := controller.NewGenreController(genreUseCase)
	
	
		
	router := gin.Default()
	routes.RegisterRoutes(router)
	
	//server.GET("/generos", genreController.GetGenres)
	
	router.Run(":8000")
}