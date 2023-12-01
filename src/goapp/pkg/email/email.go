package email

import (
	"bytes"
	"encoding/json"
	"fmt"
	"main/pkg/msgraph"
	"net/http"
	"os"
)

type EmailMessage struct {
	To      string
	Cc      string
	Subject string
	Body    string
}

func SendEmailObsolete(msg EmailMessage) error {
	endpoint := os.Getenv("EMAIL_ENDPOINT")

	postBody, _ := json.Marshal(map[string]string{
		"to":      msg.To,
		"subject": msg.Subject,
		"body":    msg.Body,
	})
	payload := bytes.NewBuffer(postBody)
	_, err := http.Post(endpoint, "application/json", payload)
	if err != nil {
		return err
	}
	return nil
}

func SendEmail(msg EmailMessage) error {
	sendMailRequest := msgraph.SendMailRequest{
		Message: msgraph.EmailMessage{
			Subject: msg.Subject,
			Body: msgraph.BodyContent{
				ContentType: "html",
				Content:     msg.Body,
			},
			ToRecipients: []msgraph.Recipient{
				{
					EmailAddress: msgraph.EmailAddress{
						Address: msg.To,
					},
				},
			},
			CcRecipients: []msgraph.Recipient{
				{
					EmailAddress: msgraph.EmailAddress{
						Address: os.Getenv("EMAIL_DEFAULT_CC_RECIPIENT"),
					},
				},
			},
		},
		SaveToSentItems: "true",
	}

	userId := os.Getenv("EMAIL_USER_ID")

	err := msgraph.SendEmail(userId, sendMailRequest)
	if err != nil {
		fmt.Println("Error sending email:", err)
		return err
	}

	return nil
}
