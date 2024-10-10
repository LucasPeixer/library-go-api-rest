package routes

import (
	"github.com/gin-gonic/gin"
)

// Criando o inicial de cada rota para concatenar os endpoints
func RegisterRoutes(router *gin.Engine){
	api := router.Group("/api/v1")
	BooksRouter(api)
}