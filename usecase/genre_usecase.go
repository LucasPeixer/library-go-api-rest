package usecase

import (
	"fmt"
	"go-api/model"
	"go-api/repository"
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
