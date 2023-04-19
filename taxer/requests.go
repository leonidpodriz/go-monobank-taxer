package taxer

const DESC = "DESC"

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ItemsRequest struct {
	Filters       Filters `json:"filters"`
	Sorting       Soring  `json:"sorting,omitempty"`
	PageNumber    int     `json:"pageNumber"`
	UserId        int     `json:"userId"`
	RecordsOnPage int     `json:"recordsOnPage"`
}

type Filters struct {
	FilterArchived int        `json:"filterArchived,omitempty"`
	FilterAccount  int        `json:"filterAccount,omitempty"`
	FilterDate     FilterDate `json:"filterDate,omitempty"`
}

type FilterDate struct {
	DateFrom int `json:"dateFrom"`
	DateTo   int `json:"dateTo"`
}

type Soring struct {
	Date string `json:"date"`
}

type OperationsCreateRequest struct {
	Operations []OperationToCreate `json:"operations"`
}

type AccountCreateRequest struct {
	UserId  int     `json:"userId"`
	Account Account `json:"account"`
}
