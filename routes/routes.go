package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/swaggo/gin-swagger"
    "github.com/swaggo/files"  // Pacote correto para servir os arquivos Swagger
    _ "go-api/docs" // Atualize conforme o nome do módulo no go.mod
)

// Routes registra todas as rotas http.
func Routes(r *gin.Engine) {
	api := r.Group("/api/v1")
	api.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("/api/v1/swagger.json")))
	api.StaticFile("/swagger.json", "./docs/swagger.json")
	UserRoutes(api)
	BookRoutes(api)
	ReservationRoutes(api)
}
	