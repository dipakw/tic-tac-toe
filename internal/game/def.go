package game

import "sync"

const (
	peersPerSession = 2
)

type CellsData struct {
	rootCount int
	rows      [][]*Cell
}

type Config struct {
	Peers     []*Peer
	RootCount int
}

type Cell struct {
	Value     string `json:"value"`
	Highlight bool   `json:"highlight"`
}

type Peer struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Desc   string `json:"desc"`
	Avatar string `json:"avatar"`
	Sign   string `json:"sign"`

	// Internal.
	// This holds the number of cells occupied by the peer.
	// It comes handy to short circut pattern checking when there are not enough slots used.
	count int
}

type Session struct {
	mu   *sync.RWMutex
	cfg  *Config
	data *CellsData
	done bool
}

type Click struct {
	PeerID string `json:"peer_id"`
	X      int    `json:"x"`
	Y      int    `json:"y"`
}

type State struct {
	Done  bool      `json:"done"`
	Peers []*Peer   `json:"peers"`
	Rows  [][]*Cell `json:"rows"`
}
