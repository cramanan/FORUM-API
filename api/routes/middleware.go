package routes

import (
	"context"
	"encoding/json"
	"net/http"
	"real-time-forum/api/database"
)

type context_key string

const contextSessionKey context_key = "session"

func Protected(next http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := database.GetSession(w, r)
		if err != nil {
			writeJSON(w, http.StatusUnauthorized, "Unauthorized")
			return
		}

		ctx := context.WithValue(r.Context(), contextSessionKey, session)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func contextSession(request *http.Request) *database.Session {
	sess, ok := request.Context().Value(contextSessionKey).(*database.Session)
	if ok {
		return sess
	}
	return nil
}

type handlerFunc func(http.ResponseWriter, *http.Request) error

func HandleFunc(fn handlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := fn(w, r); err != nil {
			writeJSON(w, http.StatusInternalServerError, nil)
		}
	}
}

func writeJSON(writer http.ResponseWriter, status int, v any) error {
	writer.Header().Add("Content-Type", "application/json")
	writer.WriteHeader(status)
	return json.NewEncoder(writer).Encode(v)
}
