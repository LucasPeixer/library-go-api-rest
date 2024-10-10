package main

import (
	"go-api/routes"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {	

		
	router := gin.Default()
	routes.RegisterRoutes(router)
	
	router.Run(":8000")
}