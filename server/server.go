package main

import (
	log "log"
	http "net/http"
	os "os"

	"github.com/bouk/httprouter"
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

	mux := httprouter.New()
	mux.GET("/", handler.Playground("GraphQL playground", "/query"))
	mux.POST("/query", handler.GraphQL(ossn_backend.NewExecutableSchema(ossn_backend.Config{Resolvers: &ossn_backend.Resolver{}})))
	registerAll(mux, "/auth/*a", models.Auth.NewServeMux())
	adminMux := http.NewServeMux()
	models.AdminResource.MountTo("/admin", adminMux)
	registerAll(mux, "/admin/*f", adminMux)
	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	// handler := cors.Default().Handler(mux)
	//  manager.SessionManager.Middleware(mux)
	log.Fatal(http.ListenAndServe(":"+port, manager.SessionManager.Middleware(mux)))
}

func registerAll(mux *httprouter.Router, path string, handler http.Handler) {
	mux.Handler("GET", path, handler)
	mux.Handler("POST", path, handler)
	mux.Handler("PUT", path, handler)
	mux.Handler("DELETE", path, handler)
}
