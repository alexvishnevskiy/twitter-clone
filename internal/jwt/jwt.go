package jwt

import (
	"fmt"
	"github.com/alexvishnevskiy/twitter-clone/internal/types"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// TODO: put secret key to .env
var jwtKey = []byte("my_secret_key")

type Claims struct {
	UserId types.UserId `json:"user_id"`
	jwt.StandardClaims
}

func GenerateJWT(id types.UserId) (string, error) {
	var err error

	// Create the Claims
	claims := &Claims{
		UserId: id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
			Issuer:    "test",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)

	if err != nil {
		fmt.Errorf("something went wrong: %s", err.Error())
		return "", err
	}

	return tokenString, nil
}

func ValidateMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			user, err := strconv.Atoi(r.FormValue("user_id"))
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			userID := types.UserId(user)

			// Get the JWT string from the header
			tokenString := r.Header.Get("Authorization")

			// Check if the token is empty
			if tokenString == "" {
				http.Error(w, "Authorization header is missing", http.StatusUnauthorized)
				return
			}

			// The token always starts with "Bearer "
			// we need to remove this part in order to be able to parse the token correctly
			tokenString = strings.TrimPrefix(tokenString, "Bearer ")

			// Parse the token
			claims := &Claims{}
			token, err := jwt.ParseWithClaims(
				tokenString, claims, func(token *jwt.Token) (interface{}, error) {
					return jwtKey, nil
				},
			)

			// If the token is expited, redirect to /login
			ve, ok := err.(*jwt.ValidationError)
			if ok && (ve.Errors&jwt.ValidationErrorExpired != 0) {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}
			if err != nil || !token.Valid {
				http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
				return
			}

			if claims.UserId != userID {
				http.Error(w, "Invalid user_id", http.StatusUnauthorized)
				return
			}

			// If everything is OK, call the next handler
			next.ServeHTTP(w, r)
		},
	)
}
