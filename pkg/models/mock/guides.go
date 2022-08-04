package mock

import (
	"database/sql"
	"time"

	"github.com/and1x/bln--h/pkg/models"
)

type GuidesModel struct{}

var mockGuide = &models.Guide{
	Id:      21,
	Title:   "Can't stop, won't stop!",
	Content: "Can't rest, can't rest, won't rest, beliving in the proccess - every day's a progress - slow steps, start to love coding - ah yes",
	Author:  "anon",
	Created: time.Now(),
	Updated: time.Now(),
}

// GetById mocked - default error is no rows can be found -
// todo: are there any other errors than can occur? - see: else if stmt in guides.go
func (g *GuidesModel) GetById(id int, inHtml bool) (*models.Guide, error) {
	if id == 21 {
		return mockGuide, nil
	} else {
		return nil, sql.ErrNoRows
	}
}

// Just a mock of the default behaviour
func (g *GuidesModel) GetAll() ([]*models.Guide, error) {
	return []*models.Guide{mockGuide}, nil
}

func (g *GuidesModel) Insert(title, content, author string) (int, error) {
	return 0, nil
}
