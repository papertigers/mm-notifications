package main

import (
	"context"
	"fmt"
	"log"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/papertigers/mm-notifications/lib/client"
	"github.com/spf13/viper"
)

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

		client := client.New(user, password, url)
		err = client.Login()
		if err != nil {
			fmt.Println(err)
		}

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		events := client.StartWatcher(ctx)
		for event := range events {
			fmt.Printf("%v\n", event)
		}
	}

}
