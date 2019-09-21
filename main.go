package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	RunSimulation(1)
	router := mux.NewRouter().StrictSlash(true)
	s := http.StripPrefix("/static/", http.FileServer(http.Dir("./static/")))
	router.HandleFunc("/", Index)
	router.HandleFunc("/search_move", GameStateAPI)
	router.PathPrefix("/static/").Handler(s)
	log.Fatal(http.ListenAndServe(":8080", router))
}
