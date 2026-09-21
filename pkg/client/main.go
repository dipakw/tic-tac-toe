package main

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"ttt-game/internal/client"
)

//go:embed ui
var uiFs embed.FS

func main() {
	fs, _ := fs.Sub(uiFs, "ui")

	client, err := client.New(&client.Config{
		ServerHost: "127.0.0.1",
		ServerPort: "16000",
		ClientHost: "0.0.0.0",
		ClientPort: "21000",
		UIFS:       fs,
	})

	if err != nil {
		log.Println("Failed to start the client:", err.Error())
		return
	}

	fmt.Println("URL:", client.URL())

	client.Run()
}
