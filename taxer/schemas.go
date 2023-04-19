package taxer

const FlowIncome = "FlowIncome"
const Custom = "custom"

type UserAccount struct {
	AccountId    int    `json:"accountId"`
	AccountEmail string `json:"accountEmail"`
	AccountName  string `json:"accountName"`
	Users        []User `json:"users"`
	AppletToken  string `json:"appletToken"`
	HashId       string `json:"hashId"`
	Agreement    bool   `json:"agreement"`
	Language     string `json:"language"`
	KeysSync     bool   `json:"keysSync"`
}

type User struct {
	Id           int          `json:"id"`
	Name         string       `json:"name"`
	TitleName    string       `json:"titleName"`
	IsCompany    bool         `json:"isCompany"`
	IsPayable    bool         `json:"isPayable"`
	Subscription Subscription `json:"subscription"`
	IdKey        string       `json:"idKey"`
}

type Subscription struct {
	IsPaid         bool   `json:"isPaid"`
	IsFree         bool   `json:"isFree"`
	WasPaid        bool   `json:"wasPaid"`
	ExpiresDate    string `json:"expiresDate"`
	AvailableToPay bool   `json:"availableToPay"`
	IsInCart       bool   `json:"isInCart"`
}

type Account struct {
	Id                int         `json:"id"`
	Balance           float64     `json:"balance"`
	Title             string      `json:"title"`
	Currency          string      `json:"currency"`
	Num               string      `json:"num"`
	Bank              string      `json:"bank"`
	Mfo               interface{} `json:"mfo"`
	Comment           interface{} `json:"comment"`
	IsArchived        bool        `json:"isArchived"`
	TfBankPlace       *string     `json:"tfBankPlace"`
	TfBankSwift       *string     `json:"tfBankSwift"`
	TfBankCorr        *string     `json:"tfBankCorr"`
	TfBankCorrPlace   *string     `json:"tfBankCorrPlace"`
	TfBankCorrSwift   *string     `json:"tfBankCorrSwift"`
	TfBankCorrAccount *string     `json:"tfBankCorrAccount"`
}

type Operation struct {
	Id                 int         `json:"id"`
	Type               string      `json:"type"`
	Comment            string      `json:"comment"`
	ExchangeDifference interface{} `json:"exchangeDifference"`
	Contents           []Content   `json:"contents"`
	ContractorName     interface{} `json:"contractorName"`
}

type OperationToCreate struct {
	UserId    int                `json:"userId"`
	Operation UncreatedOperation `json:"operation"`
}

type UncreatedOperation struct {
	Type        string           `json:"type"`
	Comment     string           `json:"comment"`
	Contents    []Content        `json:"contents"`
	Timestamp   int              `json:"timestamp"`
	PayedSum    interface{}      `json:"payedSum"`
	FinanceType string           `json:"financeType"`
	Account     OperationAccount `json:"account"`
	Total       float64          `json:"total"`
}

type OperationAccount struct {
	Currency string `json:"currency"`
	Id       int    `json:"id"`
	Title    string `json:"title"`
}

type Content struct {
	Id              int     `json:"id"`
	AccountTitle    string  `json:"accountTitle"`
	AccountCurrency string  `json:"accountCurrency"`
	Comment         string  `json:"comment"`
	Timestamp       int     `json:"timestamp"`
	SumCurrency     float64 `json:"sumCurrency"`
}

func (ua *UserAccount) GetFirstUser() User {
	return ua.Users[0]
}
