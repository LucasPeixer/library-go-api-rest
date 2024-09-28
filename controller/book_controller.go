package controller

import (
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