package main

import (
	"embed"
	"halalcloud/handlers"
	"halalcloud/routes"
	"net/http"

	"github.com/gin-gonic/gin"
)

var content embed.FS

func main() {
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		data, err := content.ReadFile("templates/index.html")
		if err != nil {
			c.String(http.StatusInternalServerError, "Error loading index.html")
			return
		}
		// 返回 HTML 文件的内容
		c.Data(http.StatusOK, "text/html", data)
	})
	r.GET("/login", handlers.Login)
	routes.SetupRoutes(r)
	r.Run(":8080")
}
