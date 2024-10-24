package email

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"main/config"
	"main/model"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

type MessageType string

const (
	HtmlMessageType MessageType = "html"
	TextMessageType MessageType = "text"
)

type EmailAddress struct {
	Address string `json:"address"`
}

type Recipient struct {
	EmailAddress EmailAddress `json:"emailAddress"`
}

type BodyContent struct {
	ContentType MessageType `json:"contentType"`
	Content     string      `json:"content"`
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

type httpEmailService struct {
	TenantID     string
	ClientID     string
	ClientSecret string
	UserId       string
	IsEnabled    bool
}

func NewHttpEmailService(config config.ConfigManager) EmailService {
	return &httpEmailService{
		TenantID:     config.GetEmailTenantID(),
		ClientID:     config.GetEmailClientID(),
		ClientSecret: config.GetEmailClientSecret(),
		UserId:       config.GetEmailUserID(),
		IsEnabled:    config.GetIsEmailEnabled(),
	}
}

func (es *httpEmailService) SendActivityHelpEmail(activityHelpEmail *model.ActivityHelpEmail) error {
	panic("unimplemented")
}

func (es *httpEmailService) SendEmail(to, cc []string, subject, content string, contentType ContentType, isSaveToSetItem bool) error {
	saveToSetItem := "false"
	if isSaveToSetItem {
		saveToSetItem = "true"
	}

	messageType := HtmlMessageType
	if contentType == Html {
		messageType = TextMessageType
	}

	sendMailRequest := SendMailRequest{
		Message: EmailMessage{
			Subject: subject,
			Body: BodyContent{
				ContentType: messageType,
				Content:     content,
			},
		},
		SaveToSentItems: saveToSetItem,
	}

	for _, recipient := range to {
		sendMailRequest.Message.ToRecipients = append(sendMailRequest.Message.ToRecipients, Recipient{
			EmailAddress: EmailAddress{
				Address: recipient,
			},
		})
	}

	for _, recipient := range cc {
		sendMailRequest.Message.CcRecipients = append(sendMailRequest.Message.CcRecipients, Recipient{
			EmailAddress: EmailAddress{
				Address: recipient,
			},
		})
	}

	return es.send(sendMailRequest)
}

func (s *httpEmailService) getToken() (*oauth2.Token, error) {
	if !s.IsEnabled {
		return nil, fmt.Errorf("email service is not enabled")
	}
	conf := &clientcredentials.Config{
		ClientID:     s.ClientID,
		ClientSecret: s.ClientSecret,
		TokenURL:     fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/v2.0/token", s.TenantID),
		Scopes:       []string{"https://graph.microsoft.com/.default"},
	}

	token, err := conf.Token(context.Background())
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (es *httpEmailService) send(sendMailRequest SendMailRequest) error {
	requestBody, err := json.Marshal(sendMailRequest)
	if err != nil {
		return fmt.Errorf("error marshalling JSON: %v", err)
	}
	apiEndpoint := fmt.Sprintf("https://graph.microsoft.com/v1.0/users/%s/sendMail", es.UserId)
	req, err := http.NewRequest("POST", apiEndpoint, bytes.NewBuffer(requestBody))
	if err != nil {
		return fmt.Errorf("error creating HTTP request: %v", err)
	}

	token, err := es.getToken()
	if err != nil {
		return fmt.Errorf("error getting token: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+token.AccessToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error performing HTTP request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		responseBody, _ := io.ReadAll(resp.Body)
		fmt.Println("Response Body:", string(responseBody))
		return fmt.Errorf("unexpected response status: %v", resp.Status)
	}

	fmt.Println("Email sent successfully!")
	return nil
}
