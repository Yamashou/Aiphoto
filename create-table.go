package main

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	dbconf := "user:pass@tcp(mysql:3306)/db"
	db, err := sql.Open("mysql", dbconf)
	if err != nil {
		log.Fatalf("ERROR: %v", err)
		return
	}
	defer db.Close()

	_, err = db.Exec(
		"CREATE TABLE `photos` (`id` int(11) NOT NULL AUTO_INCREMENT,`lat` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,`title` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,`long` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,`region` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,`season` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,`era` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,`image` text COLLATE utf8_unicode_ci,`get_type` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,`created_at` datetime NOT NULL,`updated_at` datetime NOT NULL,PRIMARY KEY (`id`)) ENGINE=InnoDB AUTO_INCREMENT=164 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;")
	if err != nil {
		log.Fatalf("CREATE TABLE: %v", err)
	}
	log.Println("SUCSESS!!")
	return
}
