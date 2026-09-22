package server

import (
	"log"
	"net/http"
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

	s.mu.Lock()
	s.users[user.peer.ID] = user
	s.mu.Unlock()

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

	s.users[user.id].conn = &LiveConn{
		ws: conn,
	}
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

	s.send(w, http.StatusOK, &common.Kv{
		"session_id": sessId,
	})
}

func (s *Server) endpointClick(w http.ResponseWriter, r *http.Request) {

}
