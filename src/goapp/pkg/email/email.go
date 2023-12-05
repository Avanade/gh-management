package email

import (
	"fmt"
	"main/pkg/msgraph"
	"os"
)

type MessageType string

const (
	HtmlMessageType MessageType = "html"
	TextMessageType MessageType = "text"
)

type EmailMessage struct {
	To      string
	Cc      string
	Subject string
	Body    string
}

type Message struct {
	Subject       string
	Body          Body
	ToRecipients  []Recipient
	CcRecipients  []Recipient
	BccRecipients *[]Recipient
}

type Body struct {
	Content string
	Type    MessageType
}

type Recipient struct {
	Email string
}

func SendEmail(message Message) error {
	sendMailRequest := msgraph.SendMailRequest{
		Message: msgraph.EmailMessage{
			Subject: message.Subject,
			Body: msgraph.BodyContent{
				ContentType: string(message.Body.Type),
				Content:     message.Body.Content,
			},
		},
		SaveToSentItems: "true",
	}

	for _, recipient := range message.ToRecipients {
		sendMailRequest.Message.ToRecipients = append(sendMailRequest.Message.ToRecipients, msgraph.Recipient{
			EmailAddress: msgraph.EmailAddress{
				Address: recipient.Email,
			},
		})
	}

	// DEFAULT CC RECIPIENT
	sendMailRequest.Message.CcRecipients = append(sendMailRequest.Message.CcRecipients, msgraph.Recipient{
		EmailAddress: msgraph.EmailAddress{
			Address: os.Getenv("EMAIL_DEFAULT_CC_RECIPIENT"),
		},
	})

	if len(message.CcRecipients) > 0 {
		for _, recipient := range message.CcRecipients {
			sendMailRequest.Message.CcRecipients = append(sendMailRequest.Message.CcRecipients, msgraph.Recipient{
				EmailAddress: msgraph.EmailAddress{
					Address: recipient.Email,
				},
			})
		}
	}

	userId := os.Getenv("EMAIL_USER_ID")

	err := msgraph.SendEmail(userId, sendMailRequest)
	if err != nil {
		fmt.Println("Error sending email:", err)
		return err
	}

	return nil
}
