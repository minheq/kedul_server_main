package main

import (
	"context"
	"net/http"

	"github.com/go-chi/jwtauth"
	"github.com/minheq/kedul_server_main/auth"
	"github.com/minheq/kedul_server_main/errors"
)

type contextKey struct {
	name string
}

var (
	userCtxKey = &contextKey{"user"}
)

func (s *server) requireAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		op := "middleware.requireAuthentication"
		token, _, err := jwtauth.FromContext(r.Context())

		if err != nil {
			s.respondError(w, r, errors.Unauthorized(op))
			return
		}

		if token == nil || !token.Valid {
			s.respondError(w, r, errors.Unauthorized(op))
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (s *server) addCurrentUserContext(authService auth.Service) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			op := "middleware.addCurrentUserContext"
			user, err := authService.GetCurrentUser(r.Context())

			if err != nil {
				s.respondError(w, r, errors.Unauthorized(op))
				return
			}

			ctx := context.WithValue(r.Context(), userCtxKey, user)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
