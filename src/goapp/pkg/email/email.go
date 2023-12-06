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
	Subject      string
	Body         Body
	ToRecipients []Recipient
	CcRecipients []Recipient
}

type Body struct {
	Content string
	Type    MessageType
}

type Recipient struct {
	Email string
}

func SendEmail(message Message, hasDefaultCc bool) error {
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

	var ccRecipients []msgraph.Recipient

	// DEFAULT CC RECIPIENT
	if hasDefaultCc {
		if os.Getenv("EMAIL_SUPPORT") != "" {
			ccRecipients = append(ccRecipients, msgraph.Recipient{
				EmailAddress: msgraph.EmailAddress{
					Address: os.Getenv("EMAIL_SUPPORT"),
				},
			})
		}
	}

	if len(message.CcRecipients) > 0 {
		for _, recipient := range message.CcRecipients {
			ccRecipients = append(ccRecipients, msgraph.Recipient{
				EmailAddress: msgraph.EmailAddress{
					Address: recipient.Email,
				},
			})
		}
	}

	if ccRecipients != nil {
		sendMailRequest.Message.CcRecipients = ccRecipients
	} else {
		sendMailRequest.Message.CcRecipients = []msgraph.Recipient{}
	}

	userId := os.Getenv("EMAIL_USER_ID")

	err := msgraph.SendEmail(userId, sendMailRequest)
	if err != nil {
		fmt.Println("Error sending email:", err)
		return err
	}

	return nil
}
