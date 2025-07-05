package models

type Transaction struct {
	ID            int    `json:"id"`
	Direction     string `json:"direction"`
	Amount        int64  `json:"amount"`
	AccountName   string `json:"account_name"`
	AccountNumber string `json:"account_number,omitempty"`
	CurrencyCode  string `json:"currency_code"`
	Status        string `json:"status,omitempty"`
	Reference     string `json:"reference"`
	BankName      string `json:"bank_name"`
	BankCode      string `json:"bank_code"`
	Narration     string `json:"narration"`
}
