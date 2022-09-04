package models

import (
	"database/sql"
	"errors"
	"html/template"
	"time"
)

var ErrNoRows error = sql.ErrNoRows
var ErrInvalidCredentials = errors.New("models: invalid credentials")
var ErrAlreadyExists = errors.New("models: Already exists")

type Guide struct {
	Id      int
	Title   string
	Content template.HTML
	UserID  int
	Created time.Time
	Updated time.Time
}

type User struct {
	Id       int
	Name     string
	Password []byte
	LNaddr   string
	Email    string
	Created  time.Time
}
