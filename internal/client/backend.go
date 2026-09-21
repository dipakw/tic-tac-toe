package client

import (
	"fmt"
	"net/http"
	"ttt-game/internal/common"
)

func newBackend(baseUrl string) *Backend {
	b := &Backend{
		baseUrl: baseUrl,
	}

	return b
}

func (b *Backend) register() (*common.User, error) {
	url := fmt.Sprintf("%s/register", b.baseUrl)

	return common.HttpRequest[common.User](http.MethodPost, url, &common.PayloadRegister{
		Name: "", // Let it generate one itself.
	})
}
