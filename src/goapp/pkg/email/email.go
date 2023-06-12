package email

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
)

type EmailMessage struct {
	To      string
	Subject string
	Body    string
}

func SendEmail(msg EmailMessage) (*http.Response, error) {
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
