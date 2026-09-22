package common

type PayloadRegister struct {
	Name string `json:"name"`
}

type PayloadStartGame struct {
	PeerID string `json:"peer_id"`
}

type PayloadClick struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type Kv map[string]any
