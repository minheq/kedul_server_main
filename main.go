package main

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/jwtauth"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

func main() {
	db, err := sql.Open("postgres", "postgres://postgres@127.0.0.1:5432/kedul?sslmode=disable")

	if err != nil {
		panic(err)
	}

	router := chi.NewRouter()

	// A good base middleware stack
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	// CORS
	corsConfiguration := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "Workspace", "X-CSRF-Token"},
		ExposedHeaders:   []string{""},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})

	tokenAuth := jwtauth.New("HS256", []byte("secret"), nil)

	router.Use(jwtauth.Verifier(tokenAuth))
	router.Use(corsConfiguration.Handler)

	logger := logrus.New()
	server := newServer(db, router, logger, tokenAuth)

	fmt.Println("Server listening at localhost:4000")
	http.ListenAndServe(":4000", server)
}
