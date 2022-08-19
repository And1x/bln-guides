package models

import (
	"database/sql"
	"html/template"
	"time"
)

type Guide struct {
	Id      int
	Title   string
	Content template.HTML
	Author  string
	Created time.Time
	Updated time.Time
}

var ErrNoRows error = sql.ErrNoRows
