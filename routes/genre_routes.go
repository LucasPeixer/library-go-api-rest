package routes

import (
	"go-api/controller"
	"go-api/db"
	"go-api/repository"
	"go-api/usecase"

	"github.com/gin-gonic/gin"
)

func GenreRouter(gr *gin.RouterGroup){
	dbConnection, err := db.ConnectDB()
	if(err != nil){
		panic(err)
	}


	genreRepository := repository.NewGenreRepository(dbConnection)
	genreUseCase := usecase.NewGenreUseCase(genreRepository)
	genreController := controller.NewGenreController(genreUseCase)

	genres := gr.Group("/generos")
	{
		genres.GET("/", genreController.GetGenres)
		genres.POST("/", genreController.CreateGenre)
		genres.DELETE("/", genreController.DeleteGenre)
	}	
}