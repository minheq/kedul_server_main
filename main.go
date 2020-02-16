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
	"github.com/minheq/kedul_server_main/auth"
	"github.com/minheq/kedul_server_main/sms"
	"github.com/sirupsen/logrus"
)

func main() {
	db, err := sql.Open("postgres", "postgres://postgres@127.0.0.1:5432/kedul?sslmode=disable")

	if err != nil {
		panic(err)
	}

	router := chi.NewRouter()

	// Setup the logger backend using sirupsen/logrus and configure
	// it to use a custom JSONFormatter. See the logrus docs for how to
	// configure the backend at github.com/sirupsen/logrus
	logger := logrus.New()
	logger.Formatter = &logrus.JSONFormatter{
		// disable, as we set our own
		DisableTimestamp: true,
		PrettyPrint:      true,
	}

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(NewStructuredLogger(logger))
	router.Use(middleware.Recoverer)

	router.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "Workspace", "X-CSRF-Token"},
		ExposedHeaders:   []string{""},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}).Handler)

	tokenAuth := jwtauth.New("HS256", []byte("secret"), nil)
	smsSender := sms.NewSender()

	authStore := auth.NewStore(db)
	authService := auth.NewService(authStore, tokenAuth, smsSender)

	router.Use(jwtauth.Verifier(tokenAuth))

	router.Post("/login_verify", HandleLoginVerify(authService))
	router.Post("/login_verify_check", HandleLoginVerifyCheck(authService))

	fmt.Println("Server listening at localhost:4000")

	http.ListenAndServe(":4000", router)
}
