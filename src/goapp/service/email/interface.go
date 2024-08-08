package email

type ContentType int

const (
	Html ContentType = iota
	Text
)

type EmailSender interface {
	SendEmail(to, cc []string, subject, content string, contentType ContentType, isSaveToSetItem bool) error
}

type EmailService interface {
	Connect() (EmailSender, error)
}
