package main

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/lambda"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

// GuardDuty Message

type Message struct {
	Id     string
	Detail MessageDetail
}

type MessageDetail struct {
	Type        string
	Service     MessageService
	Severity    int
	Title       string
	Description string
}

type MessageService struct {
	Action         MessageAction
	EventFirstSeen string
	EventLastSeen  string
	Count          int
}

type MessageAction struct {
	ActionType string
	Map        map[string]json.RawMessage
}

// Slack Message

type SlackMessage struct {
	Text        string       `json:"text"`
	Attachments []Attachment `json:"attachments"`
}

type Attachment struct {
	Text   string  `json:"text"`
	Color  string  `json:"color"`
	Fields []Field `json:"fields"`
}

type Field struct {
	Title string `json:"title"`
	Value string `json:"value"`
	Short bool   `json:"short"`
}

func Handler(ctx context.Context, message Message) (string, error) {
	// timestamp
	location, _ := time.LoadLocation("Asia/Seoul")

	var eventFirstSeen, _ = time.Parse(time.RFC3339, message.Detail.Service.EventFirstSeen)
	var eventLastSeen, _ = time.Parse(time.RFC3339, message.Detail.Service.EventLastSeen)

	var text = "GuardDuty Finding"
	var color = "#0073bb"
	var floatSeverity = float64(message.Detail.Severity)

	// Severity levels
	// - https://docs.aws.amazon.com/guardduty/latest/ug/guardduty_findings.html?icmpid=docs_gd_help_panel
	// - color는 console 배경색 참조
	if floatSeverity >= 7 && floatSeverity <= 8.9 {
		color = "#d13212"
	} else if floatSeverity >= 5 && floatSeverity <= 6.9 {
		color = "#eb5f07"
	}

	var fields []Field
	fields = append(fields, Field{Title: "Type", Value: message.Detail.Type, Short: false})
	fields = append(fields, Field{Title: "Severity", Value: strconv.Itoa(message.Detail.Severity), Short: false})
	fields = append(fields, Field{Title: "Action Type", Value: message.Detail.Service.Action.ActionType, Short: false})
	fields = append(fields, Field{Title: "Title", Value: message.Detail.Title, Short: false})
	fields = append(fields, Field{Title: "Description", Value: message.Detail.Description, Short: false})
	fields = append(fields, Field{Title: "EventFirstSeen", Value: eventFirstSeen.In(location).String(), Short: false})
	fields = append(fields, Field{Title: "EventLastSeen", Value: eventLastSeen.In(location).String(), Short: false})
	fields = append(fields, Field{Title: "Count", Value: strconv.Itoa(message.Detail.Service.Count), Short: false})

	slackMessage := SlackMessage{
		Text: text,
		Attachments: []Attachment{
			Attachment{
				Color:  color,
				Fields: fields,
			},
		},
	}

	slackBody := sendSlack(slackMessage)
	log.Printf("slackBody: %s\n", slackBody)

	return slackBody, nil
}

func sendSlack(slackMessage SlackMessage) string {
	webhookUrl := os.Getenv("WEBHOOK_URL")
	slackBody, _ := json.Marshal(slackMessage)

	req, err := http.NewRequest(http.MethodPost, webhookUrl, bytes.NewBuffer(slackBody))
	if err != nil {
		panic(err)
	}

	req.Header.Add("Content-type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	res, _ := client.Do(req)
	defer res.Body.Close()

	buf := new(bytes.Buffer)
	buf.ReadFrom(res.Body)

	return buf.String()
}

func main() {
	lambda.Start(Handler)
}