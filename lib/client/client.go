package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	//"github.com/gorilla/websocket"
	"github.com/hashicorp/errwrap"
)

type loginPayload struct {
	Login_id string `json:"login_id"`
	Password string `json:"password"`
}

type Client struct {
	client      *http.Client
	username    string
	password    string
	baseurl     string
	mmuuid      string
	mmauthtoken string
}

func New(username, password, url string) *Client {
	return &Client{
		client: &http.Client{
			Timeout: time.Second * 5,
		},
		username: username,
		password: password,
		baseurl:  url,
	}
}

func (c *Client) Login() error {
	path := fmt.Sprintf("https://%s/api/v4/users/login", c.baseurl)
	body := &loginPayload{
		Login_id: c.username,
		Password: c.password,
	}

	marshaled, err := json.Marshal(body)
	if err != nil {
		return errwrap.Wrapf("Error constructing POST body: {{err}}", err)
	}
	reqBody := bytes.NewReader(marshaled)

	req, err := http.NewRequest(http.MethodPost, path, reqBody)
	if err != nil {
		return errwrap.Wrapf("Error constructing HTTP request: {{err}}", err)
	}

	res, err := c.client.Do(req)
	if err != nil {
		return errwrap.Wrapf("Error executing HTTP request: {{err}}", err)
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return errwrap.Wrapf("Failed to read request body: {{err}}", err)
		}
		return fmt.Errorf("Status code %d, body: %s", res.StatusCode, string(bodyBytes))
	}

	c.mmuuid = res.Header.Get("token")
	for _, cookie := range res.Cookies() {
		if cookie.Name == "MMAUTHTOKEN" {
			c.mmauthtoken = cookie.Value
			break
		}
	}
	return nil
}
