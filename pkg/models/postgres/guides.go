package postgres

import (
	"bytes"
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"time"

	"github.com/and1x/bln--h/pkg/models"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
)

type GuidesModel struct {
	DB *sql.DB
}

// GetById returns a Guide by ID - if inHtml==true then convertes content from md to Html
func (g *GuidesModel) GetById(id int, inHtml bool) (*models.Guide, error) {

	stmt := `SELECT id, title, content, author, created, updated FROM guides WHERE id = $1`

	row := g.DB.QueryRow(stmt, id)

	mg := &models.Guide{}

	err := row.Scan(&mg.Id, &mg.Title, &mg.Content, &mg.Author, &mg.Created, &mg.Updated)
	if err == sql.ErrNoRows {
		return nil, sql.ErrNoRows
	} else if err != nil {
		return nil, err
	}
	if inHtml {
		mg.Content = mdToHtml(mg.Content)
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
		mg.Content = mdToHtml(mg.Content)
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
	now := time.Now()
	// res, err := g.DB.Exec(stmt, title, content, author, time.Now(), time.Now())
	err := g.DB.QueryRow(stmt, title, content, author, now, now).Scan(&id)
	if err != nil {
		return 0, err
	}

	fmt.Println(id) // todo: delete inproduction
	return id, nil
}

func (g *GuidesModel) DeleteById(id int) error {

	stmt := `DELETE FROM guides WHERE id = $1`

	_, err := g.DB.Exec(stmt, id)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (g *GuidesModel) UpdateById(title, content string, id int) error {

	stmt := `UPDATE guides 
	SET title = $1,
	content = $2,
	updated = $3
	WHERE id = $4`

	_, err := g.DB.Exec(stmt, title, content, time.Now(), id)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// todo: maybe put in another file?
// mdToHtml converts raw MD (as saved in DB) to HTML
func mdToHtml(convert template.HTML) template.HTML {

	// specify goldmark extension
	md := goldmark.New(
		goldmark.WithExtensions(extension.TaskList),
		goldmark.WithExtensions(extension.Footnote),
	)

	var buf bytes.Buffer
	source := []byte(convert)
	if err := md.Convert(source, &buf); err != nil {
		panic(err)
	}
	return template.HTML(buf.Bytes())
}
