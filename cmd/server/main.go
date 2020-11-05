package main

import (
	"log"
	"net/http"

	"github.com/ellull/salutebot"
)

const addr = ":4000"

func main() {
	http.HandleFunc("/", salutebot.SaluteHandler)
	log.Printf("Listening to %s...\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
