package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	socketio "github.com/googollee/go-socket.io"
	"github.com/joho/godotenv"
)

var (
	server *socketio.Server
)

func Handler() {
	server.OnConnect("/", func(c socketio.Conn) error {
		fmt.Println("A new User conencted: ", c.ID())
		return nil
	})
	server.OnEvent("/", "message", func(s socketio.Conn, msg string) {

	})
}

func main() {
	godotenv.Load()
	server = socketio.NewServer(nil)

	port := os.Getenv("PORT")

	go server.Serve()
	defer server.Close()

	http.Handle("/", server)
	log.Printf("Serving at localhost:%s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("0.0.0.0:%s", port), nil))
}
