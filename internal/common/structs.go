package common

type PayloadRegister struct {
	Name string `json:"name"`
}

type PayloadStartGame struct {
	PeerID string `json:"peer_id"`
}

type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Token string `json:"token"`
}
