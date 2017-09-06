package client

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/mattermost/platform/model"
)

type Client struct {
	client   *model.Client4
	user     *model.User
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
	c.client = mmclient
	user, res := mmclient.Login(c.username, c.password)
	if res.StatusCode != http.StatusOK {
		return errors.New("Failed to login")
	}
	c.user = user
	c.cookie = res.Header.Get("Set-Cookie")
	return nil
}
