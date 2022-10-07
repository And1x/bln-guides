package lnbits

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
)

// CreateInvoice return a payment-Hash and BOLT11 payment-Request
// invoiceKey is from the author of the guide that gets an upvote
func (lnb *LNbits) CreateInvoice(invoiceKey, message string, amount int) (string, string, error) {

	requestURL := lnb.Conf["host"] + lnb.Conf["paymentEp"]

	newInvoice := struct {
		Out     bool   `json:"out"`
		Amount  int    `json:"amount"`
		Memo    string `json:"memo"`
		Unit    string `json:"unit"`
		Webhook string `json:"string"`
	}{
		Out:    false,
		Amount: amount,
		Memo:   "Upvote",
		Unit:   "sat",
	}
	newInvoiceJson, err := json.Marshal(newInvoice)
	if err != nil {
		return "", "", err
	}
	req, err := http.NewRequest(http.MethodPost, requestURL, bytes.NewBuffer(newInvoiceJson))
	if err != nil {
		return "", "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("X-Api-Key", invoiceKey)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", "", err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", "", err
	}

	type invoice struct {
		PaymentHash    string `json:"payment_hash"`
		PaymentRequest string `json:"payment_request"`
	}

	var resp invoice
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return "", "", err
	}
	return resp.PaymentHash, resp.PaymentRequest, nil
}

// PayInvoice returns after successful payment the payment hash -> should be the same Hash like from paymentRequest
func (lnb *LNbits) PayInvoice(paymentRequest, adminKey string) (string, error) {

	requestURL := lnb.Conf["host"] + lnb.Conf["paymentEp"]

	payInvoice := struct {
		Out    bool   `json:"out"`
		Bolt11 string `json:"bolt11"`
	}{
		Out:    true,
		Bolt11: paymentRequest,
	}

	payInvoiceJson, err := json.Marshal(payInvoice)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(http.MethodPost, requestURL, bytes.NewBuffer(payInvoiceJson))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("X-Api-Key", adminKey)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	// LNbits throws for some errors no JSON but responds just with StatusInternalServerError.
	// with already used Invoices -- but not for malformed invoices for example // todo: explore more cases
	if res.StatusCode == http.StatusInternalServerError {
		return "", errors.New("bad request to LNbits - please check Invoice")
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	type payInvoiceResponse struct {
		PaymentHash string `json:"payment_hash"`
		Detail      string `json:"detail"`
	}

	var resp payInvoiceResponse

	err = json.Unmarshal(body, &resp)
	if err != nil {
		return "", err
	}

	// LNbits API returns detail field incase an error occurs otherwise it's empty
	if resp.Detail != "" {
		return "", errors.New(resp.Detail) // raw err from LNbits
	}

	return resp.PaymentHash, nil
}

// func (lnb *LNbits) PayInvoice(paymentRequest, paymentHash, adminKey string) (bool, error) {

// 	requestURL := lnb.Conf["host"] + lnb.Conf["paymentEp"]

// 	payInvoice := struct {
// 		Out    bool   `json:"out"`
// 		Bolt11 string `json:"bolt11"`
// 	}{
// 		Out:    true,
// 		Bolt11: paymentRequest,
// 	}

// 	payInvoiceJson, err := json.Marshal(payInvoice)
// 	if err != nil {
// 		return false, err
// 	}

// 	req, err := http.NewRequest(http.MethodPost, requestURL, bytes.NewBuffer(payInvoiceJson))
// 	if err != nil {
// 		return false, err
// 	}

// 	req.Header.Set("Content-Type", "application/json")
// 	req.Header.Add("X-Api-Key", adminKey)

// 	res, err := http.DefaultClient.Do(req)
// 	if err != nil {
// 		return false, err
// 	}
// 	defer res.Body.Close()

// 	body, err := ioutil.ReadAll(res.Body)
// 	if err != nil {
// 		return false, err
// 	}

// 	type payInvoiceResponse struct {
// 		PaymentHash string `json:"payment_hash"`
// 		Detail      string `json:"detail"`
// 	}

// 	var resp payInvoiceResponse

// 	err = json.Unmarshal(body, &resp)
// 	if err != nil {
// 		return false, err
// 	}

// 	// LNbits API returns detail field incase an error occurs otherwise it's empty
// 	if resp.Detail != "" {
// 		return false, errors.New(resp.Detail)
// 	}

// 	return paymentHash == resp.PaymentHash, nil
// }
