package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func Handler(res http.ResponseWriter, req *http.Request) {
	c, err := upgrader.Upgrade(res, req, nil)
	if err != nil {
		log.Println(err.Error())
	}
	defer c.Close()

	for {
		mt, msg, err := c.ReadMessage()
		if err != nil {
			log.Println(err.Error())
			break
		}
		log.Printf("Received: %s", msg)
		err = c.WriteMessage(mt, msg)
		if err != nil {
			log.Println("write: ", err)
			break
		}
	}
}

func main() {
	godotenv.Load()
	http.HandleFunc("/", Handler)
	port := os.Getenv("PORT")
	log.Printf("Serving at localhost:%s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("0.0.0.0:%s", port), nil))
}
