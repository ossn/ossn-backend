package main

import (
	log "log"
	http "net/http"
	os "os"

	"github.com/ossn/ossn-backend/models"

	handler "github.com/99designs/gqlgen/handler"
	ossn_backend "github.com/ossn/ossn-backend"
)

const defaultPort = "8080"

func main() {
	defer models.DBSession.Close()
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	http.Handle("/", handler.Playground("GraphQL playground", "/query"))
	http.Handle("/query", handler.GraphQL(ossn_backend.NewExecutableSchema(ossn_backend.Config{Resolvers: &ossn_backend.Resolver{}})))

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
