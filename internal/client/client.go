package client

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"
)

func New(cfg *Config) (*Client, error) {
	ctx, cancel := context.WithCancel(context.Background())

	client := &Client{
		cfg:      cfg,
		ctx:      ctx,
		cancel:   cancel,
		sseConns: map[string]*SSE{},
		mode:     "wait",
	}

	client.backend = client.newBackend(fmt.Sprintf("http://%s", net.JoinHostPort(cfg.ServerHost, cfg.ServerPort)))

	if err := client.setup(); err != nil {
		return nil, err
	}

	return client, nil
}

func (c *Client) Addr() string {
	return c.ln.Addr().String()
}

func (c *Client) URL() string {
	_, port, _ := net.SplitHostPort(c.Addr())
	return fmt.Sprintf("http://localhost:%s", port)
}

func (c *Client) Run() {
	mux := http.NewServeMux()

	// Routes.
	mux.HandleFunc("/sse", c.endpointSSE)
	mux.HandleFunc("/register", c.endpointRegister)
	mux.HandleFunc("/start-game", c.endpointStartGame)
	mux.HandleFunc("/click", c.endpointClick)
	mux.HandleFunc("/", c.endpointUI)

	c.service = &http.Server{
		Handler: mux,
	}

	if err := c.service.Serve(c.ln); err != nil &&
		err != http.ErrServerClosed {
		fmt.Printf("HTTP server error: %v\n", err)
	}
}

func (c *Client) setup() error {
	var err error

	c.ln, err = net.Listen("tcp", net.JoinHostPort(c.cfg.ClientHost, c.cfg.ClientPort))

	if err != nil {
		return fmt.Errorf("failed to listen on a random port for the HTTP server: %w", err)
	}

	return nil
}

func (c *Client) pushState() {
	state := &State{
		Mode:    c.mode,
		Profile: c.profile.WithoutToken(),
	}

	time.Sleep(100 * time.Millisecond)

	for _, sse := range c.sseConns {
		sse.Write(state)
	}
}
