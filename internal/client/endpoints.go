package client

import (
	"fmt"
	"log"
	"net/http"
	"ttt-game/internal/common"
)

func (c *Client) endpointSSE(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)

	if !ok {
		http.Error(w, "SSE not supported", http.StatusInternalServerError)
		return
	}

	// Send an initial event so the client knows the connection is alive.
	fmt.Fprint(w, ": connected\n\n")
	flusher.Flush()

	// Generate a random ID.
	sseID := common.GenerateAlphaNumId(20)

	// Add the connection to the pool.
	c.mu.Lock()

	c.sseConns[sseID] = &SSE{
		active: true,
		addr:   r.RemoteAddr,
		w:      w,
	}

	c.mu.Unlock()

	// Remove the connection when the client disconnects.
	defer func() {
		c.mu.Lock()
		delete(c.sseConns, sseID)
		c.mu.Unlock()
	}()

	// Keep the handler alive until the client disconnects.
	<-r.Context().Done()
}

func (c *Client) endpointRegister(w http.ResponseWriter, _ *http.Request) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.registered {
		if user, err := c.backend.register(); err != nil {
			log.Println(err.Error())

			c.send(w, http.StatusInternalServerError, map[string]any{
				"message": "Failed to register the user",
			})
		} else {
			c.profile = user
		}
	}

	c.send(w, http.StatusOK, c.profile.WithoutToken())
}

func (c *Client) endpointStartGame(w http.ResponseWriter, r *http.Request) {

}

func (s *Client) endpointClick(w http.ResponseWriter, r *http.Request) {

}
