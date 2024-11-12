package controller

import (
	"go-api/usecase"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type BookController interface {
	CreateBook(c *gin.Context)
	GetBooks(c *gin.Context)
	UpdateBook(c *gin.Context)
	DeleteBook(c *gin.Context)
}

type bookController struct {
	useCase usecase.BookUseCase
}

func NewBookController(useCase usecase.BookUseCase) BookController {
	return &bookController{useCase: useCase}
}

// CreateBook recebe um input JSON através do gin.Context e tenta criar um livro.
func (bc *bookController) CreateBook(c *gin.Context) {
	var i struct {
		Title    string `json:"title" binding:"required"`
		Synopsis string `json:"synopsis" binding:"required"`
		Amount   int    `json:"amount" binding:"required"`
		AuthorId int    `json:"author_id" binding:"required"`
		GenreIds []int  `json:"genre_ids" binding:"required"`
	}

	if err := c.ShouldBindJSON(&i); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book creation input"})
		return
	}

	book, err := bc.useCase.CreateBook(i.Title, i.Synopsis, i.Amount, i.AuthorId, i.GenreIds)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, book)
}

// GetBooks retorna os livros com e sem query params (title, author e genres).
func (bc *bookController) GetBooks(c *gin.Context) {
	title := c.Query("title")
	author := c.Query("author")
	genresParam := c.Query("genres")

	// Separa múltiplos gêneros por vírgula
	var genres []string
	if genresParam != "" {
		genres = strings.Split(genresParam, ",")
	}

	books, err := bc.useCase.GetBooks(title, author, genres)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, books)
}

// UpdateBook atualiza as informações de um livro existente.
func (bc *bookController) UpdateBook(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book Id"})
		return
	}

	var i struct {
		Title    string `json:"title" binding:"required"`
		Synopsis string `json:"synopsis" binding:"required"`
		Amount   int    `json:"amount" binding:"required"`
		AuthorId int    `json:"author_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&i); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book update input"})
		return
	}

	if err := bc.useCase.UpdateBook(id, i.Title, i.Synopsis, i.Amount, i.AuthorId); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Book updated successfully"})
}

// DeleteBook deleta um livro existente.
func (bc *bookController) DeleteBook(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book Id"})
		return
	}

	if err := bc.useCase.DeleteBook(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Book deleted successfully"})
}
