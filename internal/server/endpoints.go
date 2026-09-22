package server

import (
	"log"
	"net/http"
	"time"
	"ttt-game/internal/common"
	"ttt-game/internal/game"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (s *Server) endpointRegister(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	defer s.mu.Unlock()

	payload, err := common.GetRequestPayloadAs[common.PayloadRegister](r)

	if err != nil {
		s.send(w, http.StatusBadRequest, &common.Kv{
			"message": err.Error(),
		})

		return
	}

	if payload.Name == "" {
		payload.Name = getRandomName()
	}

	id := common.GenerateAlphaNumId(16)

	user := &User{
		id:    id,
		token: common.GenerateAlphaNumId(25),
		peer: &game.Peer{
			ID:     id,
			Name:   payload.Name,
			Avatar: randomAvatarUrl(),
		},
	}

	s.users[user.peer.ID] = user

	s.send(w, http.StatusOK, &common.User{
		ID:    id,
		Name:  payload.Name,
		Token: user.token,
	})
}

func (s *Server) endpointAddLiveConnection(w http.ResponseWriter, r *http.Request) {
	user := s.getAuthenticedUser(r)
	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println("upgrade error:", err)
		return
	}

	s.mu.Lock()
	s.users[user.id].conn = &LiveConn{
		ws: conn,
	}
	s.mu.Unlock()
}

func (s *Server) endpointStartGame(w http.ResponseWriter, r *http.Request) {
	payload, err := common.GetRequestPayloadAs[common.PayloadStartGame](r)

	if err != nil {
		s.send(w, http.StatusBadRequest, &common.Kv{
			"message": err.Error(),
		})

		return
	}

	var peer *User

	if peer = s.getUser(payload.PeerID); peer == nil {
		s.send(w, http.StatusUnprocessableEntity, &common.Kv{
			"message": "Requested peer is not available.",
		})

		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	me := s.getAuthenticedUser(r)

	session, err := game.New(&game.Config{
		RootCount: 3,
		Peers: []*game.Peer{
			me.peer,
			peer.peer,
		},
	})

	if err != nil {
		log.Println(err.Error())

		s.send(w, http.StatusInternalServerError, &common.Kv{
			"message": "Failed to start the game",
		})

		return
	}

	sessId := common.GenerateAlphaNumId(20)

	// Add session id to players.
	s.users[me.id].sessionId = sessId
	s.users[peer.id].sessionId = sessId

	// Add session to the list.
	s.sessions[sessId] = session

	// Assign signs.
	me.peer.Sign = "x"
	peer.peer.Sign = "o"

	// Send session ID to the client.
	go s.send(w, http.StatusOK, &common.Kv{
		"session_id": sessId,
	})

	time.Sleep(100 * time.Millisecond)

	// Let each know the they have paired.
	me.conn.ws.WriteJSON(&common.Kv{
		"event": "paired",
		"peer":  peer.peer,
	})

	peer.conn.ws.WriteJSON(&common.Kv{
		"event": "paired",
		"peer":  me.peer,
	})
}

func (s *Server) endpointClick(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	defer s.mu.Unlock()

	click, err := common.GetRequestPayloadAs[common.PayloadClick](r)

	if err != nil {
		s.send(w, http.StatusBadRequest, &common.Kv{
			"message": err.Error(),
		})

		return
	}

	me := s.getAuthenticedUser(r)
	session := s.sessions[me.sessionId]

	if session == nil {
		return
	}

	session.Click(&game.Click{
		PeerID: me.peer.ID,
		X:      click.X,
		Y:      click.Y,
	})

	go s.send(w, http.StatusOK, nil)

	gameState, _ := session.State()

	event := &common.Kv{
		"event": "cells_data",
		"data":  gameState.Rows,
	}

	// Give each peer the latest data.
	for _, peerId := range session.Peers() {
		s.users[peerId].conn.ws.WriteJSON(event)
	}
}
