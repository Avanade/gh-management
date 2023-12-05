package msgraph

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
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

func GetEmailToken() (string, error) {

	urlPath := fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/v2.0/token", os.Getenv("EMAIL_TENANT_ID"))
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	data := url.Values{}
	data.Set("client_id", os.Getenv("EMAIL_CLIENT_ID"))
	data.Set("scope", "https://graph.microsoft.com/.default")
	data.Set("client_secret", os.Getenv("EMAIL_CLIENT_SECRET"))
	data.Set("grant_type", "client_credentials")
	encodedData := data.Encode()

	req, err := http.NewRequest("POST", urlPath, strings.NewReader(encodedData))
	if err != nil {
		return "", err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	response, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	var tokenResponse TokenResponse
	err = json.NewDecoder(response.Body).Decode(&tokenResponse)
	if err != nil {
		return "", err
	}

	return tokenResponse.AccessToken, nil
}

func SendEmail(userId string, request SendMailRequest) error {
	if os.Getenv("EMAIL_ENABLED") != "true" {
		log.Println("Email not sent because it is not enabled.")
		return nil
	}

	accessToken, err := GetEmailToken()
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
		// Print the response body for additional details
		responseBody, _ := io.ReadAll(resp.Body)
		fmt.Println("Response Body:", string(responseBody))
		return fmt.Errorf("unexpected response status: %v", resp.Status)
	}

	fmt.Println("Email sent successfully!")
	return nil
}
