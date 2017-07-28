package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	imageupload "github.com/olahol/go-imageupload"
)

//ImageSaveHandler is save png
func ImageSaveHandler(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	rw.Header().Set("Access-Control-Allow-Credentials", "true")
	rw.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
	rw.Header().Set("Access-Control-Allow-Methods", "POST")
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
	nowTime := time.Now().Unix()
	err = thumb.Save(fmt.Sprintf("%d.png", nowTime))
	if err != nil {
		log.Printf("save :%s", err)
		rw.WriteHeader(http.StatusNoContent)
		return
	}
}
