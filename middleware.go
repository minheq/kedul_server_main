package main

import (
	"context"
	"net/http"

	"github.com/minheq/kedul_server_main/auth"
)

type contextKey struct {
	name string
}

var (
	userCtxKey = &contextKey{"user"}
)

func addCurrentUserContext(authService auth.Service) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, err := authService.GetCurrentUser(r.Context())

			if err != nil {
				http.Error(w, http.StatusText(401), 401)
				return
			}

			ctx := context.WithValue(r.Context(), userCtxKey, user)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
