package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	socketio "github.com/googollee/go-socket.io"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	server := socketio.NewServer(nil)

	port := os.Getenv("PORT")

	go server.Serve()
	defer server.Close()

	http.Handle("/", server)
	log.Printf("Serving at localhost:%s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("0.0.0.0:%s", port), nil))
}
