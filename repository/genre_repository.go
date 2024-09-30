package repository

import (
	"database/sql"
	"fmt"
	"go-api/model"
)

type GenreRepository struct {
	connection *sql.DB
}

func NewGenreRepository (connection *sql.DB) GenreRepository {
	return GenreRepository{
		connection: connection,
	}
}

func (pr *GenreRepository) GetGenres() ([]model.Genre, error){
	query := "SELECT * FROM genres"
	rows, err := pr.connection.Query(query)
	if(err != nil){
		fmt.Println(err)
		return []model.Genre{}, err
	}

	var genreList []model.Genre
	var genreObj model.Genre

	for rows.Next(){
		err = rows.Scan(
			&genreObj.ID,
			&genreObj.Name)

		if(err != nil){
			fmt.Println(err)
		return []model.Genre{}, err
		}

		genreList = append(genreList, genreObj)
	}

	rows.Close()

	return genreList, nil
}