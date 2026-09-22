package common

type PayloadRegister struct {
	Name string `json:"name"`
}

type PayloadStartGame struct {
	PeerID string `json:"peer_id"`
}

type Kv map[string]any
