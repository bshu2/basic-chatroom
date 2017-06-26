package main

import (
	"os"
	"log"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	server := NewServer()

	router := mux.NewRouter()
	router.HandleFunc("/", homepage_handler)
	router.HandleFunc("/room/{room_code:[a-zA-Z0-9]+}", room_handler)
	router.HandleFunc("/room/{room_code:[a-zA-Z0-9]+}/ws", func(w http.ResponseWriter, r *http.Request) {
		room_code := mux.Vars(r)["room_code"]
		hub, ok := server.hubs[room_code]
		if !(ok) {//create new hub if one does not exist
			log.Printf("new hub for %s", room_code)
			hub = NewHub(server)
			go hub.run()
			server.hubs[room_code] = hub
		}
		runClient(hub, w, r)
	})
    http.Handle("/", router)

	log.Println("starting main.go")
	err := http.ListenAndServe(":" + os.Getenv("PORT"), nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

func homepage_handler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("homepage_handler called")
	http.ServeFile(w, r, "public/home.html")
}

func room_handler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("room_handler called")
	http.ServeFile(w, r, "public/room.html")
}
