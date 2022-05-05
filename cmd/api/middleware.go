package main

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/pascaldekloe/jwt"
)

func (app *application) enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, OPTIONS, DELETE")
		next.ServeHTTP(w, r)
	})
}

func (app *application) checkToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Authorization")
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			// set an anonymous user
		}

		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 {
			app.errorJSON(w, errors.New("Invalid auth header"))
			return
		}
		// if headerParts[0] != "Bearer" {
		// 	app.errorJSON(w, errors.New("Unauthorized user"))
		// 	return
		// }
		token := headerParts[1]
		claims, err := jwt.HMACCheck([]byte(token), []byte(app.config.jwt.secret))
		if err != nil {
			app.errorJSON(w, errors.New("Unauthorized user"))
			return
		}
		if !claims.Valid(time.Now()) {
			app.errorJSON(w, errors.New("Token expired"))
			return
		}
		if !claims.AcceptAudience("mydomain.com") {
			app.errorJSON(w, errors.New("Invalid Token"))
			return
		}
		if claims.Issuer != "mydomain.com" {
			app.errorJSON(w, errors.New("Invalid user"))
			return
		}

		userID, err := strconv.ParseInt(claims.Subject, 10, 64)
		if err != nil {
			app.errorJSON(w, errors.New("Unauthorized"))
			return
		}
		log.Println("Valid user", userID)

		next.ServeHTTP(w, r)
	})
}
