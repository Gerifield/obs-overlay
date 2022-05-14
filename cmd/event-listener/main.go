package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gempir/go-twitch-irc/v3"

	"github.com/gerifield/obs-overlay/token"
)

func main() {
	channelName := flag.String("channel", "gerifield", "Twitch channel name")
	botName := flag.String("botName", "CoderBot42", "Bot name")
	clientID := flag.String("clientID", "", "Twitch App ClientID")
	clientSecret := flag.String("clientSecret", "", "Twitch App clientSecret")

	flag.Parse()

	tl := token.New(*clientID, *clientSecret)
	log.Println("Fetching token")
	token, err := tl.Get()
	if err != nil {
		log.Println(err)
		return
	}

	client := twitch.NewClient(*botName, "oauth:"+token.AccessToken)

	client.OnPrivateMessage(func(m twitch.PrivateMessage) {
		log.Println(m)
		showMessage(fmt.Sprintf(`<span style="color: %s">%s</span>: %s`, m.User.Color, m.User.DisplayName, m.Message))
	})

	client.Join(*channelName)

	log.Println("Connect with client")
	err = client.Connect()
	if err != nil {
		log.Println(err)
	}
}

func showMessage(msg string) error {
	req, err := http.NewRequest(http.MethodPost, "http://127.0.0.1:8080/http", strings.NewReader(msg))
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	return resp.Body.Close()
}
