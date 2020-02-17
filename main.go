package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/jwtauth"
	_ "github.com/lib/pq"
	"github.com/minheq/kedul_server_main/auth"
	"github.com/minheq/kedul_server_main/logger"
	"github.com/minheq/kedul_server_main/sms"
	"github.com/sirupsen/logrus"
)

func main() {
	r := chi.NewRouter()
	l := logger.NewLogger()
	dbURL := os.Getenv("DATABASE_URL")
	db, err := sql.Open("postgres", dbURL)

	if err != nil {
		l.WithFields(logrus.Fields{
			"DATABASE_URL": dbURL,
			"error":        err.Error(),
		}).Fatal("error opening database")
	}

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(logger.NewRequestLogger(l))
	r.Use(middleware.Recoverer)

	r.Use(cors.New(cors.Options{
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
	authService := auth.NewService(authStore, tokenAuth, smsSender, l)

	r.Use(jwtauth.Verifier(tokenAuth))

	r.Post("/login_verify", HandleLoginVerify(authService))
	r.Post("/login_verify_check", HandleLoginVerifyCheck(authService))
	r.Get("/current_user", HandleGetCurrentUser(authService))

	fmt.Println("Server listening at localhost:4000")

	http.ListenAndServe(":4000", r)
}
