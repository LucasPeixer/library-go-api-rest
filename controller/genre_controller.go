package controller

import (
	"fmt"
	"go-api/model"
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

func (g *genreController) CreateGenre(ctx *gin.Context){
	var newGenre model.Genre

	err := ctx.BindJSON(&newGenre)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	lastGenreCreated, err := g.genreUseCase.CreateGenre(newGenre)
	if (err != nil){
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao registrar o Genero"})
	}
	ctx.JSON(http.StatusOK, lastGenreCreated)
}

func (g *genreController) DeleteGenre (ctx *gin.Context){
	lastGenereDeleted, err := g.genreUseCase.DeleteGenre(ctx)

	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao deletar o Livro!"})
		return 
	}

	ctx.JSON(http.StatusOK,gin.H{"message": "Genero deletado com sucesso!", "Genero deletado": lastGenereDeleted})

}