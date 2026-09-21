package client

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (s *SSE) Write(data any) error {
	if !s.active || s.w == nil {
		return fmt.Errorf("SSE connection is not active")
	}

	payload, err := json.Marshal(data)

	if err != nil {
		return err
	}

	if _, err := fmt.Fprintf(s.w, "data: %s\n\n", payload); err != nil {
		s.active = false
		return err
	}

	if flusher, ok := s.w.(http.Flusher); ok {
		flusher.Flush()
	}

	return nil
}
