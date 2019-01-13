package main

import (
	log "log"
	http "net/http"
	os "os"

	"github.com/ossn/ossn-backend/models"
	"github.com/qor/session/manager"

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
	mux.Handle("/auth/", models.Auth.NewServeMux())
	models.AdminResource.MountTo("/admin", mux)
	// mux.Handle("/auth/", 	 .NewServeMux())
	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	// handler := cors.Default().Handler(mux)
	//  manager.SessionManager.Middleware(mux)
	log.Fatal(http.ListenAndServe(":"+port, manager.SessionManager.Middleware(mux)))
}
