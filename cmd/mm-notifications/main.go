package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/deckarep/gosx-notifier"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/papertigers/mm-notifications/lib/client"
	"github.com/spf13/viper"
)

type Post struct {
	Message string `json:"message"`
}

func main() {
	home, err := homedir.Dir()
	if err != nil {
		log.Fatalln(err)
	}
	viper.SetConfigName(".mattermost")
	viper.AddConfigPath(home)

	err = viper.ReadInConfig()
	if err != nil {
		fmt.Println("Config file not found...")
	} else {
		user := viper.GetString("username")
		password := viper.GetString("password")
		url := viper.GetString("url")
		highlight := viper.GetString("highlight")
		highlightSlice := strings.Split(highlight, ",")

		client := client.New(user, password, url)
		err = client.Login()
		if err != nil {
			fmt.Println(err)
		}

		wsClient := client.StartWatcher()
		for {
			mentioned := false

			select {
			case event := <-wsClient.EventChannel:
				if event == nil {
					log.Fatalln(wsClient.ListenError)
				}
				if event.Event == "posted" {
					var post Post
					data := event.Data

					channel := data["channel_display_name"].(string)
					username := data["sender_name"].(string)

					if user == username {
						// Don't notifiy if the message is from ourselves
						continue
					}

					// Assert string, panic otherwise
					raw := data["post"].(string)
					json.Unmarshal([]byte(raw), &post)

					if data["channel_type"] == "D" {
						mentioned = true
					} else {

						for _, h := range highlightSlice {
							if strings.Contains(strings.ToLower(post.Message), strings.ToLower(h)) {
								mentioned = true
							}
						}
					}

					if !mentioned {
						continue
					}

					note := gosxnotifier.NewNotification(post.Message)
					note.Title = channel
					note.Subtitle = username
					note.Push()

					// Ring bell
					fmt.Fprint(os.Stdout, "\a")

					log.Printf("user: %s, message: %s, channel: %s",
						data["sender_name"],
						post.Message,
						data["channel_display_name"])
				}
			case res := <-wsClient.ResponseChannel:
				log.Println(res)
			}
		}
	}

}
