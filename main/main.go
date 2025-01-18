package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	log.Println(mux)
}
