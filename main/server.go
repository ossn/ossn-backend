package main

import (
	"fmt"
	log "log"
	http "net/http"
	os "os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/ossn/ossn-backend/resolvers"

	handler "github.com/99designs/gqlgen/handler"
	"github.com/bouk/httprouter"
	"github.com/ossn/ossn-backend/controllers"
	"github.com/ossn/ossn-backend/middlewares"
	"github.com/ossn/ossn-backend/models"
	"github.com/qor/admin"
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

	// GraphQL route
	mux.GET(prefix+"/", handler.Playground("GraphQL playground", prefix+"/query"))
	mux.POST(prefix+"/query", handler.GraphQL(resolvers.NewExecutableSchema(resolvers.Config{Resolvers: &resolvers.Resolver{}})))

	// Open ID Connect routes
	mux.GET("/oidc/callback", controllers.HandleOAuth2Callback)
	mux.GET(prefix+"/oidc/callback", controllers.HandleOAuth2Callback)
	mux.GET(prefix+"/oidc/login", controllers.HandleRedirect)

	// Admin Routes
	adminMux := http.NewServeMux()
	models.AdminResource.MountTo("/admin", adminMux)
	if strings.EqualFold(os.Getenv("ENV"), "DEV") {
		models.AdminResource.GetRouter().GetMiddleware("csrf_check").Handler = func(context *admin.Context, middleware *admin.Middleware) { middleware.Next(context) }
	}
	registerAll(mux, "/admin/*f", middlewares.BasicAuth(adminMux))

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)

	middlewareHandler := middlewares.AuthMiddleware(mux)

	middlewareHandler = cors.New(cors.Options{
		AllowCredentials: true,
		AllowedMethods:   []string{"OPTIONS", "POST"},
		AllowedOrigins:   []string{"http://localhost:8000", "https://dev.ossn.club", "https://ossn.club"},
		AllowedHeaders:   []string{"X-Access-Token", "Content-Type"},
		ExposedHeaders:   []string{},
	}).Handler(middlewareHandler)

	// Gracefull shutdown
	var gracefulStop = make(chan os.Signal)
	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)
	go func() {
		<-gracefulStop
		controllers.CancelOpenIDContext()
		err := models.DBSession.Close()
		if err != nil {
			fmt.Println(err)
		}
		err = models.RedisClient.Close()
		if err != nil {
			fmt.Println(err)
		}
		os.Exit(0)
	}()

	log.Fatal(http.ListenAndServe(":"+port, (middlewareHandler)))
}

func registerAll(mux *httprouter.Router, path string, handler http.Handler) {
	mux.Handler("GET", path, handler)
	mux.Handler("POST", path, handler)
	mux.Handler("PUT", path, handler)
	mux.Handler("DELETE", path, handler)
}
