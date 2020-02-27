package main

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/jwtauth"
	"github.com/go-chi/render"
	"github.com/minheq/kedul_server_main/app"
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

	// auth
	authStore := auth.NewStore(s.db)
	authService := auth.NewService(authStore, tokenAuth, s.smsSender)

	// app
	businessStore := app.NewBusinessStore(s.db)
	locationStore := app.NewLocationStore(s.db)
	employeeStore := app.NewEmployeeStore(s.db)
	employeeRoleStore := app.NewEmployeeRoleStore(s.db)
	businessService := app.NewBusinessService(businessStore, locationStore, employeeStore)
	locationService := app.NewLocationService(businessStore, locationStore, employeeStore, employeeRoleStore)
	permissionService := app.NewPermissionService(employeeRoleStore, employeeStore)
	// employeeService := app.NewEmployeeService(employeeStore)
	// employeeRoleService := app.NewEmployeeRoleService(employeeStore, employeeRoleStore)

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
		r.Use(s.authenticate)
		r.Use(s.addCurrentUserContext(authService))

		r.Get("/auth/current_user", s.handleGetCurrentUser(authService))
		r.Post("/auth/update_phone_number_verify", s.handleUpdatePhoneNumberVerify(authService))
		r.Post("/auth/update_phone_number_check", s.handleUpdatePhoneNumberCheck(authService))
		r.Post("/auth/update_user_profile", s.handleUpdateUserProfile(authService))

		r.Get("/users/{userID}/businesses", s.handleGetBusinessesByUserID(businessService))
		r.Get("/users/{userID}/businesses/{businessID}/locations", s.handleGetLocationsByUserIDAndBusinessID(locationService))

		r.Post("/businesses", s.handleCreateBusiness(businessService))
		r.Get("/businesses/{businessID}", s.handleGetBusiness(businessService))
		r.Post("/businesses/{businessID}", s.handleUpdateBusiness(businessService))
		r.Delete("/businesses/{businessID}", s.handleDeleteBusiness(businessService))

		r.Post("/locations", s.handleCreateLocation(locationService))
		r.Post("/locations/{locationID}", s.handleUpdateLocation(locationService, permissionService))
		r.Get("/locations/{locationID}", s.handleGetLocation(locationService, permissionService))
		r.Delete("/locations/{locationID}", s.handleDeleteLocation(locationService))
	})
}

func (s *server) respondError(w http.ResponseWriter, r *http.Request, err error) {
	s.logger.Error((err))
	render.Render(w, r, errors.NewErrResponse(err))
}

func (s *server) respondSuccess(w http.ResponseWriter, r *http.Request, v render.Renderer) {
	render.Render(w, r, v)
}

func (s *server) decode(w http.ResponseWriter, r *http.Request, v interface{}) error {
	err := json.NewDecoder(r.Body).Decode(&v)

	return err
}
