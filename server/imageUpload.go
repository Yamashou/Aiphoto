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
	// err = thumb.Save(filename)
	s3FileName, err := uploadS3(filename, thumb.Data)
	if err != nil {
		panic(err)
	}
	db := ConectDB()
	defer db.Close()
	stmtIns, err := db.Prepare(fmt.Sprintf("INSERT INTO photos ( `lat`, `title`, `long`, `region`, `season`, `era`, `image`, `get_type`,`created_at`, `updated_at`) VALUES ( ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"))
	if err != nil {
		panic(err.Error())
	}
	defer stmtIns.Close()
	_, err = stmtIns.Exec("37.485", "test", "139.958", "Aizu", "Summer", "2017", s3FileName, "remote_test", "2017/10/8", "2017/10/8")
	if err != nil {
		panic(err.Error())
	}
}

func uploadS3(filename string, f []byte) (string, error) {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("AWS_REGION")),
	}))
	uploader := s3manager.NewUploader(sess)

	fileBytes := bytes.NewReader(f)
	fileType := http.DetectContentType(f)

	// Upload the file to S3.
	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket:      aws.String("shop-bot-view"),
		Key:         aws.String("aiphoto" + filename),
		Body:        fileBytes,
		ContentType: aws.String(fileType),
		ACL:         aws.String("public-read-write"),
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload file, %v", err)
	}
	fmt.Printf("file uploaded to, %s\n", result.Location)

	return "aiphoto" + filename, nil
}
