package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"go-api/model"
	"strconv"
)

type BookRepository interface {
	CreateBook(title, synopsis string, authorId int, genreIds []int) (*model.Book, error)
	GetBooks(title, author string, genres []string) (*[]model.Book, error)
	GetBookById(id int) (*model.Book, error)
	UpdateBook(bookId int, title, synopsis string, authorId int) error
	DeleteBook(bookId int) error
	AddStock(code, bookId int) (*model.BookStock, error)
	GetStock(code *int, bookId int) (*[]model.BookStock, error)
	GetStockById(id int) (*model.BookStock, error)
	UpdateStockStatus(id int, status string) error
	RemoveStock(id int, bookId *int) error
}

type bookRepository struct {
	db *sql.DB
}

func NewBookRepository(db *sql.DB) BookRepository {
	return &bookRepository{db: db}
}

// CreateBook cria um novo livro no banco de dados e o retorna.
func (br *bookRepository) CreateBook(title, synopsis string, authorId int, genreIds []int) (*model.Book, error) {
	query := `INSERT INTO book (title, synopsis, fk_author_id) VALUES ($1, $2, $3) RETURNING id;`

	var bookId int
	err := br.db.QueryRow(query, title, synopsis, authorId).Scan(&bookId)
	if err != nil {
		return nil, fmt.Errorf("error creating book: %v", err)
	}

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

	// Busca pelo nome do autor
	authorQuery := ` SELECT name FROM author WHERE id = $1;`

	var authorName string
	err = br.db.QueryRow(authorQuery, authorId).Scan(&authorName)
	if err != nil {
		return nil, fmt.Errorf("error fetching author name: %v", err)
	}

	book := &model.Book{
		Id:       bookId,
		Title:    title,
		Synopsis: synopsis,
		Author:   &model.Author{Id: authorId, Name: authorName},
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
func (br *bookRepository) GetBooks(title, author string, genres []string) (*[]model.Book, error) {
	books := make([]model.Book, 0)

	query := `
	SELECT b.id         AS book_id,
	       b.title      AS book_title,
	       b.synopsis   AS book_synopsis,
           COALESCE(SUM(CASE WHEN bs.status = 'available' THEN 1 ELSE 0 END), 0) - 
           COALESCE(SUM(CASE WHEN r.status = 'pending' THEN 1 ELSE 0 END), 0) AS amount,
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
	LEFT JOIN 
	        book_stock bs ON b.id = bs.fk_book_id
	LEFT JOIN
	   reservation r ON b.id = r.fk_book_id
	WHERE 
		    1=1 -- Permite adicionar condições "AND"
    `

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

	query += `
	GROUP BY 
	       b.id, b.title, b.synopsis, g.id, g.name, a.id, a.name;`

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
			bookMap[bookId] = model.NewBook(bookId, bookTitle, bookSynopsis, nil, author, []model.Genre{})
			bookMap[bookId].Amount = bookAmount
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
	return &books, nil
}

func (br *bookRepository) GetBookById(id int) (*model.Book, error) {
	query := `
    SELECT b.id         AS book_id,
           b.title      AS book_title,
           b.synopsis   AS book_synopsis,
           COALESCE(SUM(CASE WHEN bs.status = 'available' THEN 1 ELSE 0 END), 0) - 
           COALESCE(SUM(CASE WHEN r.status = 'pending' THEN 1 ELSE 0 END), 0) AS amount,
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
    LEFT JOIN 
           book_stock bs ON b.id = bs.fk_book_id
    LEFT JOIN
           reservation r ON b.id = r.fk_book_id
    WHERE 
           b.id = $1
    GROUP BY 
           b.id, b.title, b.synopsis, g.id, g.name, a.id, a.name
    `

	rows, err := br.db.Query(query, id)
	if err != nil {
		return nil, fmt.Errorf("error querying book by ID: %v", err)
	}
	defer rows.Close()

	var book *model.Book
	var genres []model.Genre

	for rows.Next() {
		var bookId int
		var title, synopsis string
		var amount int
		var genreId, authorId *int
		var genreName, authorName *string

		// Scan the row into variables
		err := rows.Scan(&bookId, &title, &synopsis, &amount, &genreId, &genreName, &authorId, &authorName)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}

		// Initialize the book only once
		if book == nil {
			book = &model.Book{
				Id:       bookId,
				Title:    title,
				Synopsis: synopsis,
				Amount:   amount,
				Author: &model.Author{
					Id:   *authorId,
					Name: *authorName,
				},
				Genres: []model.Genre{},
			}
		}

		// Append genres if present
		if genreId != nil && genreName != nil {
			genres = append(genres, model.Genre{
				Id:   *genreId,
				Name: *genreName,
			})
		}
	}

	if book == nil {
		return nil, fmt.Errorf("book with ID %d not found", id)
	}

	book.Genres = genres
	return book, nil
}

// UpdateBook atualiza as informações de um livro existente.
func (br *bookRepository) UpdateBook(id int, title, synopsis string, authorId int) error {
	query := `
        UPDATE book
        SET title = $1, synopsis = $2, fk_author_id = $3
        WHERE id = $4
        RETURNING id;
    `

	var updatedBookId int
	err := br.db.QueryRow(query, title, synopsis, authorId, id).Scan(&updatedBookId)
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

func (br *bookRepository) AddStock(code, bookId int) (*model.BookStock, error) {
	query := `INSERT INTO book_stock (code, fk_book_id) VALUES ($1, $2) RETURNING id;`

	var bookStockId int
	err := br.db.QueryRow(query, code, bookId).Scan(&bookStockId)
	if err != nil {
		return nil, fmt.Errorf("error adding book with code '%d' to stock: %v", code, err)
	}
	var bookStock model.BookStock
	bookStock.Id = bookStockId
	bookStock.Code = code
	bookStock.BookId = bookId
	bookStock.Status = "available"
	return &bookStock, nil
}

func (br *bookRepository) GetStock(code *int, bookId int) (*[]model.BookStock, error) {
	query := `SELECT id, status, code FROM book_stock WHERE 1=1 AND fk_book_id = $1`

	var args []interface{}
	args = append(args, bookId)

	if code != nil {
		query += ` AND code = $2`
		args = append(args, code)
	}

	rows, err := br.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	bookStocks := make([]model.BookStock, 0)

	for rows.Next() {
		var bookStock model.BookStock

		err := rows.Scan(&bookStock.Id, &bookStock.Status, &bookStock.Code)
		if err != nil {
			return nil, err
		}
		bookStocks = append(bookStocks, bookStock)
	}

	return &bookStocks, nil
}

func (br *bookRepository) GetStockById(id int) (*model.BookStock, error) {
	query := `SELECT id, status, code FROM book_stock WHERE id = $1`
	var bookStock model.BookStock
	err := br.db.QueryRow(query, id).Scan(&bookStock.Id, &bookStock.Status, &bookStock.Code)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("book stock with id %d not found", id)
		}
		return nil, err
	}
	return &bookStock, nil
}

func (br *bookRepository) UpdateStockStatus(id int, status string) error {
	query := `
		UPDATE book_stock 
		SET status = $1 
		WHERE id = $2
		RETURNING id;
	`

	err := br.db.QueryRow(query, status, id).Scan(&id)
	if err != nil {
		return fmt.Errorf("error updating stock status: %v", err)
	}

	return nil
}

func (br *bookRepository) RemoveStock(id int, bookId *int) error {
	query := `
        DELETE FROM book_stock 
        WHERE id = $1
    `

	var args []interface{}
	args = append(args, id)

	if bookId != nil {
		query += ` AND fk_book_id = $2`
		args = append(args, bookId)
	}

	query += ` RETURNING id;`

	var deletedBookStockId int

	err := br.db.QueryRow(query, args...).Scan(&deletedBookStockId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("book stock with id %d not found", id)
		}
		return fmt.Errorf("error deleting book stock: %v", err)
	}

	return nil
}
