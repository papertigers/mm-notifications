package client

import (
	"fmt"
	"log"

	"github.com/mattermost/platform/model"
)

func (c *Client) StartWatcher() *model.WebSocketClient {
	url := fmt.Sprintf("wss://%s", c.baseurl)
	client, err := model.NewWebSocketClient4(url, c.client.AuthToken)
	if err != nil {
		log.Fatalln(err)
	}
	return client
}
