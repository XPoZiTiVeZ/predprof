package main

import (
	"database/sql"
	"fmt"
	"net/http"
)

var secretKey = []byte("D!oGWmLFa2rikog%MR^@xqDgm6sjKSbrznz733FuTVrT$ms2pBiBwKDj%RxmxRjr")
var jwtKey    = []byte("fvpR!tRJW8&Z6Gk!&M*sxo6&jg8*#Sy#yger#ZhhKXM2w3cQFbWu&YPETsLmDTqC")
var db *sql.DB

func main() {
	fmt.Println("Running at http://127.0.0.1:8080")
	db, _ = openDB()
	defer db.Close()
	
	mux := http.NewServeMux()
	registerRoutes(mux)
	httpServer := http.Server{Addr: ":8080", Handler: mux}
	httpServer.ListenAndServe()
}