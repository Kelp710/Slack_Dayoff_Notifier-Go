package main

import (
	"fmt"
	"log"
	"os"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
)

func runSocketMode(client *slack.Client) error {
	//https://github.com/slack-go/slack/blob/master/examples/socketmode/socketmode.go
	socketMode := socketmode.New(
		client,
		// socketmode.OptionDebug(true),
		// socketmode.OptionLog(log.New(os.Stdout, "sm: ", log.Lshortfile|log.LstdFlags)),
	)
	authTest, authTestErr := client.AuthTest()
	if authTestErr != nil {
		fmt.Fprintf(os.Stderr, "SLACK_BOT_TOKEN is invalid: %v\n", authTestErr)
		os.Exit(1)
	}
	selfUserID := authTest.UserID
	fmt.Println("selfUserID", selfUserID)

	go func() {
		for envelope := range socketMode.Events {
			switch envelope.Type {
			case socketmode.EventTypeEventsAPI:
				socketMode.Ack(*envelope.Request)
				eventPayload, _ := envelope.Data.(slackevents.EventsAPIEvent)
				switch eventPayload.Type {
				case slackevents.CallbackEvent:
					switch event := eventPayload.InnerEvent.Data.(type) {
					case *slackevents.AppMentionEvent:
						log.Println("AppMentionEvent", *&event.Text)
						log.Println("AppMentionEvent", *&event.Channel)

						user, err := socketMode.GetUserInfo(*&event.User)
						if err != nil {
							log.Print(err)
						}
						userName := user.Name
						err = reply(socketMode, *&event.Text, *&event.Channel, userName)
						if err != nil {
							log.Print(err)
						}

					}
				}
				// eventHandler.HandleEvent(eventPayload, selfUserID)

			}
		}
	}()

	return socketMode.Run()
}
