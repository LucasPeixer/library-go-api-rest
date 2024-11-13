package repository

import (
	"database/sql"
	"fmt"
	"go-api/model"
	"strconv"
)

type BookRepository interface {
	CreateBook(title, synopsis string, amount, authorId int, genreIds []int) (*model.Book, error)
	GetBooks(title, author string, genres []string) ([]model.Book, error)
	UpdateBook(bookId int, title, synopsis string, amount, authorId int) error
	DeleteBook(bookId int) error
}

type bookRepository struct {
	db *sql.DB
}

func NewBookRepository(db *sql.DB) BookRepository {
	return &bookRepository{db: db}
}

// CreateBook cria um novo livro no banco de dados e o retorna.
func (br *bookRepository) CreateBook(title, synopsis string, amount, authorId int, genreIds []int) (*model.Book, error) {
	query := `
        INSERT INTO book (title, synopsis, amount, fk_author_id)
        VALUES ($1, $2, $3, $4)
        RETURNING id;
    `

	var bookId int
	err := br.db.QueryRow(query, title, synopsis, amount, authorId).Scan(&bookId)
	if err != nil {
		return nil, fmt.Errorf("error creating book: %v", err)
	}

	// Busca pelo nome do autor
	authorQuery := `
        SELECT name
        FROM author
        WHERE id = $1;
    `
	var authorName string
	err = br.db.QueryRow(authorQuery, authorId).Scan(&authorName)
	if err != nil {
		return nil, fmt.Errorf("error fetching author name: %v", err)
	}

	book := model.NewBook(
		bookId,
		title,
		synopsis,
		amount,
		&model.Author{Id: authorId, Name: authorName},
		[]model.Genre{},
	)

	// Se houver gêneros, associa o livro criado com o Id do gênero
	if len(genreIds) > 0 {
		genreQuery := `INSERT INTO book_genre (fk_book_id, fk_genre_id) VALUES ($1, $2);`
		for _, genreId := range genreIds {
			_, err := br.db.Exec(genreQuery, bookId, genreId)
			if err != nil {
				return nil, fmt.Errorf("error inserting book-genre link: %v", err)
			}
		}
	}

	// Busca os gêneros associados ao livro e os adiciona ao objeto.
	genreRepo := NewGenreRepository(br.db)
	for _, genreId := range genreIds {
		genre, err := genreRepo.GetGenreById(genreId)
		if err != nil {
			continue
		}
		book.Genres = append(book.Genres, *genre)
	}
	return book, nil
}

// GetBooks retorna livros com base em filtros de título, autor e gêneros (podem ser uma string vazia).
func (br *bookRepository) GetBooks(title, author string, genres []string) ([]model.Book, error) {
	books := make([]model.Book, 0)

	query := `
	SELECT b.id         AS book_id,
	       b.title      AS book_title,
	       b.synopsis   AS book_synopsis,
	       b.amount     AS book_amount,
	       g.id         AS genre_id,
	       g.name       AS genre_name,
	       a.id         AS author_id,
	       a.name       AS author_name
	FROM 
	       book b
	LEFT JOIN
	       book_genre bg ON b.id = bg.fk_book_id
	LEFT JOIN
	       genre g ON bg.fk_genre_id = g.id
	LEFT JOIN
	       author a ON b.fk_author_id = a.id
	WHERE 1=1` // Permite adicionar condições "AND"

	var args []interface{}

	// Aplica os filtros se não forem strings vazias
	if title != "" {
		query += ` AND b.title ILIKE $` + strconv.Itoa(len(args)+1)
		args = append(args, "%"+title+"%")
	}

	if author != "" {
		query += ` AND a.name ILIKE $` + strconv.Itoa(len(args)+1)
		args = append(args, "%"+author+"%")
	}

	if len(genres) > 0 {
		query += ` AND g.name IN (`
		// Adiciona um placeholder '$%d' para cada gênero
		for i := range genres {
			query += fmt.Sprintf("$%d", len(args)+1)
			if i < len(genres)-1 {
				query += ", "
			}
			args = append(args, genres[i])
		}
		query += `)`
	}

	rows, err := br.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	bookMap := make(map[int]*model.Book)

	for rows.Next() {
		var bookId int
		var bookTitle string
		var bookSynopsis string
		var bookAmount int
		var genreId, authorId *int
		var genreName, authorName *string

		err := rows.Scan(&bookId, &bookTitle, &bookSynopsis, &bookAmount, &genreId, &genreName, &authorId, &authorName)
		if err != nil {
			return nil, err
		}

		// Se o livro ainda não foi mapeado, cria uma nova instância
		if _, exists := bookMap[bookId]; !exists {
			var author *model.Author
			if authorId != nil {
				author = &model.Author{Id: *authorId, Name: *authorName}
			}
			bookMap[bookId] = model.NewBook(bookId, bookTitle, bookSynopsis, bookAmount, author, []model.Genre{})
		}

		// Adiciona o gênero ao livro correspondente
		if genreId != nil {
			bookMap[bookId].Genres = append(bookMap[bookId].Genres, model.Genre{
				Id:   *genreId,
				Name: *genreName,
			})
		}
	}

	// Converte o map em uma lista de livros
	for _, book := range bookMap {
		books = append(books, *book)
	}
	return books, nil
}

// UpdateBook atualiza as informações de um livro existente.
func (br *bookRepository) UpdateBook(id int, title, synopsis string, amount, authorId int) error {
	query := `
        UPDATE book
        SET title = $1, synopsis = $2, amount = $3, fk_author_id = $4
        WHERE id = $5
        RETURNING id;
    `

	var updatedBookId int
	err := br.db.QueryRow(query, title, synopsis, amount, authorId, id).Scan(&updatedBookId)
	if err != nil {
		return fmt.Errorf("error updating book: %v", err)
	}

	return nil
}

// DeleteBook deleta um livro do banco de dados.
func (br *bookRepository) DeleteBook(bookId int) error {
	query := `
        DELETE FROM book
        WHERE id = $1
        RETURNING id;
    `

	var deletedBookId int
	err := br.db.QueryRow(query, bookId).Scan(&deletedBookId)
	if err != nil {
		return fmt.Errorf("error deleting book: %v", err)
	}

	return nil
}
