package middlewares

import (
	"context"
	"encoding/json"
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
	jwtSecret          []byte
	invalidJWTResponse []byte
)

func init() {
	secret := os.Getenv("JWT_SECRET")
	helpers.CheckEnvVariable(&secret, "JWT_SECRET")
	jwtSecret = []byte(secret)
	var err error
	invalidJWTResponse, err = json.Marshal(map[string]string{"error": "Invalid access token"})
	if err != nil {
		panic(err)
	}
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("X-Access-Token")
		if len(token) == 0 {
			next.ServeHTTP(w, r)
			return
		}

		userID, err := models.RedisClient.Get(token).Uint64()
		if err != nil || !ValidateToken(&token) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusForbidden)
			_, err := w.Write(invalidJWTResponse)
			if err != nil {
				fmt.Println(err)
			}
			return
		}

		user := &models.User{}
		err = models.DBSession.Where("id = ?", userID).First(user).Error
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusForbidden)
			_, err := w.Write(invalidJWTResponse)
			if err != nil {
				fmt.Println(err)
			}
			return
		}

		// put it in context
		ctx := context.WithValue(r.Context(), helpers.UserCtxKey, user)

		session := models.Session{
			UserID: user.ID,
			Token:  token,
		}
		ctx = context.WithValue(ctx, helpers.SessionCtxKey, &session)

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

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return false
	}

	now := time.Now()

	nbf, err := time.Parse(time.RFC3339, claims["nbf"].(string))
	if err != nil || now.Before(nbf) {
		return false
	}

	iat, err := time.Parse(time.RFC3339, claims["iat"].(string))
	if err != nil || now.Before(iat) {
		return false
	}

	exp, err := time.Parse(time.RFC3339, claims["exp"].(string))
	if err != nil || now.After(exp) {
		return false
	}

	return true

}
