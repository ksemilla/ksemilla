package model

type Invoice struct {
	ID          string  `json:"_id" bson:"_id"`
	DateCreated string  `json:"datecreated"`
	From        string  `json:"from"`
	Address     string  `json:"address"`
	Amount      float64 `json:"amount"`
}

type NewInvoice struct {
	DateCreated string  `json:"created"`
	From        string  `json:"from"`
	Address     string  `json:"address"`
	Amount      float64 `json:"amount"`
}
