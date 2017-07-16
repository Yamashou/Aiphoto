package main

import (
	"database/sql"
	"log"
	"net/http"
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
