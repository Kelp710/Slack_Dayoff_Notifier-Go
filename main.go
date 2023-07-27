package main

import (
	// "fmt"
	"log"
	"os"

	"github.com/slack-go/slack"
	// "github.com/slack-go/slack/slackevents"
	// "github.com/slack-go/slack/socketmode"
)

func main() {
	appToken := os.Getenv("SLACK_APP_TOKEN")
	oauthToken := os.Getenv("SLACK_BOT_TOKEN")

	app := slack.New(
		oauthToken,
		slack.OptionDebug(true),
		slack.OptionAppLevelToken(appToken),
	)

	err := runSocketMode(app)
	if err != nil {
		log.Fatal(err)
	}
	// client := socketmode.New(app, socketmode.OptionDebug(true),)
	// fmt.Println("client", client.Events)
	// go func() {
	// 	for socketEvent := range client.Events {
	// 		switch socketEvent.Type {
	// 		case socketmode.EventTypeConnecting:
	// 			fmt.Println("Connecting to Slack with Socket Mode...")
	// 		case socketmode.EventTypeConnectionError:
	// 			fmt.Println("Connection failed. Retrying later...")
	// 		case socketmode.EventTypeConnected:
	// 			fmt.Println("Connected to Slack with Socket Mode.")
	// 			app.SendMessage("C01UZJZQZ9M", slack.MsgOptionText("Hello world", false))
	// 		case socketmode.EventType(slackevents.AppMention):
	// 			fmt.Println("AppMentionEvent")
	// 		}
	// 	}
	// }()

	// err := client.Run()
	// if err != nil {
	// 	log.Print(err)
	// }
}
