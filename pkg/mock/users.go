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
	Password:      []byte("nnn7c45d237d2c93594e7d0ed795e4bef5f050314cc725844cde607ee63d2a49ba6"),
	LNaddr:        "Satu@payme.com",
	Email:         "Satu@Naku.com",
	Created:       time.Now(),
	LNbUID:        "123456798abc",
	LNbAdminKey:   "987654321zyx",
	LNbInvoiceKey: "no456sfdjfo22",
	Upvote:        "25",
}

func (u *UserModel) New(name, password, lnaddr, email string) error {
	return nil
}

func (u *UserModel) UpdateLNbByName(lnbuid, lnbadminkey, lnbinvoice, name string) error {
	return nil
}

func (u *UserModel) UpdateByUid(id int, lnaddr, email, upvote string) error {
	return nil
}

func (u *UserModel) UpdatePwByUid(id int, password string) error {
	return nil
}

func (u *UserModel) GetById(id int) (*models.User, error) {
	return mockUser, nil
}

func (u *UserModel) GetInvoiceKey(id int) (string, error) {

	if id == 7 {
		return "1234no456sfdjfo22", nil
	} else {
		return "", errors.New("couldnt get invoiceKey")
	}
}
func (u *UserModel) GetAdminKeyAndUpvoteAmount(id int) (string, int, error) {

	if id == 7 {
		return "987654321zyx", 25, nil
	} else {
		return "", 0, errors.New("couldnt get Adminkey and upvote amount")
	}
}

// todo: hash of pw used -- better use bcrypt to compare pw with hash??
func (u *UserModel) Authenticate(name, password string) (int, error) {
	if name == "Satu Naku" && password == "nnn7c45d237d2c93594e7d0ed795e4bef5f050314cc725844cde607ee63d2a49ba6" {
		return 7, nil
	} else {
		return 0, errors.New("username or password wrong")
	}
}
