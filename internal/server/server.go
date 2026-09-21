package server

import (
	"net"
	"net/http"
	"sync"
	"ttt-game/internal/game"
)

func New(cfg *Config) (*Server, error) {
	server := &Server{
		cfg:      cfg,
		users:    map[string]*User{},
		sessions: map[string]*game.Session{},
		mu:       &sync.RWMutex{},
	}

	if err := server.setup(); err != nil {
		return nil, err
	}

	return server, nil
}

func (s *Server) setup() error {
	var err error

	if s.ln, err = net.Listen("tcp", net.JoinHostPort(s.cfg.Host, s.cfg.Port)); err != nil {
		return err
	}

	mux := http.NewServeMux()

	// Register routes.
	mux.HandleFunc("/register", s.endpointRegister)
	mux.HandleFunc("/start-game", s.authenticate(s.endpointStartGame))
	mux.HandleFunc("/click", s.authenticate(s.endpointClick))

	s.server = &http.Server{
		Handler: mux,
	}

	if err := s.server.Serve(s.ln); err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}
