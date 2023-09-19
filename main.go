package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	mongoClient *mongo.Client
	collection  *mongo.Collection
)

type Config struct {
	MongoDBURI string `json:"MONGODB_URI"`
}

type Produtos struct {
	Name        string `json:"name"`
	Value       string `json:"value"`
	Category    string `json:"category"`
	Description string `json:"description"`
}

func main() {
	createServer()
}

func createServer() {
	// Crie um novo servidor Gin
	server := gin.New()

	server.LoadHTMLGlob("static/*.html")

	// Configure o roteamento para servir o arquivo "index.html" como página principal
	server.GET("/", func(c *gin.Context) {
		c.File("static/index.html")
	})

	// Configure o roteamento para servir arquivos estáticos da pasta "static"
	server.Static("/static", "./static")

	connectToDataBase()
	defer mongoClient.Disconnect(context.Background())

	server.GET("/produtos.html", func(c *gin.Context) {
		// Crie uma consulta para recuperar todos os usuários.
		filter := bson.D{}

		// Execute a consulta.
		cursor, err := collection.Find(context.Background(), filter)
		if err != nil {
			c.AbortWithStatus(500)
			return
		}

		var produtos []Produtos

		// Itere pelos documentos no cursor.
		for cursor.Next(context.Background()) {
			var produto Produtos
			if err := cursor.Decode(&produto); err != nil {
				c.AbortWithStatus(500)
				return
			}
			produtos = append(produtos, produto)
		}
		c.HTML(http.StatusOK, "produtos.html", gin.H{
			"Produtos": produtos,
		})
	})

	// Inicia o servidor
	server.Run()
}

func connectToDataBase() {

	// Abra o arquivo de configuração
	configFile, err := os.Open("config.json")
	if err != nil {
		fmt.Println("Erro ao abrir o arquivo de configuração:", err)
		return
	}
	defer configFile.Close()

	// Decode do arquivo JSON para a estrutura Config
	var config Config
	decoder := json.NewDecoder(configFile)
	err = decoder.Decode(&config)
	if err != nil {
		fmt.Println("Erro ao decodificar o arquivo de configuração:", err)
		return
	}

	// Defina a variável de ambiente localmente
	os.Setenv("MONGODB_URI", config.MongoDBURI)

	// Agora você pode acessar a variável de ambiente normalmente
	mongodbURI := os.Getenv("MONGODB_URI")

	client, err := mongo.NewClient(options.Client().ApplyURI(mongodbURI))
	if err != nil {
		panic(err)
	}

	ctx := context.TODO()
	err = client.Connect(ctx)
	if err != nil {
		panic(err)
	}

	mongoClient = client
	collection = client.Database("rolegourmet").Collection("produtos")
}
