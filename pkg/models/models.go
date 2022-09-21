package models

import (
	"database/sql"
	"errors"
	"html/template"
	"time"
)

var ErrNoRows error = sql.ErrNoRows // todo: may change
var ErrInvalidCredentials = errors.New("models: invalid credentials")
var ErrNameAlreadyUsed = errors.New("models: name already used") // todo: maybe better to 1 var like ErrAlreadyExists istead of 3?
var ErrLnaddrAlreadyUsed = errors.New("models: lnaddresse already used")
var ErrEmailAlreadyUsed = errors.New("models: Email already used")

type Guide struct {
	Id      int
	Title   string
	Content template.HTML
	UserID  int
	Created time.Time
	Updated time.Time
}

type User struct {
	Id            int
	Name          string
	Password      []byte
	LNaddr        string
	Email         string
	Created       time.Time
	LNbUID        string
	LNbAdminKey   string
	LNbInvoiceKey string
}
