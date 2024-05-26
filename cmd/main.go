package main

import (
	"log"
	"net/http"
	"github.com/NathanPokharel/Connect4/internals"
	"golang.org/x/net/websocket"
)

func main() {
	server := internals.NewServer()

	http.Handle("/ws", websocket.Handler(server.HandleConnection))
    http.HandleFunc("/",func(w http.ResponseWriter, r *http.Request) {
        http.ServeFile(w,r,"frontend/index.html")
    })
	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
