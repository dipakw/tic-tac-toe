package server

import (
	"context"
	"net/http"
)

type AuthUserID string

const userIDKey AuthUserID = "userID"

func (s *Server) authenticate(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var userID string
		token := r.URL.Query().Get("token")

		s.mu.RLock()
		defer s.mu.RUnlock()

		for _, u := range s.users {
			if u.token == token {
				userID = u.id
				break
			}
		}

		if userID == "" {
			s.send(w, http.StatusUnauthorized, map[string]any{
				"message": "Unauthorized",
			})
			return
		}

		ctx := context.WithValue(r.Context(), userIDKey, userID)
		next(w, r.WithContext(ctx))
	}
}
