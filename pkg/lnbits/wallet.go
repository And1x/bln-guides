// wallet.go contains the functions to create a new Wallet and User in LNbits
// after User registers

package lnbits

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

type LNbits struct {
	Conf map[string]string
}

// todo: shrink or make anonymous struct
type createdUserWallet struct {
	UserID   string `json:"id"`
	UserName string `json:"name"`
	AdminID  string `json:"admin"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Wallets  []struct {
		Id         string `json:"id"`
		Admin      string `json:"admin"` // main admin == userID from Wallet who manages "User Manger Extension"
		Name       string `json:"name"`
		User       string `json:"user"`
		AdminKey   string `json:"adminkey"`
		InvoiceKey string `json:"inkey"`
	} `json:"wallets"`
}

// Create a new User and Wallet in LNbits on Registration
func (lnb *LNbits) CreateUserWallet(userName string) (string, string, string, error) {

	requestURL := lnb.Conf["host"] + lnb.Conf["userMgmtEp"]

	newUser := struct {
		AdminID    string `json:"admin_id"`
		UserName   string `json:"user_name"`
		WalletName string `json:"wallet_name"`
	}{
		AdminID:    lnb.Conf["adminUID"],
		UserName:   userName,
		WalletName: userName + "_wallet",
	}

	newUserJson, err := json.Marshal(newUser)
	if err != nil {
		return "", "", "", err
	}
	req, err := http.NewRequest(http.MethodPost, requestURL, bytes.NewBuffer(newUserJson))
	if err != nil {
		//log.Printf("error creating Request: %v", err)
		return "", "", "", err
	}
	apiKey := lnb.Conf["adminAPIkey"]
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("X-Api-Key", apiKey)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		// log.Printf("error making Request: %v", err)
		return "", "", "", err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		// log.Printf("error reading Response: %v", err)
		return "", "", "", err
	}

	//fmt.Println(string(body) + "\n")

	var resp createdUserWallet
	err = json.Unmarshal(body, &resp)
	if err != nil {
		// log.Printf("error unmarshal body: %v", err)
		return "", "", "", err
	}

	// fmt.Println(resp.Wallets[0].User, resp.Wallets[0].AdminKey, resp.Wallets[0].InvoiceKey)

	return resp.Wallets[0].User, resp.Wallets[0].AdminKey, resp.Wallets[0].InvoiceKey, nil
}

func (lnb *LNbits) GetBalance(invoiceKey string) (int, error) {

	requestURL := lnb.Conf["host"] + lnb.Conf["balanceEp"]

	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		return 0, err
	}
	apiKey := invoiceKey
	req.Header.Set("X-Api-Key", apiKey)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return 0, err
	}

	// todo: don't need the whole struct just Balance...
	type balance struct {
		Id      string `json:"id"`
		Name    string `json:"name"`
		Balance int    `json:"balance"`
	}
	var resp balance

	err = json.Unmarshal(body, &resp)
	if err != nil {
		return 0, err
	}

	return resp.Balance, nil
}
