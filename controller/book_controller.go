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

func (b *bookController) GetAllBooks(ctx *gin.Context){
	books, err := b.bookUseCase.GetAllBooks()
	if(err != nil){
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, books)
}

func (b *bookController) GetBooks(ctx *gin.Context){
	books, err := b.bookUseCase.GetBooks(ctx)

	if(err != nil){
		ctx.JSON(http.StatusInternalServerError, err)
	}

	ctx.JSON(http.StatusOK, books)
}

func (b *bookController) CreateBook(ctx *gin.Context) {
	var newBook model.Book

	err := ctx.BindJSON(&newBook)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	lastInsertID, err := b.bookUseCase.CreateBook(newBook)
	if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao registrar o livro"})
			return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Livro registrado com sucesso", "book_id": lastInsertID})
}

func (b *bookController) DeleteBook(ctx *gin.Context){
	lastBookDeleted, err := b.bookUseCase.DeleteBook(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao deletar o Livro"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Livro deletado com sucesso!", "Informações do Livro": lastBookDeleted})
}
