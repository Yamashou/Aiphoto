package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", ImageSaveHandler)
	if err := http.ListenAndServe("localhost:8000", nil); err != nil {
		log.Fatal(err)
	}
	return
}
