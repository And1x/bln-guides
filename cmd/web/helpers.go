package main

import (
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"
	"strconv"
)

// serverError logs a stack trace + error msg and sends a http Status Error to the User
func (app *app) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n,%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// clientError sends http Status Error to the User - e.g. when there is a "Bad Request"
func (app *app) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

// authUserId returns userID form the user session
func (app *app) authUserId(r *http.Request) int {
	// todo: it's there to test
	if !app.inProduction {
		return 1
	}
	return app.session.GetInt(r, "userID")
}

// getUserName return userName from the user session
func (app *app) getUserName(r *http.Request) string {
	// todo: it's there to test
	if !app.inProduction {
		return "Satu Naku"
	}
	return app.session.GetString(r, "userName")
}

// isAuthorized checks if the users session ID fits to the guide.UserId he wants to edit/delete
// and POST request(form) - seems more cumbersome than just use it in handlers
func (app *app) isAuthorized(guideId int, w http.ResponseWriter, r *http.Request) bool {

	// todo:  it's there to test
	if !app.inProduction {
		return true
	}

	guide, err := app.guides.GetById(guideId, false)
	if err != nil {
		app.errorLog.Printf("couldn't get guide from DB: %v", err)
		return false
	}

	if app.session.GetInt(r, "userID") == guide.UserID {
		return true
	} else {
		return false
	}
}

// upvoteGuide sends a Payment to the author of a guide
// 1. create Invoice (author)
// 2. pay Invoice (Upvoter)
// 3. update DB
func (app *app) upvoteGuide(r *http.Request, guideId string) error {

	// get InvoiceKey from author of the guide
	gid, err := strconv.Atoi(guideId)
	if err != nil {
		return err
	}

	// query DB to get authors user_Id from guide_id
	uid, err := app.guides.GetUidByID(gid)
	if err != nil {
		return err
	}

	ik, err := app.users.GetInvoiceKey(uid)
	if err != nil {
		return err
	}

	// current user to get his Upvote settings and AdminKey to pay invoices
	payer := app.authUserId(r)

	ak, amount, err := app.users.GetAdminKeyAndUpvoteAmount(payer)
	if err != nil {
		return err
	}

	// PaymentRequest needed to pay Invoice
	phIn, pr, err := app.lnProvider.CreateInvoice(ik, "upvote", amount)
	if err != nil {
		return err
	}

	// pay invoice with user who is currently logged in
	phOut, err := app.lnProvider.PayInvoice(pr, ak)
	if err != nil {
		return err
	}

	// compare payment hashes from created and payed invoice // have to be equal for successful payment
	if phIn != phOut {
		return errors.New("payment hashes are not equal")
	}

	// payment successful when reached to this point - errors below are only in DB not but payment related...
	// todo: better to put them directly into the handlers?
	err = app.guides.AddToUpvotes(gid, amount)
	if err != nil {
		app.session.Put(r, "flashMsg", "upvote successful but cannot be shown at the moment")
		app.infoLog.Println(err)
		return nil // return a nil err hence upvote was successful
	}

	err = app.guides.AddToUpvoteUserCount(gid, payer)
	if err != nil {
		app.session.Put(r, "flashMsg", "upvote successful but cannot be shown at the moment")
		app.infoLog.Println(err)
		return nil // return a nil err hence upvote was successful
	}

	return nil
}
