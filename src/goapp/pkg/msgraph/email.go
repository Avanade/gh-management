package msgraph

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type EmailAddress struct {
	Address string `json:"address"`
}

type Recipient struct {
	EmailAddress EmailAddress `json:"emailAddress"`
}

type BodyContent struct {
	ContentType string `json:"contentType"`
	Content     string `json:"content"`
}

type EmailMessage struct {
	Subject      string      `json:"subject"`
	Body         BodyContent `json:"body"`
	ToRecipients []Recipient `json:"toRecipients"`
	CcRecipients []Recipient `json:"ccRecipients"`
}

type SendMailRequest struct {
	Message         EmailMessage `json:"message"`
	SaveToSentItems string       `json:"saveToSentItems"`
}

func SendEmail(userId string, request SendMailRequest) error {
	accessToken, err := GetToken()
	if err != nil {
		return err
	}

	requestBody, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("error marshalling JSON: %v", err)
	}

	// Use a different API endpoint based on your requirements
	apiEndpoint := fmt.Sprintf("https://graph.microsoft.com/v1.0/users/%s/sendMail", userId)
	req, err := http.NewRequest("POST", apiEndpoint, bytes.NewBuffer(requestBody))
	if err != nil {
		return fmt.Errorf("error creating HTTP request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error performing HTTP request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("unexpected response status: %v", resp.Status)
	}

	fmt.Println("Email sent successfully!")
	return nil
}
