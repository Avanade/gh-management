package email

import (
	"fmt"
	"main/pkg/msgraph"
	"os"
)

type EmailMessage struct {
	To      string
	Cc      string
	Subject string
	Body    string
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
