package email

import "main/model"

type ContentType int

const (
	Html ContentType = iota
	Text
)

type EmailSender interface {
	SendEmail(to, cc []string, subject, content string, contentType ContentType, isSaveToSetItem bool) error
	CustomEmailSender
}

type EmailService interface {
	Connect() (EmailSender, error)
}

type CustomEmailSender interface {
	SendActivityHelpEmail(activityHelpEmail *model.ActivityHelpEmail) error
}
