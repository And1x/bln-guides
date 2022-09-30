package mock

import (
	"errors"
	"time"

	"github.com/and1x/bln--h/pkg/models"
)

type GuidesModel struct{}

var mockGuide = &models.Guide{
	Id:           21,
	Title:        "Cant stop, wont stop!",
	Content:      "Cant rest, cant rest, wont rest, beliving in the proccess - every days a progress - slow steps...",
	UserID:       1,
	Created:      time.Now(),
	Updated:      time.Now(),
	UpvoteAmount: 155,
	UpvoteUsers:  7,
}

// GetById mocked - default error is no rows can be found -
// todo: are there any other errors than can occur? - see: else if stmt in guides.go
func (g *GuidesModel) GetById(id int, inHtml bool) (*models.Guide, error) {
	if id == 21 {
		return mockGuide, nil
	} else {
		return nil, models.ErrNoRows //sql.ErrNoRows
	}
}

func (g *GuidesModel) GetUidByID(id int) (int, error) {
	return 8, nil
}

// Just a mock of the default behaviour
func (g *GuidesModel) GetAll() ([]*models.Guide, error) {
	return []*models.Guide{mockGuide}, nil
}

func (g *GuidesModel) Insert(title, content string, userId int) (int, error) {
	return 0, nil
}

// Just Id error considered// not specific DB errors - they get tested seperately
func (g *GuidesModel) DeleteById(id int) error {
	if id == 21 {
		return nil
	}
	return errors.New("cannot delete missing id")
}

func (g *GuidesModel) UpdateById(id int, title, content string) error {
	return nil
}

func (g *GuidesModel) AddToUpvotes(id, amount int) error {
	return nil
}

func (g *GuidesModel) AddToUpvoteUserCount(id, payerUid int) error {
	return nil
}
