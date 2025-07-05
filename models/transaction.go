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

	Narration string `json:"narration"`
}

// PubSubMessage represents the structure of a Google Cloud Pub/Sub push message
type PubSubMessage struct {
	Message struct {
		Data        string            `json:"data"`
		Attributes  map[string]string `json:"attributes"`
		MessageID   string            `json:"messageId"`
		PublishTime string            `json:"publishTime"`
	} `json:"message"`
	Subscription string `json:"subscription"`
}
