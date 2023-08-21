package notification

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

var (
	token tokenInfo
)

type tokenInfo struct {
	AccessToken string
	ExpiresIn   time.Time
}

type Response struct {
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	ExtExpiresIn int    `json:"ext_expires_in"`
	AccessToken  string `json:"access_token"`
}

type MessageBody interface {
	Send() error
}

type MessageType string

const (
	RepositoryHasBeenCreatedMessageType MessageType = "InnerSource.RepositoryHasBeenCreated"
	RepositoryPublicApprovalMessageType MessageType = "InnerSource.RepositoryPublicApproval"
	OrganizationInvitationMessageType   MessageType = "InnerSource.OrganizationInvitation"
)

type Contract struct {
	RequestId   string
	MessageType MessageType
	MessageBody interface{}
}

type RepositoryHasBeenCreatedMessageBody struct {
	Recipients       []string
	GitHubAppLink    string
	OrganizationName string
	RepoLink         string
	RepoName         string
}

type OrganizationInvitationMessageBody struct {
	Recipients       []string
	InvitationLink   string
	OrganizationLink string
	OrganizationName string
}

type RepositoryPublicApprovalMessageBody struct {
	Recipients   []string
	ApprovalLink string
	ApprovalType string
	RepoLink     string
	RepoName     string
	UserName     string
}

func (messageBody RepositoryHasBeenCreatedMessageBody) Send() error {
	messageBody.Recipients = setRecipients(messageBody.Recipients)

	contract := Contract{
		RequestId:   uuid.New().String(),
		MessageType: RepositoryHasBeenCreatedMessageType,
		MessageBody: messageBody,
	}

	err := sendNotification(contract)
	if err != nil {
		return err
	}

	return nil
}

func (messageBody OrganizationInvitationMessageBody) Send() error {
	messageBody.Recipients = setRecipients(messageBody.Recipients)

	contract := Contract{
		RequestId:   uuid.New().String(),
		MessageType: OrganizationInvitationMessageType,
		MessageBody: messageBody,
	}

	err := sendNotification(contract)
	if err != nil {
		return err
	}

	return nil
}

func (messageBody RepositoryPublicApprovalMessageBody) Send() error {
	messageBody.Recipients = setRecipients(messageBody.Recipients)

	contract := Contract{
		RequestId:   uuid.New().String(),
		MessageType: RepositoryPublicApprovalMessageType,
		MessageBody: messageBody,
	}

	err := sendNotification(contract)
	if err != nil {
		return err
	}

	return nil
}

func requestNewToken() (*Response, error) {
	urlPath := fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/v2.0/token", os.Getenv("NOTIFICATION_TENANT_ID"))
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	data := url.Values{}
	data.Set("client_id", os.Getenv("NOTIFICATION_CLIENT_ID"))
	data.Set("scope", os.Getenv("NOTIFICATION_SCOPE"))
	data.Set("client_secret", os.Getenv("NOTIFICATION_CLIENT_SECRET"))
	data.Set("grant_type", "client_credentials")
	encodedData := data.Encode()

	req, err := http.NewRequest("POST", urlPath, strings.NewReader(encodedData))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var response Response
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func setToken() error {
	if token.AccessToken != "" {
		if token.ExpiresIn.After(time.Now()) {
			return nil
		}
	}

	newToken, err := requestNewToken()
	if err != nil {
		log.Println(err.Error())
		return err
	}

	const ALLOWANCE_TIME_BEFORE_EXPIRATION = 99

	duration, _ := time.ParseDuration(fmt.Sprint(newToken.ExpiresIn-ALLOWANCE_TIME_BEFORE_EXPIRATION, "s"))

	expiresin := time.Now().Add(duration)

	token.AccessToken = newToken.AccessToken
	token.ExpiresIn = expiresin

	return nil
}

func setRecipients(recipients []string) []string {
	if os.Getenv("NOTIFICATION_RECIPIENT") != "" {
		return strings.Split(os.Getenv("NOTIFICATION_RECIPIENT"), ",")
	}
	return recipients
}

func sendNotification(c Contract) error {
	err := setToken()
	if err != nil {
		return err
	}

	client := &http.Client{
		Timeout: time.Second * 10,
	}

	postBody, _ := json.Marshal(c)

	reqBody := bytes.NewBuffer(postBody)

	req, err := http.NewRequest("POST", os.Getenv("NOTIFICATION_ENDPOINT"), reqBody)
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", "Bearer "+token.AccessToken)
	req.Header.Add("Content-Type", "application/json")

	log.Printf("REQUEST ID : %s | MESSAGE TYPE : %s", c.RequestId, c.MessageType)

	response, err := client.Do(req)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	return nil
}
