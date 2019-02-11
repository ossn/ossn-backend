package main

import (
	log "log"
	http "net/http"
	os "os"

	handler "github.com/99designs/gqlgen/handler"
	"github.com/bouk/httprouter"
	ossn_backend "github.com/ossn/ossn-backend"
	"github.com/ossn/ossn-backend/controllers"
	"github.com/ossn/ossn-backend/middlewares"
	"github.com/ossn/ossn-backend/models"
	"github.com/rs/cors"
)

const (
	prefix      = "/api/v1.0"
	defaultPort = "8080"
)

func main() {
	defer models.DBSession.Close()
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	mux := httprouter.New()

	adminMux := http.NewServeMux()
	models.AdminResource.MountTo(prefix+"/admin", adminMux)

	//TODO: Remove this once migration is done
	models.AdminResource.MountTo("/admin", adminMux)

	//TODO: Remove this once migration is done
	mux.GET("/", handler.Playground("GraphQL playground", prefix+"/query"))
	mux.POST("/query", handler.GraphQL(ossn_backend.NewExecutableSchema(ossn_backend.Config{Resolvers: &ossn_backend.Resolver{}})))

	mux.GET(prefix+"/", handler.Playground("GraphQL playground", prefix+"/query"))
	mux.POST(prefix+"/query", handler.GraphQL(ossn_backend.NewExecutableSchema(ossn_backend.Config{Resolvers: &ossn_backend.Resolver{}})))

	mux.GET(prefix+"/oidc/callback", controllers.HandleOAuth2Callback)
	mux.GET(prefix+"/oidc/login", controllers.HandleRedirect)
	//TODO: Remove this once Mozilla is ready
	mux.GET("/oidc/callback", controllers.HandleOAuth2Callback)
	mux.GET("/oidc/login", controllers.HandleRedirect)

	registerAll(mux, prefix+"/admin/*f", adminMux)
	//TODO: Remove this once migration is done
	registerAll(mux, "/admin/*f", adminMux)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)

	middlewareHandler := middlewares.AuthMiddleware(mux)

	middlewareHandler = cors.New(cors.Options{
		AllowCredentials: true,
		AllowedMethods:   []string{"OPTIONS", "POST"},
		AllowedOrigins:   []string{"http://localhost:8000", "https://dev.ossn.club", "https://ossn.club"},
		AllowedHeaders:   []string{"X-Access-Token", "Content-Type"},
		ExposedHeaders:   []string{},
	}).Handler(middlewareHandler)

	log.Fatal(http.ListenAndServe(":"+port, (middlewareHandler)))
}

func registerAll(mux *httprouter.Router, path string, handler http.Handler) {
	mux.Handler("GET", path, handler)
	mux.Handler("POST", path, handler)
	mux.Handler("PUT", path, handler)
	mux.Handler("DELETE", path, handler)
}
