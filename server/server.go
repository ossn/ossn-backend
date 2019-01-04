package main

import (
	log "log"
	http "net/http"
	os "os"

	"github.com/ossn/ossn-backend/models"
	"github.com/rs/cors"

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

	mux := http.NewServeMux()
	mux.Handle("/", handler.Playground("GraphQL playground", "/query"))
	mux.Handle("/query", handler.GraphQL(ossn_backend.NewExecutableSchema(ossn_backend.Config{Resolvers: &ossn_backend.Resolver{}})))

	models.Admin.MountTo("/admin", mux)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	handler := cors.Default().Handler(mux)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}
