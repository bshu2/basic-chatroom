package main

import (
	"os"
	"log"
	"net/http"
)

func main() {
	hub := NewHub()
	go hub.run()

	fs := http.FileServer(http.Dir("public"))
	http.Handle("/", fs)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		runClient(hub, w, r)
	})

	log.Println("starting main.go")
	err := http.ListenAndServe(":" + os.Getenv("PORT"), nil)
	//err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err);
	}
}