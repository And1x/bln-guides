package postgres

import (
	"database/sql"
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

// Get returns all User Information
func (m *UserModel) Get(id int) (*models.User, error) {
	return nil, nil
}

// Authenticate checks if user is in DB, returns UserID if successful
func (m *UserModel) Authenticate(name, password string) (int, error) {
	return 0, nil
}
