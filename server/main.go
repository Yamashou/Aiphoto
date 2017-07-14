package main

import (
	"log"
	"net/http"
	"time"
)

//Photo is image data struct
type Photo struct {
	Id        int32     `json:"id"`
	Lat       string    `json:"lat"`
	Title     string    `json:"title"`
	Region    string    `json:"region"`
	Season    string    `json:"season"`
	Era       string    `json:"era"`
	Image     string    `json:"image"`
	GetType   string    `json:"gettype"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

//Photos ....
type Photos struct {
	Photos []Photo `json:"Photos"`
	Size   int     `json:"size"`
}

func main() {
	http.HandleFunc("/", ImageSaveHandler)
	http.HandleFunc("/list", getImageHeader)
	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatal(err)
	}
	return
}
