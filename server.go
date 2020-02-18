package main

import (
	"database/sql"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/jwtauth"
	"github.com/go-chi/render"
	"github.com/minheq/kedul_server_main/auth"
	"github.com/minheq/kedul_server_main/errors"
	"github.com/minheq/kedul_server_main/logger"
	"github.com/minheq/kedul_server_main/phone"
)

type server struct {
	db        *sql.DB
	router    *chi.Mux
	logger    *logger.Logger
	smsSender phone.SMSSender
}

func newServer(
	db *sql.DB,
	router *chi.Mux,
	logger *logger.Logger,
	smsSender phone.SMSSender,
) *server {
	s := &server{
		db:        db,
		router:    router,
		logger:    logger,
		smsSender: smsSender,
	}

	s.routes()

	return s
}

func (s *server) routes() {
	// dependencies
	tokenAuth := jwtauth.New("HS256", []byte("secret"), nil)
	authStore := auth.NewStore(s.db)
	authService := auth.NewService(authStore, tokenAuth, s.smsSender)

	// middlewares
	s.router.Use(middleware.RequestID)
	s.router.Use(middleware.RealIP)
	s.router.Use(logger.NewRequestLogger(s.logger))
	s.router.Use(middleware.Recoverer)
	s.router.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "Workspace", "X-CSRF-Token"},
		ExposedHeaders:   []string{""},
		AllowCredentials: true,
		MaxAge:           300,
	}).Handler)

	// public handlers
	s.router.Group(func(r chi.Router) {
		s.router.Post("/auth/login_verify", s.handleLoginVerify(authService))
		s.router.Post("/auth/login_check", s.handleLoginCheck(authService))
	})

	// protected handlers
	s.router.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(tokenAuth))
		r.Use(jwtauth.Authenticator)
		r.Use(addCurrentUserContext(authService))

		r.Get("/auth/current_user", s.handleGetCurrentUser(authService))
		r.Post("/auth/update_phone_number_verify", s.handleUpdatePhoneNumberVerify(authService))
		r.Post("/auth/update_phone_number_check", s.handleUpdatePhoneNumberCheck(authService))
		r.Post("/auth/update_user_profile", s.handleUpdateUserProfile(authService))
	})
}

func (s *server) respondError(w http.ResponseWriter, r *http.Request, err error) {
	s.logger.Error((err))
	_ = render.Render(w, r, errors.NewErrResponse(err))
}
