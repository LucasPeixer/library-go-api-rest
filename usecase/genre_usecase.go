package usecase

import (
	"fmt"
	"go-api/model"
	"go-api/repository"

	"github.com/gin-gonic/gin"
)

type GenreUsecase struct {
	repository repository.GenreRepository
}

func NewGenreUseCase(repo repository.GenreRepository) GenreUsecase {
	return GenreUsecase{
		repository: repo,
	}
}

func (gu *GenreUsecase) GetGenres() ([]model.Genre,error){
	return gu.repository.GetGenres()
}

func (gu *GenreUsecase) CreateGenre(genre model.Genre) (string, error){
	lastGenreCreated, err := gu.repository.CreateGenre(genre)
	if(err != nil) {
		fmt.Println(err)
		return "", err
	}

	return lastGenreCreated, nil
}

func (gu *GenreUsecase) DeleteGenre(c *gin.Context) (string, error){
	lastGenereDeleted, err := gu.repository.DeleteGenre(c)
	if err != nil{
		fmt.Println(err)
		return "", err
	}

	return lastGenereDeleted, nil
}
