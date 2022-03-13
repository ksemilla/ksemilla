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

type InvoiceFilterInput struct {
	ID          *string  `json:"id"`
	DateCreated *string  `json:"DateCreated"`
	From        *string  `json:"From"`
	Address     *string  `json:"Address"`
	Amount      *float64 `json:"Amount"`
}

type InvoiceInput struct {
	ID          string  `json:"id"`
	DateCreated string  `json:"DateCreated"`
	From        string  `json:"From"`
	Address     string  `json:"Address"`
	Amount      float64 `json:"Amount"`
}

type PaginatedInvoicesReturn struct {
	Data  []*Invoice `json:"data"`
	Total int        `json:"total"`
}
