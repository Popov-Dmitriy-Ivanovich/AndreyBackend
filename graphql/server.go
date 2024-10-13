package graphql

import (
	"backend/models/articles"
	"backend/models/discounts"
	"backend/models/productmedia"
	"backend/models/products"
	"backend/models/users"

	// "encoding/json"
	"graphql/graph"
	// "io"
	//"log"
	// "net/http"
	"os"
	"sync"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	// "github.com/99designs/gqlgen/graphql/playground"
	// "github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)
type MutexedData [T any] struct {
	Data T
	Mtx sync.Mutex
}
var globalSingletonServerDataInst *GlobalSingletonServerData = nil
var globalSingletonServerDataInstMutex sync.Mutex
type GlobalSingletonServerData struct {
	Db *gorm.DB
	Tokens MutexedData[[]string]
	JwtKey []byte
	GqlUser *handler.Server
	GqlAdmin *handler.Server
	Initialized bool
}

type Claims struct {
	Username string
	jwt.RegisteredClaims
}

func GetGlobalServerData() *GlobalSingletonServerData {
	globalSingletonServerDataInstMutex.Lock()
	defer globalSingletonServerDataInstMutex.Unlock()
	if globalSingletonServerDataInst == nil {
		globalSingletonServerDataInst = new(GlobalSingletonServerData)
	}
	return globalSingletonServerDataInst
}

func generateJWT() (string, error) {
	expirationTime := time.Now().Add(5 * time.Hour)
	claims := &Claims{
        Username: "username",
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(expirationTime),
        },
    }
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	globalData := GetGlobalServerData()
    return token.SignedString(globalData.JwtKey)
}

func InitGlobalServerData() error {
	const DEFAULTPORT= "8080"
	const DEFAULTDBPATH = "server.db"
	const DEFAULTJWTKEY = "JWTKEY"
	globalServerData := GetGlobalServerData()
	
	globalServerData.Tokens.Mtx.Lock()
	defer globalServerData.Tokens.Mtx.Unlock()
	globalServerData.Tokens.Data = make([]string, 0)
	port := os.Getenv("PORT")
	if port == "" {
		port = DEFAULTPORT
	}
	dbPath := os.Getenv("DBPATH")
	if dbPath == "" {
		dbPath = DEFAULTDBPATH
	}
	jwtKey := os.Getenv("JWTKEY")
	if jwtKey == "" {
		jwtKey = DEFAULTJWTKEY
	}

	globalServerData.JwtKey = []byte(jwtKey)
	dsn := "host=dev-pg-postgresql user=admin password=admin dbname=gorm port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}
	globalServerData.Db = db

	db.AutoMigrate(
		&users.User{},
		&products.Product{},
		&discounts.Discount{},
		&productmedia.ProductMedia{},
		&products.Collection{},
		&products.Category{},
		&products.Cart{},
		&products.Advert{},
		&products.Review{},
		&articles.Article{},
		&articles.ArticleMedia{},
		&discounts.Discount{},
		&productmedia.ProductMedia{},		
	)

	userResolver := graph.Resolver{Db: db, IsAdmin: false}
	adminResolver := graph.Resolver{Db: db, IsAdmin: false}

	globalServerData.GqlUser = handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &userResolver}))
	globalServerData.GqlAdmin = handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &adminResolver}))
	globalServerData.Initialized = true
	return nil
}

const defaultPort = "8080"
const DBPATH = "server.db"
// func main() {
// 	port := os.Getenv("PORT")
// 	if port == "" {
// 		port = defaultPort
// 	}
// 	resolver := graph.Resolver{}
// 	db, err := gorm.Open(sqlite.Open("server.db"), &gorm.Config{})
// 	db.Create(&products.Product{Name: "TestProduct", Price: 123.123})
// 	if err != nil {
// 		panic("Db is not open")
// 	}
// 	resolver.Db = db
// 	resolver.IsAdmin = false
// 	resolver.Db.AutoMigrate(
// 		&products.Product{},
// 		&discounts.Discount{},
// 		&productmedia.ProductMedia{},
// 		&products.Collection{},
// 		&products.Category{},
// 		&products.Cart{},
// 	)
// 	// adminResolver := graph.Resolver{Db: resolver.Db, IsAdmin: true}
// 	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &resolver}))
	
// 	// http.Handle("/", playground.Handler("GraphQL playground", "/query"))
// 	// http.Handle("/query", srv)

// 	// log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
// 	// log.Fatal(http.ListenAndServe(":"+port, nil))

// 	r := gin.Default()
// 	r.POST("/login", func (c *gin.Context){
// 		reqBody := c.Request.Body
// 		type LoginData struct {
// 			Login string
// 			Password string
// 		}
// 		if err != nil {
// 			c.JSON(http.StatusTeapot, gin.H{"error": "Can't parse request body"})
// 		}
// 		var body LoginData
// 		rowBody, errReader := io.ReadAll(reqBody)
// 		if errReader != nil {
// 			c.JSON(http.StatusTeapot, gin.H{"error": "Can't parse request body"})
// 		}
// 		json.Unmarshal(rowBody, &body)
// 		c.JSON(http.StatusOK, gin.H{"login": body.Login})
// 	})
// 	r.Run()
// }
