// A Go package that provides bindings to the Reach SMS REST API
//
// See https://www.reach-interactive.com/developers/api/
package reachsms

import (
	"bytes"
	"encoding/json"
	"errors"
	_ "fmt"
	"io"
	_ "log"
	"net/http"
	"net/url"
)

const (
	apiUrl    = "http://api.reach-interactive.com"
	userAgent = "ReachSms Go Wrapper (" + version + ") - " + repo
	version   = "0.1.0"
	repo      = "github.com/brentmountford/go-reachsms"
)

type ReachSmsApi struct {
	ApiUrl     *url.URL
	UserAgent  string
	Username   string
	Password   string
	httpClient http.Client
}

func Create(username, password string) (*ReachSmsApi, error) {
	url, err := url.Parse(apiUrl)
	if err != nil {
		return nil, err
	}
	reachSmsApi := &ReachSmsApi{
		ApiUrl:    url,
		UserAgent: userAgent,
		Username:  username,
		Password:  password,
	}

	return reachSmsApi, nil
}

type Balance struct {
	Success     bool   `json:"success"`
	Balance     string `json:"balance"`
	Description string `json:"description"`
}

func (c *ReachSmsApi) GetBalance() (balance Balance, err error) {
	req, err := c.newRequest("GET", "/sms/balance", nil)
	if err != nil {
		return
	}
	_, err = c.do(req, &balance)
	if !balance.Success {
		err = errors.New(balance.Description)
	}
	return
}

type MessageDetails struct {
	Method        string `json:"Method"`
	To            string `json:"To"`
	Originator    string `json:"Originator"`
	Text          string `json:"Text"`
	SentDate      string `json:"Sent Date"`
	Status        string `json:"Message Status"`
	DeliveredDate string `json:"Delivered Date"`
	DlrCode       string `json:"DlrCode"`
	Description   string `json:"Description"`
	Reference     string `json:"Reference"`
	Success       bool   `json:"Success"`
}

func (c *ReachSmsApi) GetMessage(id string) (message []MessageDetails, err error) {
	req, err := c.newRequest("GET", "/sms/message/"+id, nil)
	if err != nil {
		return
	}
	_, err = c.do(req, &message)
	return
}

type CommonResponse struct {
	Success     bool   `json:"success"`
	Id          string `json:"id"`
	Description string `json:"description"`
}

// func (c *ReachSmsApi) DeleteMessage(id string) (response CommonResponse, err error) {
// 	response.Success = true
// 	response.Id = "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"
// 	response.Description = ""
// 	return

// }

type Message struct {
	To          string `json:"to"`
	From        string `json:"from"`
	Message     string `json:"message"`
	Valid       string `json:"valid"`
	Reference   string `json:"reference"`
	Callbackurl string `json:"callbackUrl"`
	Scheduled   string `json:"scheduled"`
	Coding      string `json:"coding"`
	Udh         string `json:"udh"`
}

func NewMessage(to, from, message string) *Message {
	return &Message{
		To:      to,
		From:    from,
		Message: message,
		Valid:   "72",
		Coding:  "1",
	}
}

func (c *ReachSmsApi) SendMessage(message *Message) (response []CommonResponse, err error) {
	req, err := c.newRequest("POST", "/sms/message", message)
	if err != nil {
		return
	}
	_, err = c.do(req, &response)
	return
}

func (c *ReachSmsApi) do(req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(v)
	if err != nil {
		return nil, err
	}
	return resp, err
}

func (c *ReachSmsApi) newRequest(method, path string, body interface{}) (*http.Request, error) {
	rel := &url.URL{Path: path}
	u := c.ApiUrl.ResolveReference(rel)
	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.UserAgent)
	req.Header.Add("username", c.Username)
	req.Header.Add("password", c.Password)

	return req, nil
}
