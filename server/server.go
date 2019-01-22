package main

import (
	log "log"
	http "net/http"
	os "os"

	"github.com/rs/cors"

	"github.com/bouk/httprouter"
	"github.com/ossn/ossn-backend/models"
	"github.com/qor/session/manager"

	handler "github.com/99designs/gqlgen/handler"
	ossn_backend "github.com/ossn/ossn-backend"
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

	// mux.GET(prefix+"/oidc/callback", controllers.HandleOAuth2Callback)
	// mux.GET(prefix+"/oidc/login", controllers.HandleRedirect)
	// //TODO: Remove this once Mozilla is ready
	// mux.GET("/oidc/callback", controllers.HandleOAuth2Callback)
	// mux.GET("/oidc/login", controllers.HandleRedirect)

	registerAll(mux, prefix+"/auth/*a", models.Auth.NewServeMux())
	registerAll(mux, prefix+"/admin/*f", adminMux)
	//TODO: Remove this once migration is done
	registerAll(mux, "/auth/*a", models.Auth.NewServeMux())
	registerAll(mux, "/admin/*f", adminMux)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)

	sessionHandler := manager.SessionManager.Middleware(mux)
	corsHandler := cors.Default().Handler(sessionHandler)
	log.Fatal(http.ListenAndServe(":"+port, (corsHandler)))
}

func registerAll(mux *httprouter.Router, path string, handler http.Handler) {
	mux.Handler("GET", path, handler)
	mux.Handler("POST", path, handler)
	mux.Handler("PUT", path, handler)
	mux.Handler("DELETE", path, handler)
}
