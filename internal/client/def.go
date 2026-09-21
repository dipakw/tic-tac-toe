package client

import (
	"context"
	"io/fs"
	"net"
	"net/http"
	"sync"
	"ttt-game/internal/common"
	"ttt-game/internal/game"

	"github.com/gorilla/websocket"
)

type Config struct {
	ServerHost string
	ServerPort string
	ClientHost string
	ClientPort string
	UIFS       fs.FS
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
	mode       string
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

type State struct {
	Mode    string         `json:"mode"`
	Rows    [][]*game.Cell `json:"rows"`
	Profile *common.User   `json:"profile"`
}
