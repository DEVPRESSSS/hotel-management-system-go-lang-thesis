package dto

type PaymentIntentPayMongo struct {
	Data struct {
		Attributes struct {
			Amount               int64    `json:"amount"`
			Currency             string   `json:"currency"`
			PaymentMethodAllowed []string `json:"payment_method_allowed"`
		} `json:"attributes"`
	} `json:"data"`
}
