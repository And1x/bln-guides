package postgres

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/and1x/bln--h/pkg/models"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type UserModel struct {
	DB *sql.DB
}

// New inserts a new user into DB
func (u *UserModel) New(name, password, lnaddr, email string) error {

	hashPw, err := bcrypt.GenerateFromPassword([]byte(password), 9)
	if err != nil {
		return err
	}

	stmt := `INSERT INTO users (name, password, lnaddress, email, created)
	VALUES($1, $2, $3, $4, $5)`

	// check for errors like duplicate mailaddresses, lnaddr and name
	_, err = u.DB.Exec(stmt, name, string(hashPw), lnaddr, email, time.Now())
	if err != nil {
		// get pq specific err to distinguish
		// check if err contains name_unique constraint; see in postgres table users -> constraints
		// check : https://github.com/lib/pq/blob/master/error.go
		// pq Code -> "23505": "unique_violation"
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code.Name() == "unique_violation" { // true when a pq.Error is there which contains "unique_violation"
			switch {
			case strings.Contains(pqErr.Message, "name_unique"):
				return models.ErrNameAlreadyUsed
			case strings.Contains(pqErr.Message, "lnaddress_unique"):
				return models.ErrLnaddrAlreadyUsed
			case strings.Contains(pqErr.Message, "email_unique"):
				return models.ErrEmailAlreadyUsed
			}
		}
	}
	return err
}

// Update updates the users Data, email,lnaddr and password are possible to update
func (u *UserModel) UpdateByUid(id int, lnaddr, email string) error {

	stmt := `UPDATE users
		SET lnaddress = $1,
		email = $2
		WHERE id = $3`

	_, err := u.DB.Exec(stmt, lnaddr, email, id)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code.Name() == "unique_violation" {
			switch {
			case strings.Contains(pqErr.Message, "lnaddress_unique"):
				return models.ErrLnaddrAlreadyUsed
			case strings.Contains(pqErr.Message, "email_unique"):
				return models.ErrEmailAlreadyUsed
			}
		}
	}
	return err
}

// Get returns all User Information
func (m *UserModel) GetById(id int) (*models.User, error) {

	if id < 1 { // todo: just check 0 better?
		return nil, errors.New("ivalid UserID")
	}

	stmt := `SELECT id, name, password, lnaddress, email, created FROM users WHERE id = $1`

	row := m.DB.QueryRow(stmt, id)

	mu := &models.User{}

	err := row.Scan(&mu.Id, &mu.Name, &mu.Password, &mu.LNaddr, &mu.Email, &mu.Created)
	if err == sql.ErrNoRows {
		return nil, sql.ErrNoRows
	} else if err != nil {
		return nil, err
	}

	return mu, nil
}

// Authenticate checks if user-name is in DB, compares password-hashes
// returns UserID if successful
func (m *UserModel) Authenticate(name, password string) (int, error) {
	var id int
	var hashPw []byte

	stmt := `SELECT id, password FROM users WHERE name = $1`

	err := m.DB.QueryRow(stmt, name).Scan(&id, &hashPw)

	if err == sql.ErrNoRows {
		return 0, models.ErrInvalidCredentials
	} else if err != nil {
		return 0, err
	}

	err = bcrypt.CompareHashAndPassword(hashPw, []byte(password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return 0, models.ErrInvalidCredentials
	} else if err != nil {
		return 0, nil
	}

	return id, nil

}
