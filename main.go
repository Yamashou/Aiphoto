package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/fcgi"
	"strings"
	"time"

	imageupload "github.com/olahol/go-imageupload"
)

func main() {
	l, err := net.Listen("tcp", ":8000")
	if err != nil {
		return
	}
	http.HandleFunc("/", ImageSaveHandler)
	// if err := http.ListenAndServe("localhost:8000", nil); err != nil {
	// 	log.Fatal(err)
	// }
	fcgi.Serve(l, nil)
}

//ImageSaveHandler is save png
func ImageSaveHandler(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	rw.Header().Set("Access-Control-Allow-Credentials", "true")
	rw.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
	rw.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
	if req.Method == "OPTIONS" {
		s := req.Header.Get("Access-Control-Request-Headers")
		log.Println(s)
		if strings.Contains(strings.ToLower(s), "authorization") {
			rw.WriteHeader(204)
			return
		}
		rw.WriteHeader(400)
		return
	}

	if req.Method == "POST" {
		img, err := imageupload.Process(req, "file")
		if err != nil {
			log.Printf("Process :%s", err)
			rw.WriteHeader(501)
			return
		}
		thumb, err := imageupload.ThumbnailPNG(img, 300, 300)
		if err != nil {
			log.Printf("ThumbanilPNG :%s", err)
			rw.WriteHeader(501)
			return
		}

		err = thumb.Save(fmt.Sprintf("%d.png", time.Now().Unix()))
		if err != nil {
			log.Printf("save :%s", err)
			rw.WriteHeader(501)
			return
		}
		return
	}
	if req.Method != http.MethodPost {
		return
	}
}
