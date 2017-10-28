package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	imageupload "github.com/olahol/go-imageupload"
)

//Photo is image data struct
type Photo struct {
	ID        int            `json:"id"`
	Lat       sql.NullString `json:"lat"`
	Title     sql.NullString `json:"title"`
	Long      sql.NullString `json:"long"`
	Region    sql.NullString `json:"region"`
	Season    sql.NullString `json:"season"`
	Era       sql.NullString `json:"era"`
	Image     sql.NullString `json:"image"`
	GetType   sql.NullString `json:"get_type"`
	CreateAt  string         `json:"create_at"`
	UpdatedAt string         `json:"updated_at"`
}

type Upload struct {
	Title string `json:"title"`
	Lat   string `json:"lat"`
	Long  string `json:"long"`
	Image string `json:"image"`
	Date  string `json:"date"`
}

//Photos ....
type Photos struct {
	Photos []Photo `json:"Photos"`
	Size   int     `json:"size"`
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", ImageSaveHandler).Methods(http.MethodPost)
	r.HandleFunc("/list", getImageHeader).Methods(http.MethodGet)
	r.HandleFunc("/list/{season}", getSeasonList).Methods(http.MethodGet)
	r.NotFoundHandler = http.HandlerFunc(NotFoundHandler)
	if err := http.ListenAndServe(":8000", r); err != nil {
		log.Fatal(err)
	}
	return
}

func ConectDB() *sql.DB {
	dbconf := "user:pass@tcp(mysql:3306)/db"
	db, err := sql.Open("mysql", dbconf)
	if err != nil {
		log.Fatalf("ERROR: %v", err)
	}
	return db
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Gorilla!\nNot Found\n"))
}

func NewUpload(bs []byte, filename string) *imageupload.Image {
	img := &imageupload.Image{
		Filename:    filename,
		ContentType: "image/png",
		Data:        bs,
		Size:        len(bs),
	}
	return img
}
