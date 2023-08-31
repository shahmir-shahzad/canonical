package controllers

import (
	"github.com/dgrijalva/jwt-go"
	"log"
	"net/http"
	"strings"
)

func AuthenticateAndHandle(next http.HandlerFunc) http.HandlerFunc {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Our middleware logic goes here...
		authHeader := r.Header.Get("Authorization")

		headerValues := strings.Split(authHeader, " ")

		if len(headerValues) < 2 {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Token not found!"))
			return
		}

		tokenValue := headerValues[1]
		log.Printf("The authorization headers are: %+v", authHeader)

		token, parseErr := jwt.Parse(tokenValue, func(t *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if parseErr != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unable to parse token"))
			return
		}

		if !token.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Invalid Token"))
			return
		}

		next.ServeHTTP(w, r)
	})
}
