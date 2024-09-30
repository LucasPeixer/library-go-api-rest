package controller

import (
	"go-api/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type genreController struct {
	genreUseCase usecase.GenreUsecase
}

func NewGenreController(usecase usecase.GenreUsecase) genreController{
	return genreController{
		genreUseCase: usecase,
	}
}

func (g *genreController) GetGenres(ctx *gin.Context){

	genres, err := g.genreUseCase.GetGenres()

	if(err != nil){
		ctx.JSON(http.StatusInternalServerError, err)
	}
	ctx.JSON(http.StatusOK, genres)
}