package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
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
	var m Upload
	b, _ := ioutil.ReadAll(req.Body)
	json.Unmarshal(b, &m)

	bs := base64StringtoByte(m.Image)

	nowTime := time.Now().Unix()
	filename := fmt.Sprintf("%d.png", nowTime)

	img := NewUpload(bs, filename)

	thumb, err := imageupload.ThumbnailPNG(img, 300, 300)
	if err != nil {
		panic(err)
	}

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

	t, err := strconv.ParseFloat(m.Date, 64)
	if err != nil {
		panic(err)
	}
	fmt.Println(t / 1000)
	dateU := time.Unix(int64(t)/1000, 0)
	fmt.Println(dateU)
	_, err = stmtIns.Exec(m.Lat, m.Title, m.Long, "Aizu", getSeason(dateU), string(dateU.Year()), s3FileName, "remote_test", time.Now(), time.Now())
	if err != nil {
		panic(err.Error())
	}

}

func uploadS3(filename string, f []byte) (string, error) {
	sess := session.Must(session.NewSession(&aws.Config{
		Region:           aws.String(os.Getenv("AWS_REGION")),
		S3ForcePathStyle: aws.Bool(true),
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

func base64StringtoByte(s string) []byte {
	r, _ := base64.StdEncoding.DecodeString(s)
	return r
}

func getSeason(t time.Time) string {
	if isSpling(t) {
		return "Spling"
	} else if isSummer(t) {
		return "Summer"
	} else if isFall(t) {
		return "Fall"
	} else if isWinter(t) {
		return "Winter"
	}
	return time.Now().Month().String()
}

func isSpling(t time.Time) bool {
	month := t.Month()
	return month == time.May || month == time.April || month == time.March
}

func isSummer(t time.Time) bool {
	month := t.Month()
	return month == time.August || month == time.July || month == time.June
}

func isFall(t time.Time) bool {
	month := t.Month()
	return month == time.November || month == time.October || month == time.September
}

func isWinter(t time.Time) bool {
	month := t.Month()
	return month == time.February || month == time.January || month == time.December
}
