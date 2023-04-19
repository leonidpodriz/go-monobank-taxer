package taxer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	url2 "net/url"
	"time"
)

const BaseURL = "https://taxer.ua"
const LoginURI = "/api/user/login/login"
const AccountsLoadURI = "/api/finances/account/load"
const AccountCreateURI = "/api/finances/account/create"
const OperationsLoadURI = "/api/finances/operation/load"
const OperationsCreateURI = "/api/finances/operation/create"

const RecordsOnPage = 100

var loginCookies = []string{"PHPSESSID", "XSRF-TOKEN"}

type Client struct {
	email       string
	pass        string
	httpClient  *http.Client
	userAccount *UserAccount
}

func NewClient(email, pass string) (*Client, error) {
	jar, err := cookiejar.New(nil)

	if err != nil {
		return nil, err
	}

	httpClient := &http.Client{Jar: jar, Timeout: 5 * time.Second}

	return &Client{
		email:      email,
		pass:       pass,
		httpClient: httpClient,
	}, nil
}

func (c *Client) GetUserAccount() (*UserAccount, error) {
	if !c.IsLoggedIn() {
		if err := c.Login(); err != nil {
			return nil, err
		}
	}

	return c.userAccount, nil
}

func (c *Client) GetAllAccounts(userId int) ([]Account, error) {
	var accounts []Account

	for i := 1; ; i++ {
		res, err := c.GetAccounts(userId, i)

		if err != nil {
			return nil, err
		}

		accounts = append(accounts, res.Accounts...)

		if res.Paginator.CurrentPage == res.Paginator.TotalPages {
			break
		}
	}

	return accounts, nil
}

func (c *Client) GetAccounts(userId, pageNumber int) (*AccountsLoadResponse, error) {
	var res AccountsLoadResponse

	err := c.PerformLoggedInRequest(http.MethodPost, AccountsLoadURI, ItemsRequest{
		Filters: Filters{
			FilterArchived: 0,
		},
		PageNumber:    pageNumber,
		UserId:        userId,
		RecordsOnPage: RecordsOnPage,
	}, &res)

	return &res, err
}

func (c *Client) GetAllOperationsForPeriod(userId, accountId int, from, to time.Time) ([]Operation, error) {
	var operations []Operation

	for i := 1; ; i++ {
		res, err := c.GetOperationsForPeriod(userId, i, accountId, from, to)

		if err != nil {
			return nil, err
		}

		operations = append(operations, res.Operations...)

		if res.Paginator.CurrentPage >= res.Paginator.TotalPages {
			break
		}
	}

	return operations, nil
}

func (c *Client) GetOperationsForPeriod(userId, pageNumber, accountId int, from, to time.Time) (*OperationsLoadResponse, error) {
	var res OperationsLoadResponse

	err := c.PerformLoggedInRequest(http.MethodPost, OperationsLoadURI, ItemsRequest{
		Filters: Filters{
			FilterArchived: 0,
			FilterAccount:  accountId,
			FilterDate: FilterDate{
				DateFrom: int(from.Unix()),
				DateTo:   int(to.Unix()),
			},
		},
		Sorting:       Soring{DESC},
		PageNumber:    pageNumber,
		UserId:        userId,
		RecordsOnPage: RecordsOnPage,
	}, &res)

	return &res, err
}

func (c *Client) CreateOperations(userId int, operations []UncreatedOperation) error {
	var opToCreate []OperationToCreate

	for _, op := range operations {
		if op.Type == "" {
			op.Type = FlowIncome
		}

		if op.FinanceType == "" {
			op.FinanceType = Custom
		}

		opToCreate = append(opToCreate, OperationToCreate{
			UserId:    userId,
			Operation: op,
		})
	}

	req := OperationsCreateRequest{opToCreate}

	return c.PerformLoggedInRequest(http.MethodPost, OperationsCreateURI, req, nil)
}

func (c *Client) PerformLoggedInRequest(method, uri string, payload any, v interface{}) error {
	if !c.IsLoggedIn() {
		if err := c.Login(); err != nil {
			return err
		}
	}

	return c.performRequest(method, uri, payload, v)
}

func (c *Client) IsLoggedIn() bool {
	acc := 0

	url, err := url2.Parse(BaseURL)

	if err != nil {
		return false
	}

	for _, cookie := range c.httpClient.Jar.Cookies(url) {
		for _, loginCookie := range loginCookies {
			if cookie.Name == loginCookie {
				acc++
			}
		}
	}

	return acc == len(loginCookies)
}

func (c *Client) Login() (err error) {
	var res LoginResponse

	err = c.performRequest(http.MethodPost, LoginURI, LoginRequest{Email: c.email, Password: c.pass}, &res)
	c.userAccount = &res.Account

	return err
}

func (c *Client) performRequest(method, uri string, payload any, v interface{}) error {
	var err error
	var jsonBody []byte

	if jsonBody, err = json.Marshal(payload); err != nil {
		return err
	}

	var req *http.Request

	if req, err = http.NewRequest(method, BaseURL+uri, bytes.NewBuffer(jsonBody)); err != nil {
		return err
	}

	var res *http.Response

	if res, err = c.httpClient.Do(req); err != nil {
		return err
	}

	if res.StatusCode > 300 {
		body, err := io.ReadAll(res.Body)

		if err != nil {
			return err
		}
		return fmt.Errorf("status code: %d, response: %s", res.StatusCode, string(body))
	}

	defer func() {
		_ = res.Body.Close()
	}()

	if v == nil {
		return nil
	}

	if err := json.NewDecoder(res.Body).Decode(v); err != nil {
		return err
	}

	return nil
}

func (c *Client) CreateAccount(userId int, account Account) error {
	req := AccountCreateRequest{
		Account: account,
		UserId:  userId,
	}

	return c.PerformLoggedInRequest(http.MethodPost, AccountCreateURI, req, nil)
}
