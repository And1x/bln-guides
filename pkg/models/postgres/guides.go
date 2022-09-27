package postgres

import (
	"bytes"
	"database/sql"
	"html/template"
	"log"
	"time"

	"github.com/and1x/bln--h/pkg/models"
	"github.com/lib/pq"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
)

type GuidesModel struct {
	DB *sql.DB
}

// GetById returns a Guide by ID - if inHtml==true then convertes content from md to Html
func (g *GuidesModel) GetById(id int, inHtml bool) (*models.Guide, error) {

	stmt := `SELECT id, title, content, user_id, created, updated, up_total, ups_by_uid FROM guides WHERE id = $1`

	row := g.DB.QueryRow(stmt, id)

	mg := &models.Guide{}
	upvotesByUid := []sql.NullInt16{} // used for conversion of psql array field

	err := row.Scan(&mg.Id, &mg.Title, &mg.Content, &mg.UserID, &mg.Created, &mg.Updated, &mg.UpvoteAmount, pq.Array(&upvotesByUid))
	//err := row.Scan(&mg.Id, &mg.Title, &mg.Content, &mg.UserID, &mg.Created, &mg.Updated, &mg.UpvoteAmount, pq.Array(&mg.UpvoteUsers))
	if err == sql.ErrNoRows {
		return nil, sql.ErrNoRows
	} else if err != nil {
		return nil, err
	}
	if inHtml {
		mg.Content = mdToHtml(mg.Content)
	}

	// upvote is an array of uid in DB - just use length to show Upvote amount
	mg.UpvoteUsers = len(upvotesByUid)

	return mg, nil
}

// GetUidByID gets the User_ID of the guide
func (g *GuidesModel) GetUidByID(id int) (int, error) {

	stmt := `SELECT user_id FROM guides WHERE id = $1`

	row := g.DB.QueryRow(stmt, id)

	var uid int

	err := row.Scan(&uid)
	if err == sql.ErrNoRows {
		return 0, sql.ErrNoRows
	} else if err != nil {
		return 0, err
	}

	return uid, nil
}

func (g *GuidesModel) GetAll() ([]*models.Guide, error) {

	stmt := `SELECT id, title, content, user_id, created, updated, up_total, ups_by_uid 
			FROM guides
			ORDER BY id ASC`

	rows, err := g.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	guides := []*models.Guide{}

	for rows.Next() {
		mg := &models.Guide{}
		upvotesByUid := []sql.NullInt16{}

		err := rows.Scan(&mg.Id, &mg.Title, &mg.Content, &mg.UserID, &mg.Created, &mg.Updated, &mg.UpvoteAmount, pq.Array(&upvotesByUid))
		if err != nil {
			return nil, err
		}
		// convert before appending to []guides
		mg.Content = mdToHtml(mg.Content)
		mg.UpvoteUsers = len(upvotesByUid)

		guides = append(guides, mg)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	//fmt.Println("length...", guides[3].UpvoteUsers, len(guides[3].UpvoteUsers))

	return guides, nil
}

func (g *GuidesModel) Insert(title, content string, userID int) (int, error) {

	stmt := `INSERT INTO guides (title, content, user_id, created, updated)
	VALUES($1, $2, $3, $4, $5)
	RETURNING id`

	var id int
	now := time.Now()
	// res, err := g.DB.Exec(stmt, title, content, author, time.Now(), time.Now())
	err := g.DB.QueryRow(stmt, title, content, userID, now, now).Scan(&id)
	if err != nil {
		return 0, err
	}

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

func (g *GuidesModel) UpdateById(id int, title, content string) error {

	stmt := `UPDATE guides 
			SET title = $1,
			content = $2,
			updated = $3
			WHERE id = $4`

	_, err := g.DB.Exec(stmt, title, content, time.Now(), id)
	if err != nil {
		return err
	}
	return nil
}

func (g *GuidesModel) AddToUpvotes(id, amount int) error {

	stmt := `UPDATE guides
			SET up_total = up_total + $1
			WHERE id = $2`

	_, err := g.DB.Exec(stmt, amount, id)
	if err != nil {
		return err
	}
	return nil
}

func (g *GuidesModel) AddToUpvoteUserCount(id, payerUid int) error {

	stmt := `UPDATE guides
			SET ups_by_uid = CASE WHEN CAST($1 AS INTEGER) = ANY(ups_by_uid) 
			THEN ups_by_uid ELSE ARRAY_APPEND(ups_by_uid, CAST($1 AS INTEGER)) END
			WHERE id = $2`

	_, err := g.DB.Exec(stmt, payerUid, id)
	if err != nil {
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
