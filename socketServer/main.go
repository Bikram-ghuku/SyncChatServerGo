package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-redis/redis/v8"
	socketio "github.com/googollee/go-socket.io"
	"github.com/joho/godotenv"
)

var (
	server *socketio.Server
)

var redisClient = redis.NewClient(&redis.Options{
	Addr: "localhost:6379",
})

var ctx = context.Background()

func Handler() {
	server.OnConnect("/", func(c socketio.Conn) error {
		c.SetContext("")
		fmt.Println("A new User conencted: ", c.ID())
		return nil
	})
	server.OnEvent("/", "message", func(s socketio.Conn, msg string) {
		if err := redisClient.Publish(ctx, "message", msg).Err(); err != nil {
			fmt.Println("Error publishing to redis: ", err.Error())
		}
	})
}

func main() {
	godotenv.Load()
	server = socketio.NewServer(nil)
	Handler()
	port := os.Getenv("PORT")

	go server.Serve()
	defer server.Close()

	http.Handle("/", server)
	log.Printf("Serving at localhost:%s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("0.0.0.0:%s", port), nil))
}
