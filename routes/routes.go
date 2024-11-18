package routes

import (
	"github.com/gin-gonic/gin"
)

// Routes registra todas as rotas http.
func Routes(r *gin.Engine) {
	api := r.Group("/api/v1")
	UserRoutes(api)
	BookRoutes(api)
	ReservationRoutes(api)
}
