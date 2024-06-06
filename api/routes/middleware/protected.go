package middleware

import (
	"context"
	"net/http"
	"real-time-forum/api/database"
)

type context_key string

const ContextSessionKey context_key = "session"

func Protected(next http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := database.GetSession(w, r)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), ContextSessionKey, session)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
