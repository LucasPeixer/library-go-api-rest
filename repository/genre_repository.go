package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"go-api/model"
)

type GenreRepository interface {
	GetGenreById(id int) (*model.Genre, error)
}

type genreRepository struct {
	db *sql.DB
}

func NewGenreRepository(db *sql.DB) GenreRepository {
	return &genreRepository{db: db}
}

func (gr *genreRepository) GetGenreById(id int) (*model.Genre, error) {
	query := `
        SELECT id, name
        FROM genre
        WHERE id = $1;
    `

	var genre model.Genre
	err := gr.db.QueryRow(query, id).Scan(&genre.Id, &genre.Name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("genre with id %d not found", id)
		}
		return nil, err
	}

	return &genre, nil
}
