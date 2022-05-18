package email

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"text/template"
)

func SendEmail(msg TypEmailMessage) (*http.Response, error) {
	endpoint := os.Getenv("EMAIL_ENDPOINT")

	postBody, _ := json.Marshal(map[string]string{
		"to":      msg.To,
		"subject": msg.Subject,
		"body":    msg.Body,
	})
	payload := bytes.NewBuffer(postBody)
	resp, err := http.Post(endpoint, "application/json", payload)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func ComposeEmail(data TypEmailData) (string, error) {
	t, err := template.ParseFiles("templates/email.html", "templates/buttons.html")
	if err != nil {
		return "", err
	}

	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

type TypEmailMessage struct {
	To      string
	Subject string
	Body    string
}

type TypEmailData struct {
	Subject     string
	Body        string
	ApproveText string
	RejectText  string
	ApproveUrl  string
	RejectUrl   string
}
