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
	err = fcgi.Serve(l, nil)
	if err != nil {
		log.Println(err)
	}

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
			rw.WriteHeader(http.StatusNoContent)
			return
		}
		rw.WriteHeader(http.StatusNoContent)
		return
	}

	if req.Method == "POST" {
		img, err := imageupload.Process(req, "file")
		if err != nil {
			log.Printf("Process :%s", err)
			rw.WriteHeader(http.StatusNoContent)
			return
		}
		thumb, err := imageupload.ThumbnailPNG(img, 300, 300)
		if err != nil {
			log.Printf("ThumbanilPNG :%s", err)
			rw.WriteHeader(http.StatusNoContent)
			return
		}

		err = thumb.Save(fmt.Sprintf("%d.png", time.Now().Unix()))
		if err != nil {
			log.Printf("save :%s", err)
			rw.WriteHeader(http.StatusNoContent)
			return
		}
		return
	}
	if req.Method != http.MethodPost {
		return
	}
}
