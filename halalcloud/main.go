package main

import (
	"halalcloud/handlers"
	"halalcloud/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.File("templates/index.html")
	})
	r.GET("/login", handlers.Login)
	routes.SetupRoutes(r)
	r.Run(":8080")
}
