package main

import (
	"os"
)

func main() {

	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "3003"
	}
	go create_link_queue_receive()
	go hit_link_queue_receive()
	server := NewServer()
	server.Run(":" + port)

}
