package main

import (
	"flag"
	"fmt"
	"ttt-game/internal/server"
)

func main() {
	port := flag.String("port", "16000", "server port")

	flag.Parse()

	fmt.Println("The server will start on port:", *port)

	_, err := server.New(&server.Config{
		Host: "0.0.0.0",
		Port: *port,
	})

	if err != nil {
		fmt.Println(err.Error())
		return
	}
}
