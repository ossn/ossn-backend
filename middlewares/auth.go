package middlewares

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/ossn/ossn-backend/helpers"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/ossn/ossn-backend/models"
)

const AuthCookie = "sessionCookie"

var (
	jwtSecret []byte
)

func init() {
	secret := os.Getenv("JWT_SECRET")
	helpers.CheckEnvVariable(&secret, "JWT_SECRET")
	jwtSecret = []byte(secret)
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// c, err := r.Cookie(AuthCookie)
		// // Allow unauthenticated users in
		// if err != nil || c == nil {
		// 	next.ServeHTTP(w, r)
		// 	return
		// }
		token := r.Header.Get("X-Access-Token")
		if len(token) == 0 {
			next.ServeHTTP(w, r)
			return
		}

		session := &models.Session{}
		err := models.DBSession.Where("token = ?", token).First(session).Error
		if err != nil || !ValidateToken(&session.Token) {
			http.Error(w, "Invalid cookie or access token", http.StatusForbidden)
			return
		}

		user := &models.User{}
		err = models.DBSession.Where("id = ?", session.UserID).First(user).Error
		if err != nil {
			http.Error(w, "Invalid cookie or access token", http.StatusForbidden)
			return
		}

		// put it in context
		ctx := context.WithValue(r.Context(), helpers.UserCtxKey, user)
		ctx = context.WithValue(ctx, helpers.SessionCtxKey, session)

		// and call the next with our new context
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func SignToken(user *models.User) (string, error) {
	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	now := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    user.ID,
		"email": user.Email,
		"nbf":   now,
		"iat":   now,
		"exp":   now.Add(time.Minute * 15),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString(jwtSecret)

	return tokenString, err
}

func ValidateToken(tokenString *string) bool {
	// Parse takes the token string and a function for looking up the key. The latter is especially useful if you use multiple keys for your application.  The standard is to use 'kid' in the head of the token to identify which key to use, but the parsed token (head and claims) is provided to the callback, providing flexibility.
	token, err := jwt.Parse(*tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return jwtSecret, nil
	})
	if err != nil {
		return false
	}
	_, ok := token.Claims.(jwt.MapClaims)
	return ok && token.Valid
}
