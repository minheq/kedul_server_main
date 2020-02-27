package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/dgrijalva/jwt-go"
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

func (s *server) addCurrentUserContext(authService auth.Service) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			const op = "server.addCurrentUserContext"
			user, err := authService.GetCurrentUser(r.Context())

			if err != nil {
				s.respondError(w, r, errors.Unauthorized(op, err))
				return
			}

			ctx := context.WithValue(r.Context(), userCtxKey, user)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func (s *server) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		const op = "server.authenticate"
		token, _ := r.Context().Value(jwtauth.TokenCtxKey).(*jwt.Token)

		var claims jwt.MapClaims

		if token != nil {
			if tokenClaims, ok := token.Claims.(jwt.MapClaims); ok {
				claims = tokenClaims
			} else {
				s.respondError(w, r, errors.Unauthorized(op, fmt.Errorf("jwtauth: unknown type of Claims: %T", token.Claims)))
			}
		} else {
			claims = jwt.MapClaims{}
		}

		err, _ := r.Context().Value(jwtauth.ErrorCtxKey).(error)

		if err != nil {
			s.respondError(w, r, errors.Unauthorized(op, err))
			return
		}

		if token == nil || !token.Valid {
			s.respondError(w, r, errors.Unauthorized(op, fmt.Errorf("invalid token")))
			return
		}

		if claims == nil {
			s.respondError(w, r, errors.Unauthorized(op, fmt.Errorf("invalid claims")))
			return
		}

		next.ServeHTTP(w, r)
	})
}
