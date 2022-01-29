package main

import (
	"fmt"
	"ksemilla/graph"
	"ksemilla/graph/generated"
	"ksemilla/middlewares"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"

	"github.com/joho/godotenv"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	godotenv.Load(filepath.Join(".", ".env"))
	// err := godotenv.Load(filepath.Join(".", ".env"))
	// if err != nil {
	// 	log.Fatal("Error loading .env file")
	// }

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middlewares.GetCorsHandler())
	r.Use(middlewares.UserContextBody)

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))

	// http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	// http.Handle("/query", srv)

	r.Handle("/", playground.Handler("GraphQL playground", "/query"))
	r.Handle("/query", srv)

	test := os.Getenv("ACCESS_KEY")
	test2 := os.Getenv("AWS_ACCESS_KEY")
	fmt.Println(test, "test", test2)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	// log.Fatal(http.ListenAndServe(":"+port, nil))
	log.Fatal(http.ListenAndServe(":"+port, r))
}
