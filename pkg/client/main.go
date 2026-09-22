package main

import (
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"net"
	"ttt-game/internal/client"
)

//go:embed ui
var uiFs embed.FS

func main() {
	serverHost := flag.String("server-host", "0.0.0.0", "server host")
	serverPort := flag.String("server-port", "16000", "server port")

	flag.Parse()

	fmt.Println("Client will connect to the server:", net.JoinHostPort(*serverHost, *serverPort))

	fs, _ := fs.Sub(uiFs, "ui")

	client, err := client.New(&client.Config{
		ServerHost: *serverHost,
		ServerPort: *serverPort,
		ClientHost: "0.0.0.0",
		ClientPort: "0",
		UIFS:       fs,
	})

	if err != nil {
		log.Println("Failed to start the client:", err.Error())
		return
	}

	fmt.Println("Play game at:", client.URL())

	client.Run()
}
