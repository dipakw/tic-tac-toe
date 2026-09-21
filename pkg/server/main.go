package main

import (
	"fmt"
	"ttt-game/internal/server"
)

func main() {
	_, err := server.New(&server.Config{
		Host: "0.0.0.0",
		Port: "16000",
	})

	if err != nil {
		fmt.Println(err.Error())
		return
	}
}
