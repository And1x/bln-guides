package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/and1x/bln--h/pkg/models"
	"github.com/and1x/bln--h/pkg/models/postgres"
	_ "github.com/lib/pq"
)

// instead of const use env vars or cli args
const (
	host     = "localhost"
	port     = 5432
	user     = "and1"
	password = "4kn0way"
	dbname   = "blnguide"
)

type app struct {
	guides interface { // GuidesModel in guides.go satisfies interface guides hence it implements all methods
		GetById(id int, inHtml bool) (*models.Guide, error)
		GetAll() ([]*models.Guide, error)
		Insert(title, content, author string) (int, error)
		DeleteById(id int) error
		UpdateById(title, content string, id int) error
	}
}

func main() {

	connectPsql := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := openDB(connectPsql)
	if err != nil {
		log.Panic(err)
	}
	log.Println("Connected to PostgreSQL")
	defer db.Close()

	a := app{
		guides: &postgres.GuidesModel{DB: db},
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", a.HomeSiteHandler)
	mux.HandleFunc("/showguides", a.ShowGuidesHandler)
	mux.HandleFunc("/createguide", a.CreateGuideHandler)
	mux.HandleFunc("/editguide", a.EditGuidesHandler)

	fs := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fs))

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
