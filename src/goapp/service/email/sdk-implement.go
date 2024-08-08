package email

import (
	"context"
	"main/config"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
	graphmodels "github.com/microsoftgraph/msgraph-sdk-go/models"
	graphusers "github.com/microsoftgraph/msgraph-sdk-go/users"
)

type sdkEmailService struct {
	TenantID     string
	ClientID     string
	ClientSecret string
	UserId       string
	IsEnabled    bool
}

func NewSdkEmailService(config config.ConfigManager) EmailService {
	return &sdkEmailService{
		TenantID:     config.GetEmailTenantID(),
		ClientID:     config.GetEmailClientID(),
		ClientSecret: config.GetEmailClientSecret(),
		UserId:       config.GetEmailUserID(),
		IsEnabled:    config.GetIsEmailEnabled(),
	}
}

func (s *sdkEmailService) Connect() (EmailSender, error) {
	cred, _ := azidentity.NewClientSecretCredential(
		s.TenantID,
		s.ClientID,
		s.ClientSecret,
		nil,
	)

	graphClient, _ := msgraphsdk.NewGraphServiceClientWithCredentials(
		cred, []string{"https://graph.microsoft.com/.default"})

	return &sdkEmailSender{
		GraphServiceClient: graphClient,
		UserId:             s.UserId,
	}, nil
}

type sdkEmailSender struct {
	*msgraphsdk.GraphServiceClient
	UserId string
}

func (es *sdkEmailSender) SendEmail(to, cc []string, subject, content string, contentType ContentType, isSaveToSentItem bool) error {
	requestBody := graphusers.NewItemSendMailPostRequestBody()
	message := graphmodels.NewMessage()
	message.SetSubject(&subject)
	body := graphmodels.NewItemBody()
	bodyType := graphmodels.TEXT_BODYTYPE
	if contentType == Html {
		bodyType = graphmodels.HTML_BODYTYPE
	}
	body.SetContentType(&bodyType)
	body.SetContent(&content)
	message.SetBody(body)

	var toRecipients []graphmodels.Recipientable
	for _, v := range to {
		recipient := graphmodels.NewRecipient()
		emailAddress := graphmodels.NewEmailAddress()
		address := v
		emailAddress.SetAddress(&address)
		recipient.SetEmailAddress(emailAddress)
		toRecipients = append(toRecipients, recipient)
	}
	message.SetToRecipients(toRecipients)

	var ccRecipients []graphmodels.Recipientable
	for _, v := range cc {
		recipient := graphmodels.NewRecipient()
		emailAddress := graphmodels.NewEmailAddress()
		address := v
		emailAddress.SetAddress(&address)
		recipient.SetEmailAddress(emailAddress)
		ccRecipients = append(ccRecipients, recipient)
	}
	message.SetCcRecipients(ccRecipients)

	requestBody.SetMessage(message)
	requestBody.SetSaveToSentItems(&isSaveToSentItem)

	return es.GraphServiceClient.Users().ByUserId(es.UserId).SendMail().Post(context.Background(), requestBody, nil)
}
