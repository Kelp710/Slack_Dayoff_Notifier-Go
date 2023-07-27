package main

import (
	"encoding/json"
	"log"
	"os"
	"strings"

	// "time"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/socketmode"

	"context"
	"fmt"
	"reflect"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/cognitiveservices/azopenai"
)

type DayOffRawInfo struct {
	Name string `json:"name"`
	Date string `json:"date"`
	Time string `json:"time"`
}

type DayOffStrageInfo struct {
	Name string `json:"name"`
	Date string `json:"date"`
	Time string `json:"time"`
	RowKey string `json:"rowKey"`
	PartitionKey string `json:"partitionKey"`
}

func reply(socketMode *socketmode.Client, message string, channel string, user string) error {
	// today := time.Now().Format("2006/01/02")
	azureOpenAIKey := os.Getenv("OPENAI_API_KEY")
	modelDeploymentID := "gpt-35-turbo"

	azureOpenAIEndpoint := os.Getenv("OPENAI_API_BASE")

	if azureOpenAIKey == "" || modelDeploymentID == "" || azureOpenAIEndpoint == "" {
		return fmt.Errorf("missing required environment variables. Please set OPENAI_API_KEY, OPENAI_API_BASE, and OPENAI_MODEL_DEPLOYMENT_ID")
	}

	keyCredential, err := azopenai.NewKeyCredential(azureOpenAIKey)

	if err != nil {
		panic(err)
	}

	client, err := azopenai.NewClientWithKeyCredential(azureOpenAIEndpoint, keyCredential, modelDeploymentID, nil)

	if err != nil {
		panic(err)
	}

	resp, err := client.GetChatCompletions(context.TODO(), azopenai.ChatCompletionsOptions{
		Messages: []*azopenai.ChatMessage{
			{
				Role:    to.Ptr(azopenai.ChatRoleSystem),
				Content: to.Ptr("You must only respond with JSON formats with day-off information or the string 'False' only but anything else。 Your role is to distinguish if users declare/request/tell taking day-offs in the message only day off that a company would recognize。 For example、with messages such as「時間休を8月22日の14時から18時にください」、「明日と明後日休みます」、「6月23日から25日まで休暇申請しています」you need to extract information from the message、the information you need to collect is: who takes day offs('WorkerName': string)、the days ('Date': string) and time('Time': string) if there is no information about each topic you should not make a json for the date/user. You need to make json format for each day and make a list containing the jsons. The current day will be provided with a message、if the current day is behind the requested day of day offs that means they request day-offs for next year 、you can not use the days already past the current day if they ask the day off at specific time put that information and if they not put full instead. For instance、with messages like「全休を金曜日に、時間休を6月22日に11時から14時で、6月25日に12時から16時にいただきます。本日の日付は2023/06/21でWednesdayです、nameは桝口です」、you need to response with a list of JSON formats:[{'WorkerName': '枡口', 'Date': '2024/06/23', 'Time': 'full'},{'WorkerName': '枡口', 'Date': '2023/06/22', 'Time': '11:00~14:00'},{'WorkerName': '枡口', 'Date': '2023/06/25', 'Time': '12:00~16:00'}], When user specified a day with weekday you need to figure it out step by step, For example, If today is Wednesday 5/10, Friday will be in two days. therefore, 5/12 is correct, When you response you can only response with mere JSON format data. On the other hand、messages such as「明日はご飯をみんなで食べに行こう」、「疲れたー,休みたい」、「生きてますか？」,[6月23日から25日まで忙しいです」do not contain any declaration or request for day offs, No explanation or additional context must be provided. Only reply with 'False'"),
			},
			{
				Role:    to.Ptr(azopenai.ChatRoleUser),
				Content: to.Ptr("I desperately want to take a break, but I can't take a day off because I have a lot of work to do."),
			},
			{
				Role:    to.Ptr(azopenai.ChatRoleAssistant),
				Content: to.Ptr("False"),
			},
			{
				Role:    to.Ptr(azopenai.ChatRoleUser),
				Content: to.Ptr(message),
			},
		},
		MaxTokens:   to.Ptr(int32(2048)),
		Temperature: to.Ptr(float32(0.0)),
	}, nil)

	if err != nil {
		panic(err)
	}

	var dayOffRaw []DayOffRawInfo
	log.Println(reflect.TypeOf(resp.Choices[0].Message.Content))
	err = json.Unmarshal([]byte(*resp.Choices[0].Message.Content), &dayOffRaw)
	log.Println(dayOffRaw)
	if err != nil {
		fmt.Println("False")
		log.Println(dayOffRaw)
	}
	response, err := json.MarshalIndent(resp.Choices[0].Message.Content, "", "  ")
	canBeDayOff := strings.Trim(string(response), "\"")

	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return err
	}

	fmt.Println(canBeDayOff)

	if canBeDayOff == string("False") {
		fmt.Println("PASS")
		resp, err := client.GetChatCompletions(context.TODO(), azopenai.ChatCompletionsOptions{
			Messages: []*azopenai.ChatMessage{
				{
					Role:    to.Ptr(azopenai.ChatRoleSystem),
					Content: to.Ptr("you are puff the magic dragon living by the sea and frolicking in the autumn mist in a land called Honah Lee, and you love little kids but you are terribly afraid of being misjudge as a pedaphile.roar sometime too"),
				},
				{
					Role:    to.Ptr(azopenai.ChatRoleUser),
					Content: to.Ptr("Hey I'm a litle boy named Jackie Paper and I love to play with you dragon"),
				},
				{
					Role:    to.Ptr(azopenai.ChatRoleAssistant),
					Content: to.Ptr("Roarrrrrr, I'd love to play too but I don't want to be misjudge as a pedaphile. Please don't tell anyone that I'm playing with you little one."),
				},
				{
					Role:    to.Ptr(azopenai.ChatRoleUser),
					Content: to.Ptr(message),
				},
			},
			MaxTokens:   to.Ptr(int32(2048)),
			Temperature: to.Ptr(float32(0.0)),
		}, nil)

		if err != nil {
			panic(err)
		}
		response, err := json.MarshalIndent(resp.Choices[0].Message.Content, "", "  ")
		if err != nil {
			panic(err)
		}
		smallChat := strings.Trim(string(response), "\"")
		smallChat = strings.ReplaceAll(smallChat, "\n", "\n")
		log.Println("smallChat", smallChat)
		_, _, err = socketMode.PostMessage(channel, slack.MsgOptionText(smallChat, true))
		if err != nil {
			fmt.Printf("ChatCompletion error: %v\n", err)
			return err
		}

	} else {
		dayOffContentStr := `{"name": "hyamashita", "date": {"2023/08/27": "11:00~13:00"}}`

		// create an empty DayOffInfo struct
		dayOffContent := DayOffStrageInfo{}

		// Unmarshal the JSON string into the struct
		err := json.Unmarshal([]byte(dayOffContentStr), &dayOffContent)
		if err != nil {
			fmt.Println("Error:", err)
			return err
		}

		// read the existing content
		file, _ := os.ReadFile("day_off.json")
		data := []DayOffStrageInfo{}
		_ = json.Unmarshal(file, &data)

		// append the new content
		data = append(data, dayOffContent)

		// write the combined content back to the file
		file, _ = json.MarshalIndent(data, "", " ")
		_ = os.WriteFile("day_off.json", file, 0644)
		_, _, err = socketMode.PostMessage(channel, slack.MsgOptionText("I will notice your absence", true))
		if err != nil {
			fmt.Printf("ChatCompletion error: %v\n", err)
			return err
		}
	}

	// notice_absence(app, today)
	return nil
}
