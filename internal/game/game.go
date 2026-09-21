package game

import (
	"errors"
	"fmt"
	"sync"
)

func New(cfg *Config) (*Session, error) {
	if cfg == nil {
		return nil, errors.New("config must be provided")
	}

	if len(cfg.Peers) != peersPerSession {
		return nil, fmt.Errorf("expected number of peers is %d but got %d", peersPerSession, len(cfg.Peers))
	}

	for _, peer := range cfg.Peers {
		if peer == nil {
			return nil, errors.New("one or more of the provided peers is nil")
		}
	}

	session := &Session{
		mu:  &sync.RWMutex{},
		cfg: cfg,
	}

	// Set initial cells.
	session.data = session.newCellsData()

	return session, nil
}
