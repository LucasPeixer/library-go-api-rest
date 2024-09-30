package controller

import (
	"go-api/model"
	"go-api/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type bookController struct {
	bookUseCase usecase.BookUsecase
}

func NewBookController(usecase usecase.BookUsecase) bookController {
	return bookController{
		bookUseCase: usecase,
	}
}

func (b *bookController) GetBooks(ctx *gin.Context){
	books, err := b.bookUseCase.GetBooks()
	if(err != nil){
		ctx.JSON(http.StatusInternalServerError, err)
	}
	ctx.JSON(http.StatusOK, books)
}

func (bc *bookController) CreateBook(ctx *gin.Context) {
	var newBook model.Book

	// Binda o JSON da requisição para o Model (Direciona corretamente os dados para os campos do Model)
	if err := ctx.ShouldBindJSON(&newBook); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
	}

	lastInsertID, err := bc.bookUseCase.CreateBook(newBook)
	if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao registrar o livro"})
			return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Livro registrado com sucesso", "book_id": lastInsertID})
}
