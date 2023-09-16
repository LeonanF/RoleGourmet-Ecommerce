package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	// Crie um novo servidor Gin
	server := gin.New()

	// Configure o roteamento para servir o arquivo "index.html" como página principal
	server.GET("/", func(c *gin.Context) {
		c.File("static/index.html")
	})

	// Configure o roteamento para servir arquivos estáticos da pasta "static"
	server.Static("/static", "./static")

	// Roteie a página "produtos.html"
	server.GET("/produtos.html", func(c *gin.Context) {
		c.File("static/produtos.html")
	})

	// Inicia o servidor
	server.Run()
}
