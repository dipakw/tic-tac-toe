package client

import (
	"errors"
	"fmt"
	"net/http"
	"ttt-game/internal/common"
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

func (b *Backend) syncPeers(sessionId, token string) {

}

func (b *Backend) click(token string) {

}
