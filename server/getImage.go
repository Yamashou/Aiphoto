package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func getImageHeader(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
	w.Header().Set("Access-Control-Allow-Methods", ", GET")
	dbconf := "user:pass@tcp(mysql:3306)/db"
	db, err := sql.Open("mysql", dbconf)
	if err != nil {
		log.Fatalf("ERROR: %v", err)
		return
	}
	defer db.Close()
	if req.Method == http.MethodGet {
		rows, err := db.Query("SELECT * FROM photos LIMIT ?,?", 0, 25)
		items := make([]Photo, 25)
		if err != nil {
			log.Printf("SELECT LIST ERRER:%v", err)
		}
		i := 0
		for rows.Next() {
			var item Photo
			if err := rows.Scan(&item.Id); err != nil {
				log.Printf("Scan ERRER:%v", err)
			}
			items[i] = item
		}
		photos := Photos{Photos: items, Size: len(items)}
		if err := json.NewEncoder(w).Encode(photos); err != nil {
			panic(err)
		}
	}

}
