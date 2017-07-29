package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func getSeasonList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	params := r.URL.Query()
	season := mux.Vars(r)["season"]
	lim, err := strconv.Atoi(params.Get("lim"))
	if err != nil {
		lim = 10
	}
	page, err := strconv.Atoi(params.Get("page"))
	if err != nil {
		page = 1
	}
	photos := seasonList(season, lim, page)
	if err := json.NewEncoder(w).Encode(photos); err != nil {
		panic(err)
	}
}

func seasonList(season string, lim, page int) Photos {
	db := ConectDB()
	defer db.Close()
	rows, err := db.Query("SELECT `id`, `lat`, `title`, `long`, `region`, `season`, `era`, `image`, `get_type`,`created_at`, `updated_at` FROM photos WHERE `season` = ? LIMIT ?,?;", season, (page-1)*lim, lim)
	items := make([]Photo, lim)
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
	return photos
}
