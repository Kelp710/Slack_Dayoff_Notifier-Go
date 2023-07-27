package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/slack-go/slack"
)

	type DayOffInfo struct {
	Name string `json:"name"`
	Date map[string]string `json:"date"`
}

func notice_absence(app *slack.Client, today string) {
	// read the content from day_off.json and print
	fmt.Println(today)
	file, _ := os.ReadFile("day_off.json")
	data := []DayOffInfo{}
	_ = json.Unmarshal(file, &data)
	fmt.Println(data)
	for _, dayOff := range data {
		for date, time := range dayOff.Date {
			if date == today {
				if time == "full" {
					notification := fmt.Sprintf("%s is day off today", dayOff.Name)
					_, _, _ = app.PostMessage("#random", slack.MsgOptionText(notification, true))
				} else {
					notification := fmt.Sprintf("%s is day off today %s", dayOff.Name, time)
					_, _, _ = app.PostMessage("#random", slack.MsgOptionText(notification, true))
				}
			}
		}
	}

}
