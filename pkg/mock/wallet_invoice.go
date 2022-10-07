package mock

type LNbits struct{}

func (lnb *LNbits) CreateUserWallet(userName string) (string, string, string, error) {
	return "", "", "", nil
}

func (lnb *LNbits) GetBalance(invoiceKey string) (int, error) {
	return 0, nil
}

func (lnb *LNbits) CreateInvoice(invoiceKey, message string, amount int) (string, string, error) {
	return "", "", nil
}

func (lnb *LNbits) PayInvoice(paymentRequest, adminKey string) (string, error) {
	return "", nil
}
