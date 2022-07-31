package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "and1"
	password = "4kn0way"
	dbname   = "blnguide"
)

var DB *sql.DB // maybe just db enough hence it's only package level??

func main() {

	connectPsql := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	var err error
	DB, err = openDB(connectPsql)
	if err != nil {
		//panic(err)
		log.Panic(err)
	}
	defer DB.Close()

	mux := http.NewServeMux()
	mux.HandleFunc("/", HomeSiteHandler)
	mux.HandleFunc("/showguides", ShowGuidesHandler)
	mux.HandleFunc("/createguides", CreateGuideHandler)

	log.Println("Starting Server on Port :8080")
	err = http.ListenAndServe(":8080", mux)
	log.Fatal(err)
}

func openDB(settings string) (*sql.DB, error) {
	db, err := sql.Open("postgres", settings)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
