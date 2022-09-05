package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

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
	infoLog       *log.Logger
	errorLog      *log.Logger
	templateCache map[string]*template.Template
	guides        interface { // GuidesModel in guides.go & mockguidesModel(for tests) satisfies interface guides hence it implements all methods
		GetById(id int, inHtml bool) (*models.Guide, error)
		GetAll() ([]*models.Guide, error)
		Insert(title, content string, userId int) (int, error)
		DeleteById(id int) error
		UpdateById(id int, title, content string) error
	}
	users interface {
		New(name, password, lnaddr, email string) error
		Get(id int) (*models.User, error)
		Authenticate(name, password string) (int, error)
	}
}

func main() {
	// leveled logging - todo: if needed add warningLog
	infoLog := log.New(os.Stdout, "INFO:\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR:\t", log.Ldate|log.Ltime|log.Lshortfile)

	// DB Postgresql
	connectPsql := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := openDB(connectPsql)
	if err != nil {
		errorLog.Panic(err)
	}
	infoLog.Println("Connected to PostgreSQL")
	defer db.Close()

	// TemplateCache
	templateCache, err := createTemplateCache("./ui/templates/")
	if err != nil {
		errorLog.Fatal(err)
	}

	//  App
	app := &app{
		infoLog:       infoLog,
		errorLog:      errorLog,
		templateCache: templateCache,
		guides:        &postgres.GuidesModel{DB: db},
		users:         &postgres.UserModel{DB: db},
	}

	// HTTP-Server
	srv := &http.Server{
		Addr:     ":8080",
		Handler:  app.routes(),
		ErrorLog: errorLog,
	}

	infoLog.Println("Starting Server on Port :8080")
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
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
