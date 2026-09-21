package client

import (
	"context"
	"net"
	"net/http"
	"sync"
	"ttt-game/internal/common"

	"github.com/gorilla/websocket"
)

type Config struct {
	ServerHost string
	ServerPort string
	ClientHost string
	ClientPort string
}

type Client struct {
	cfg      *Config
	mu       sync.RWMutex
	ln       net.Listener
	service  *http.Server
	ctx      context.Context
	cancel   context.CancelFunc
	sseConns map[string]*SSE
	backend  *Backend

	// Registration info.
	registered bool
	profile    *common.User
}

type SSE struct {
	active bool
	addr   string
	w      http.ResponseWriter
}

type Backend struct {
	baseUrl string
	ws      *websocket.Conn
}
