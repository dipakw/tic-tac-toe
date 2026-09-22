package server

import (
	"context"
	"net/http"
	"ttt-game/internal/common"
)

type AuthUserID string

const userIDKey AuthUserID = "userID"

func (s *Server) authenticate(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var userID string
		token := r.URL.Query().Get("token")

		s.mu.RLock()

		for _, u := range s.users {
			if u.token == token {
				userID = u.id
				break
			}
		}

		if userID == "" {
			defer s.mu.RUnlock()

			s.send(w, http.StatusUnauthorized, &common.Kv{
				"message": "Unauthorized",
			})

			return
		}

		s.mu.RUnlock()

		ctx := context.WithValue(r.Context(), userIDKey, userID)
		next(w, r.WithContext(ctx))
	}
}
