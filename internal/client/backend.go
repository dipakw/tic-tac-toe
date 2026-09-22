package client

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"ttt-game/internal/common"
	"ttt-game/internal/game"

	"github.com/gorilla/websocket"
)

func (c *Client) newBackend(baseUrl string) *Backend {
	b := &Backend{
		baseUrl: baseUrl,
		c:       c,
	}

	return b
}

func (b *Backend) register() (*common.User, error) {
	url := fmt.Sprintf("%s/register", b.baseUrl)

	user, _, err := common.HttpRequest[common.User](http.MethodPost, url, &common.PayloadRegister{
		Name: "", // Let it generate one itself.
	})

	return user, err
}

func (b *Backend) addLiveConnection(token string) (*websocket.Conn, error) {
	url := fmt.Sprintf("%s/add-live-connection?token=%s", b.baseUrl, token)
	url = strings.Replace(url, "http://", "ws://", 1)
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	return conn, err
}

func (b *Backend) playWith(peerId, token string) (string, error) {
	if peerId == "" {
		return "", errors.New("peer id is required")
	}

	url := fmt.Sprintf("%s/start-game?token=%s", b.baseUrl, token)

	res, _, err := common.HttpRequest[common.Kv](http.MethodPost, url, &common.PayloadStartGame{
		PeerID: peerId,
	})

	if err != nil {
		return "", err
	}

	return (*res)["session_id"].(string), nil
}

func (b *Backend) listen() {
	var msg common.Kv

	for {
		err := b.ws.ReadJSON(&msg)

		if err != nil {
			b.ws.Close()
			log.Println("live connection with the server has broken")
			break
		}

		event, ok := msg["event"].(string)

		if !ok || event == "" {
			continue
		}

		switch event {
		case "paired":
			if peer, err := common.DecodeAnyAs[game.Peer](msg["peer"]); err != nil {
				log.Println("failed to get the peer info:", err.Error())
			} else {
				b.c.peer = &common.User{
					ID:   peer.ID,
					Name: peer.Name,
				}

				b.c.mode = "play"
				b.c.pushState()
			}

		default:
		}
	}
}

func (b *Backend) click(token string) {

}
