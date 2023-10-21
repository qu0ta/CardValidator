package main

import (
	"cardValidator/server"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", server.Handler)
	err := http.ListenAndServe("localhost:8080", nil)
	if err != nil {
		log.Println(err.Error())
	}
}
