package server

import (
	"net"
	"net/http"
	"sync"
	"ttt-game/internal/game"

	"github.com/gorilla/websocket"
)

type Config struct {
	Host string
	Port string
}

type Server struct {
	cfg      *Config
	ln       net.Listener
	users    map[string]*User
	sessions map[string]*game.Session
	server   *http.Server
	mu       *sync.RWMutex
}

type LiveConn struct {
	ws *websocket.Conn
}

type User struct {
	id        string
	conn      *LiveConn
	token     string
	peer      *game.Peer
	sessionId string
}
