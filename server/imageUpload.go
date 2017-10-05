package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
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
	filename := fmt.Sprintf("%d.png", nowTime)
	err = thumb.Save(filename)
	uploadS3(filename)
	if err != nil {
		log.Printf("save :%s", err)
		rw.WriteHeader(http.StatusNoContent)
		return
	}
}

func uploadS3(filename string) error {
	// The session the S3 Uploader will use
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("AWS_REGION")),
	}))
	// Create an uploader with the session and default options
	uploader := s3manager.NewUploader(sess)

	f, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open file %q, %v", filename, err)
	}
	defer f.Close()
	fileInfo, _ := f.Stat()
	size := fileInfo.Size()
	buffer := make([]byte, size) // read file content to buffer

	f.Read(buffer)
	fileBytes := bytes.NewReader(buffer)
	fileType := http.DetectContentType(buffer)

	// Upload the file to S3.
	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket:      aws.String("shop-bot-view"),
		Key:         aws.String("aiphoto" + filename),
		Body:        fileBytes,
		ContentType: aws.String(fileType),
		ACL:         aws.String("public-read-write"),
	})
	if err != nil {
		return fmt.Errorf("failed to upload file, %v", err)
	}
	fmt.Printf("file uploaded to, %s\n", result.Location)
	return nil
}
