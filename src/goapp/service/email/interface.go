package email

import "main/model"

type ContentType int

const (
	Html ContentType = iota
	Text
)

type EmailService interface {
	SendEmail(to, cc []string, subject, content string, contentType ContentType, isSaveToSetItem bool) error
	CustomEmailSender
}

type CustomEmailSender interface {
	SendActivityHelpEmail(activityHelpEmail *model.ActivityHelpEmail) error
}
