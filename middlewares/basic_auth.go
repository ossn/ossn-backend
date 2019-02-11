package middlewares

import (
	"net/http"
	"os"
)

var (
	requiredUser     = os.Getenv("ADMIN_USER")
	requiredPassword = os.Getenv("ADMIN_PASSWORD")
)

func BasicAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the Basic Authentication credentials
		user, password, hasAuth := r.BasicAuth()

		if hasAuth && user == requiredUser && password == requiredPassword {
			// Delegate request to the given handle
			next.ServeHTTP(w, r)
			return
		}

		// Request Basic Authentication otherwise
		w.Header().Set("WWW-Authenticate", "Basic realm=Restricted")
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)

	})
}
