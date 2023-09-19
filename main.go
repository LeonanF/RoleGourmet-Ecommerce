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
	// Cria um novo servidor Gin
	server := gin.New()

	//Configura o servidor para carregar/renderizar templates HTML
	server.LoadHTMLGlob("static/*.html")

	// Configura a rota raiz / para servir o arquivo "index.html" como página principal
	server.GET("/", func(c *gin.Context) {
		c.File("static/index.html")
	})

	// Configura a rota /static para servir arquivos estáticos da pasta "static"
	server.Static("/static", "./static")

	//Chama a função de conexão ao banco de dados
	connectToDataBase()
	//Faz com que o banco de dados desconecte somente quando o programa parar de executar
	defer mongoClient.Disconnect(context.Background())

	//Configura a rota no servidor para quando for feita uma solicitação GET para a rota especificada
	//Quando a solicitação for feita, ou seja, quando o usuário acessar o site com a rota /produtos.html, a função anônima será executada
	server.GET("/produtos.html", func(c *gin.Context) {

		//Cria uma variável de filtro vazia, para que todos os dados sejam retornados
		filter := bson.D{}

		//Executa a consulta e retorna os elementos para o cursor(ponteiro)
		//Chama o método find na collection (váriavel que representa a coleção de dados, definida globalmente)
		//Um contexto de fundo neutro é criado e usado na consulta junto ao filtro definido anteriormente
		cursor, err := collection.Find(context.Background(), filter)

		//Verifica se ocorreu algum erro durante a consulta e, caso sim, envia uma resposta HTTP com o código de erro e encerra a função
		if err != nil {
			c.AbortWithStatus(500)
			return
		}

		//Cria um slice do tipo Produtos (struct definido no início)
		//O slice irá armazenar a lista de produtos retornada pelo banco de dados
		var produtos []Produtos

		// Itera pelos documentos no cursor.
		//O método Next() retorna falso quando não houver mais documentos a serem lidos
		for cursor.Next(context.Background()) {

			//É criada uma váriavel produto do tipo Produtos (struct definido no início)
			//Não confundir com o slice produtos, produto armazena cada documento temporariamente, e produtos é o slice que vai acumular todas as variáveis "produto"
			var produto Produtos

			//O método decode é aplicado no cursor para decodificar os dados do BD e armazenar na váriavel produto.
			//Se a decodificação funcionar, não haverá erros e err será nulo. Do contrário, o erro será retornado para a váriavel err e o código condicional tratará a exceção
			/* Notas de estudo: O código a seguir equivale a:
				err := cursor.Decode(&produto);
				if err != nil {
					c.AbortWithStatus(500)
					return
				}
				No entanto, na forma realmente utilizada, a váriavel err tem escopo de bloco e não fica definida fora do if
			}
			*/
			if err := cursor.Decode(&produto); err != nil {
				c.AbortWithStatus(500)
				return
			}
			//A seguinte linha adiciona a lista de dados do produto atual ao slice de produtos anterior
			produtos = append(produtos, produto)
		}

		//O código a seguir gera uma resposta para o arquivo html produtos.html, definindo o status da resposta como tratada com sucesso e passando um map com os dados a serem acessados pelo arquivo html.
		c.HTML(http.StatusOK, "produtos.html", gin.H{
			"Produtos": produtos,
		})
	})

	//O código busca o valor da váriavel de ambiente PORT
	port := os.Getenv("PORT")

	//É feito um teste para definir se o valor da variável está vazio/variável não definida. Caso esteja, é atribuído um valor padrão à variável
	if port == "" {
		port = "8080"
	}

	// Inicia o servidor com a porta definida acima
	server.Run(":" + port)
}

func connectToDataBase() {

	if os.Getenv("MONGODB_URI") == "" {
		// Abra o arquivo config.json e em caso de erro, imprime o erro e encerra a função
		configFile, err := os.Open("config.json")
		if err != nil {
			fmt.Println("Erro ao abrir o arquivo de configuração:", err)
			return
		}

		//Garante que o arquivo seja fechado no fim da execução
		defer configFile.Close()

		//Criada váriavel do tipo Config (struct no início) para armazenar os dados do arquivo .json
		var config Config

		// É criado um novo decodificador para ler os dados do arquivo configFile
		decoder := json.NewDecoder(configFile)

		//O código a seguir é que vai, efetivamente, decodificar e converter o JSON em uma estrutura de dados do tipo Config
		//Para entender a criação de err como váriavel de bloco, verificar documentação da função createServer()
		//Em caso de erro, a função é imediatamente encerrada e o erro é impresso
		if err := decoder.Decode(&config); err != nil {
			fmt.Println("Erro ao decodificar o arquivo de configuração:", err)
			return
		}

		//Seta a variável de ambiente com o valor de MongoDBURI (propriedade do struct Config)
		os.Setenv("MONGODB_URI", config.MongoDBURI)
	}
	//É criada uma váriavel Go para armazenar o valor da váriavel de ambiente
	mongodbURI := os.Getenv("MONGODB_URI")

	//Primeiro é criado um cliente MongoDB, e usada a função Connect para conectar com o banco de dados
	//É fornecido um contexto de fundo neutro, e configurado a URI de conexão (com a váriavel criada acima)
	//Devido à importância da conexão com o BD, se ocorrer algum erro, o programa inteiro é encerrado
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mongodbURI))
	if err != nil {
		panic(err)
	}

	//O cliente criado é atribuído à váriavel global para que outras partes do código possam acessar e executar operações no banco de dados
	mongoClient = client

	//É especificada a coleção de documentos a ser trabalhado em cima.
	collection = client.Database("rolegourmet").Collection("produtos")

}
