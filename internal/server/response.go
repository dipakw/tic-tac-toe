package server

import (
	"encoding/json"
	"net/http"
)

func (s *Server) send(w http.ResponseWriter, status int, data any) {
	w.WriteHeader(status)
	b, _ := json.Marshal(data)
	w.Write(b)
}
