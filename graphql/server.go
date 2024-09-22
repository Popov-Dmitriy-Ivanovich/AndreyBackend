package main

import (
	"backend/models/products"
	"backend/models/discounts"
	"backend/models/productmedia"
	"graphql/graph"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"gorm.io/gorm"
	"gorm.io/driver/sqlite"
)

const defaultPort = "8080"
const DBPATH = "server.db"
func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}
	resolver := graph.Resolver{}
	db, err := gorm.Open(sqlite.Open("server.db"), &gorm.Config{})
	db.Create(&products.Product{Name: "TestProduct", Price: 123.123})
	if err != nil {
		panic("Db is not open")
	}
	resolver.Db = db
	resolver.Db.AutoMigrate(
		&products.Product{},
		&discounts.Discount{},
		&productmedia.ProductMedia{},
		&products.Collection{},
		&products.Category{},
		&products.Cart{},
	)

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &resolver}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
