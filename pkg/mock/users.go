package mock

import (
	"errors"
	"time"

	"github.com/and1x/bln--h/pkg/models"
)

type UserModel struct{}

var mockUser = &models.User{
	Id:            7,
	Name:          "Satu Naku",
	Password:      []byte("oldpassword"),
	LNaddr:        "Satu@payme.com",
	Email:         "Satu@Naku.com",
	Created:       time.Now(),
	LNbUID:        "123456798abc",
	LNbAdminKey:   "987654321zyx",
	LNbInvoiceKey: "no456sfdjfo22",
	Upvote:        "25",
}

func (u *UserModel) New(name, password, lnaddr, email string) error {

	if name == "Satu Naku" {
		return models.ErrNameAlreadyUsed
	} else if lnaddr == "Satu@payme.com" {
		return models.ErrLnaddrAlreadyUsed
	} else if email == "Satu@Naku.com" {
		return models.ErrEmailAlreadyUsed
	}

	return nil
}

func (u *UserModel) UpdateLNbByName(lnbuid, lnbadminkey, lnbinvoice, name string) error {
	return nil
}

func (u *UserModel) UpdateByUid(id int, lnaddr, email, upvote string) error {

	if lnaddr == "Satu@payme.com" {
		return models.ErrLnaddrAlreadyUsed
	} else if email == "Satu@Naku.com" {
		return models.ErrEmailAlreadyUsed
	}

	return nil
}

func (u *UserModel) UpdatePwByUid(id int, password string) error {
	return nil
}

func (u *UserModel) GetById(id int) (*models.User, error) {
	return mockUser, nil
}

func (u *UserModel) GetInvoiceKey(id int) (string, error) {

	if id == 1 {
		return "1234no456sfdjfo22", nil
	} else {
		return "", errors.New("couldnt get invoiceKey")
	}
}
func (u *UserModel) GetAdminKeyAndUpvoteAmount(id int) (string, int, error) {

	if id == 1 {
		return "987654321zyx", 25, nil
	} else {
		return "", 0, errors.New("couldnt get Adminkey and upvote amount")
	}
}

// todo: hash of pw used -- better use bcrypt to compare pw with hash??
func (u *UserModel) Authenticate(name, password string) (int, error) {
	if name == "Satu Naku" && password == "oldpassword" {
		return 7, nil
	} else {
		return 0, models.ErrInvalidCredentials
	}
}
