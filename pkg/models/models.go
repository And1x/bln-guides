package models

import (
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
