package game

import (
	"fmt"
)

func (s *Session) Click(click *Click) (*State, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.done {
		return nil, fmt.Errorf("session is already complete")
	}

	peer, peerIndex := s.getPeer(click.PeerID)

	if peer == nil {
		return nil, fmt.Errorf("peer id \"%s\" doesn't exits in the session", click.PeerID)
	}

	x, y := click.X, click.Y
	cell := s.data.getCell(x, y)

	if cell == nil {
		return nil, fmt.Errorf("cell %d/%d doesn't exist", x, y)
	}

	if cell.Value != "" {
		return nil, fmt.Errorf("cell %d/%d is already used", x, y)
	}

	if err := s.data.setValue(x, y, peer.Sign); err != nil {
		return nil, err
	}

	// Update the count for the peer.
	s.cfg.Peers[peerIndex].count += 1

	// Update the state.
	if err := s.update(click); err != nil {
		return nil, err
	}

	return s.getState()
}

func (s *Session) State() (*State, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.getState()
}

func (s *Session) Peers() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	peers := []string{}

	for _, peer := range s.cfg.Peers {
		peers = append(peers, peer.ID)
	}

	return peers
}

// Helper methods.
func (s *Session) getPeer(id string) (*Peer, int) {
	for i, peer := range s.cfg.Peers {
		if peer.ID == id {
			return peer, i
		}
	}

	return nil, -1
}

// Update calculates and checks if one won the game.
func (s *Session) update(click *Click) error {
	peer, _ := s.getPeer(click.PeerID)

	// If the user hasn't clicked enough,
	// do not check win pattern.
	if peer.count < s.cfg.RootCount {
		return nil
	}

	pattern := s.data.getPattern(click.X, click.Y)

	// Pattern is not complete.
	if len(pattern) != s.cfg.RootCount {
		return nil
	}

	for _, pos := range pattern {
		x, y := pos[0], pos[1]

		if err := s.data.setHighlight(x, y, true); err != nil {
			return err
		}
	}

	return nil
}

func (s *Session) getState() (*State, error) {
	state := &State{
		Done:  s.done,
		Peers: s.cfg.Peers,
		Rows:  s.data.rows,
	}

	return state, nil
}
