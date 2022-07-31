package postgres

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/and1x/bln--h/pkg/models"
)

type GuidesModel struct {
	DB *sql.DB
}

func (g *GuidesModel) GetById(id int) (*models.Guide, error) {

	stmt := `SELECT id, title, content, author, created, updated FROM guides WHERE id= $1`

	row := g.DB.QueryRow(stmt, id)

	mg := &models.Guide{}
	//var mg models.Guide

	err := row.Scan(&mg.Id, &mg.Title, &mg.Content, &mg.Author, &mg.Created, &mg.Updated)
	if err == sql.ErrNoRows {
		return nil, sql.ErrNoRows
	} else if err != nil {
		return nil, err
	}

	return mg, nil
}

func (g *GuidesModel) GetAll() ([]*models.Guide, error) {

	stmt := `SELECT id, title, content, author, created, updated FROM guides`

	rows, err := g.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	guides := []*models.Guide{}

	for rows.Next() {
		mg := &models.Guide{}
		err := rows.Scan(&mg.Id, &mg.Title, &mg.Content, &mg.Author, &mg.Created, &mg.Updated)
		if err != nil {
			return nil, err
		}
		guides = append(guides, mg)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return guides, nil
}

func (g *GuidesModel) Insert(title, content, author string) (int, error) {

	stmt := `INSERT INTO guides (title, content, author, created, updated)
	VALUES($1, $2, $3, $4, $5)
	RETURNING id`

	var id int
	// res, err := g.DB.Exec(stmt, title, content, author, time.Now(), time.Now())
	err := g.DB.QueryRow(stmt, title, content, author, time.Now(), time.Now()).Scan(&id)
	if err != nil {
		return 0, err
	}

	fmt.Println(id)
	return id, nil
}
