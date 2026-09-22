package server

import "net/http"

func (s *Server) getAuthenticedUser(r *http.Request) *User {
	userID, ok := r.Context().Value(userIDKey).(string)

	if !ok {
		return nil
	}

	return s.getUser(userID)
}

func (s *Server) getUser(userId string) *User {
	return s.users[userId]
}
