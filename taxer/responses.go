package taxer

type LoginResponse struct {
	Language string      `json:"language"`
	Account  UserAccount `json:"account"`
}

type ItemsLoadResponse struct {
	Paginator Paginator `json:"paginator"`
}

type Paginator struct {
	CurrentPage   int `json:"currentPage"`
	RecordsOnPage int `json:"recordsOnPage"`
	TotalPages    int `json:"totalPages"`
	TotalRecords  int `json:"totalRecords"`
}

type AccountsLoadResponse struct {
	ItemsLoadResponse
	Accounts           []Account `json:"accounts"`
	AccountsCurrencies []string  `json:"accountsCurrencies"`
}

type OperationsLoadResponse struct {
	Operations []Operation `json:"operations"`
	Paginator  struct {
		CurrentPage   int `json:"currentPage"`
		RecordsOnPage int `json:"recordsOnPage"`
		TotalPages    int `json:"totalPages"`
		TotalRecords  int `json:"totalRecords"`
	} `json:"paginator"`
	Currencies []string `json:"currencies"`
}
