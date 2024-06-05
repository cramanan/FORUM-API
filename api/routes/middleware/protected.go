package middleware

import (
	"net/http"
	"real-time-forum/api/database"
)

type ContextKey string

const ContextSessionKey ContextKey = "session"

func Protected(next http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		/*session*/ _, err := database.GetSession(w, r)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// ctx := context.WithValue(r.Context(), ContextSessionKey, session)
		next.ServeHTTP(w, r /*.WithContext(ctx)*/)
	})
}
