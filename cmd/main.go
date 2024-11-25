package main

import (
	"go-api/initializers"
	"go-api/middleware"
	"go-api/routes"
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func init() {
	initializers.LoadEnv()
	initializers.InitDB()
}

func main() {
	defer initializers.DB.Close()

	r := gin.Default()
	r.Use(middleware.CORSMiddleware())
	routes.Routes(r)

	log.Fatal(r.Run(":" + initializers.Port))
}
