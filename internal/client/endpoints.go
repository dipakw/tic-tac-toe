package client

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"ttt-game/internal/common"
)

func (c *Client) endpointUI(w http.ResponseWriter, r *http.Request) {
	http.FileServer(http.FS(c.cfg.UIFS)).ServeHTTP(w, r)
}

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

	end := func() {
		c.send(w, http.StatusOK, map[string]any{})
		go c.pushState()
	}

	if c.registered {
		end()
		return
	}

	// 1. Register the user.
	user, err := c.backend.register()

	if err != nil {
		log.Println(err.Error())

		c.send(w, http.StatusInternalServerError, map[string]any{
			"message": "Failed to register the user",
		})

		return
	}

	// 2. Add live connection between this client and server.
	conn, err := c.backend.addLiveConnection(user.Token)

	if err != nil {
		log.Println(err.Error())

		c.send(w, http.StatusInternalServerError, map[string]any{
			"message": "Failed to add the live connection",
		})

		return
	}

	// 3. Listen to server messages.
	go c.backend.listen()

	c.backend.ws = conn
	c.registered = true
	c.profile = user
	c.mode = "ask" // After registration, ask who they want to play with.

	end()
}

func (c *Client) endpointStartGame(w http.ResponseWriter, r *http.Request) {
	c.mu.Lock()
	defer c.mu.Unlock()

	peerId := r.URL.Query().Get("peer_id")
	peerId = strings.TrimSpace(peerId)

	if peerId == "" || peerId == c.profile.ID {
		log.Println("the peer id", peerId, "is empty or self one")
		c.send(w, http.StatusBadRequest, nil)
		return
	}

	sessionId, err := c.backend.playWith(peerId, c.profile.Token)

	if err != nil {
		c.send(w, http.StatusInternalServerError, nil)
		log.Println("failed to start the game:", err.Error())
		return
	}

	c.sessionId = sessionId

	c.send(w, http.StatusOK, &common.Kv{
		"session_id": sessionId,
	})
}

func (s *Client) endpointClick(w http.ResponseWriter, r *http.Request) {

}
