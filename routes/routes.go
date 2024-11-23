package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/swaggo/gin-swagger"
    "github.com/swaggo/files"
    _ "go-api/docs"
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
	