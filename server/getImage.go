package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func getImageHeader(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
	w.Header().Set("Access-Control-Allow-Methods", ", GET")
	db := ConectDB()
	defer db.Close()
	if req.Method == http.MethodGet {
		params := req.URL.Query()
		lim, err := strconv.Atoi(params["lim"][0])
		if err != nil {
			lim = 10
		}
		page, err := strconv.Atoi(params["page"][0])
		if err != nil {
			page = 1
		}
		rows, err := db.Query("SELECT `id`, `lat`, `title`, `long`, `region`, `season`, `era`, `image`, `get_type`,`created_at`, `updated_at` FROM photos LIMIT ?,?;", (page-1)*lim, lim)
		items := make([]Photo, 10)
		if err != nil {
			log.Printf("SELECT LIST ERRER:%v", err)
		}
		i := 0
		for rows.Next() {
			var t Photo
			if err := rows.Scan(&t.ID, &t.Lat, &t.Title, &t.Long, &t.Region, &t.Season, &t.Era, &t.Image, &t.GetType, &t.CreateAt, &t.UpdatedAt); err != nil {
				log.Fatal(err)
			}
			items[i] = t
			i++
		}
		photos := Photos{Photos: items, Size: len(items)}
		if err := json.NewEncoder(w).Encode(photos); err != nil {
			panic(err)
		}
	}

}
