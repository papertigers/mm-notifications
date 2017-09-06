package client

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/mattermost/platform/model"
)

type Client struct {
	client   *model.User
	username string
	password string
	baseurl  string
	cookie   string
}

func New(username, password, url string) *Client {
	return &Client{
		username: username,
		password: password,
		baseurl:  url,
	}
}

func (c *Client) Login() error {
	url := fmt.Sprintf("https://%s", c.baseurl)
	mmclient := model.NewAPIv4Client(url)
	user, res := mmclient.Login(c.username, c.password)
	if res.StatusCode != http.StatusOK {
		return errors.New("Failed to login")
	}
	c.client = user
	c.cookie = res.Header.Get("Set-Cookie")
	return nil
}
