package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/and1x/bln--h/pkg/lnbits"
	"github.com/and1x/bln--h/pkg/models"
	"github.com/and1x/bln--h/pkg/models/postgres"
	"github.com/golangcollege/sessions"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// load .env vars to get config settings
// todo: extra config package with os.Lookup use if var doesn't exist or use CLI Arguments
func init() {

	if err := godotenv.Load(); err != nil {
		if err := godotenv.Load("./../../.env"); err != nil { // todo: find better solution - only needed for testing
			log.Print(".env File missing")

		}
	}

	// todo: need to switch to this path totest -> improve..
	// see: https://github.com/joho/godotenv/issues/43
	// if err := godotenv.Load("./../../.env"); err != nil {

	// 	log.Print(".env File missing")
	// }
}

type app struct {
	infoLog       *log.Logger
	errorLog      *log.Logger
	session       *sessions.Session
	templateCache map[string]*template.Template
	guides        interface { // GuidesModel in guides.go & mockguidesModel(for tests) satisfies interface guides hence it implements all methods
		GetById(id int, inHtml bool) (*models.Guide, error)
		GetUidByID(id int) (int, error)
		GetAll() ([]*models.Guide, error)
		Insert(title, content string, userId int) (int, error)
		DeleteById(id int) error
		UpdateById(id int, title, content string) error
		AddToUpvotes(id, amount int) error
		AddToUpvoteUserCount(id, payerUid int) error
	}
	users interface {
		New(name, password, lnaddr, email string) error
		UpdateLNbByName(lnbuid, lnbadminkey, lnbinvoice, name string) error
		UpdateByUid(id int, lnaddr, email, upvote string) error
		UpdatePwByUid(id int, password string) error
		GetById(id int) (*models.User, error)
		GetInvoiceKey(id int) (string, error)                   // todo: needed? more in users.go
		GetAdminKeyAndUpvoteAmount(id int) (string, int, error) // todo: needed? more in users.go
		Authenticate(name, password string) (int, error)
	}
	lnProvider interface {
		CreateUserWallet(userName string) (string, string, string, error)
		GetBalance(invoiceKey string) (int, error)
		CreateInvoice(invoiceKey string, amount int) (string, string, error)
		PayInvoice(paymentRequest, paymentHash, adminKey string) (bool, error)
	}
}

func main() {
	infoLog := log.New(os.Stdout, "INFO:\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR:\t", log.Ldate|log.Ltime|log.Lshortfile)

	// Setup session
	seSecret := os.Getenv("SESSION_SECRET")
	session := sessions.New([]byte(seSecret))
	session.Lifetime = 8 * time.Hour

	// Setup DB Postgresql
	host := os.Getenv("bln_pq_host")
	port := os.Getenv("bln_pq_port")
	user := os.Getenv("bln_pq_user")
	password := os.Getenv("bln_pq_password")
	dbname := os.Getenv("bln_pq_name")

	connectPsql := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
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

	// LNbits API cofig // todo: other way to inject config??
	lnbitsConf := map[string]string{
		"host":        os.Getenv("lnb_host"),
		"userMgmtEp":  os.Getenv("lnb_apiUserManager"), // lnbits endpoint to create new users
		"paymentEp":   os.Getenv("lnb_apiInvoice"),     // lnbits endpoint incoming/outgoing invoices
		"balanceEp":   os.Getenv("lnb_apiBalance"),
		"adminUID":    os.Getenv("lnb_umAdminUID"),
		"adminAPIkey": os.Getenv("lnb_umAPIkey"),
	}

	// App
	app := &app{
		infoLog:       infoLog,
		errorLog:      errorLog,
		session:       session,
		templateCache: templateCache,
		guides:        &postgres.GuidesModel{DB: db},
		users:         &postgres.UserModel{DB: db},
		lnProvider:    &lnbits.LNbits{Conf: lnbitsConf},
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
