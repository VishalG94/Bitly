package main

import (
	"os"
	//"github.com/codegangsta/negroni"
	//"github.com/rs/cors"
)

func main() {


	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "3001"
	}
	// start server listen
	// with error handling
	server := NewServer()
	server.Run(":" + port)
}
