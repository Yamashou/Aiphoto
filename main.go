package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	imageupload "github.com/olahol/go-imageupload"
)

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

func handler(rw http.ResponseWriter, req *http.Request) {
	rw = NewHead(rw)
	if req.Method == "OPTIONS" {
		s := req.Header.Get("Access-Control-Request-Headers")
		if strings.Contains(s, "authorization") == true || strings.Contains(s, "Authorization") == true {
			rw.WriteHeader(204)
		}
		rw.WriteHeader(400)
		return
	}

	if req.Method == "POST" {
		img, err := imageupload.Process(req, "file")
		if err != nil {
			fmt.Printf("Process :%s", err)
			return
		}
		thumb, err := imageupload.ThumbnailPNG(img, 300, 300)
		if err != nil {
			fmt.Printf("ThumbanilPNG :%s", err)
			return
		}
		thumb.Save(fmt.Sprintf("%d.png", time.Now().Unix()))
		return
	}
}

//NewHead is request header create
func NewHead(rw http.ResponseWriter) http.ResponseWriter {
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	rw.Header().Set("Access-Control-Allow-Credentials", "true")
	rw.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
	rw.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	rw.Header().Set("Content-Type", "application/json")
	return rw
}
