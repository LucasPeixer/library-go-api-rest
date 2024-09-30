package usecase

import (
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
