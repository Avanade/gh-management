package email

import (
	"context"
	"fmt"
	"main/config"
	"main/model"
	"os"
	"strings"

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

func (es *sdkEmailSender) SendActivityHelpEmail(activityHelpEmail *model.ActivityHelpEmail) error {
	if os.Getenv("EMAIL_SUPPORT") == "" {
		return fmt.Errorf("EMAIL_SUPPORT is not set")
	}

	bodyTemplate := `<tr>
						<th class="center-table">
							<table style="width: 100%; max-width: 700px;" class="margin-auto">
								<tr>
									<td style="font-size: 18px; font-weight: 600; padding-top: 20px">
										Activity - Request for Help
									</td>
								</tr>
							</table>
						</th>
					</tr>
					<tr>
						<td class="center-table" align="center">
							<table style="width: 100%; max-width: 700px;" class="margin-auto">
								<tr>
									<td style="padding-top: 20px">
										|RequestorName| added an <a href="|ActivityLink|">activity</a> and is requesting for help
										below are the details of the request.<br>
										DETAILS: |Details|
									</td>
								</tr>
							</table>
						</td><br>
					</tr>`

	replacer := strings.NewReplacer(
		"|RequestorName|", activityHelpEmail.RequestorName,
		"|ActivityLink|", activityHelpEmail.ActivityLink,
		"|Details|", activityHelpEmail.Details,
	)

	htmlBody := buildHtmlBody(bodyTemplate, replacer)

	subject := fmt.Sprintf("%s - Request for Help", activityHelpEmail.HelpType)

	to := []string{
		os.Getenv("EMAIL_SUPPORT"),
	}

	err := es.SendEmail(to, nil, subject, htmlBody, Html, true)
	if err != nil {
		return err
	}
	return nil
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

// Email Template
func headStyle() string {
	return `<head>
				<style>
					table, th, tr, td {
						border: 0;
						border-collapse: collapse;
						vertical-align: middle;
					}
					.thead {
						padding: 15px;
					}
					.center-table {
						text-align: -webkit-center;
					}
					.margin-auto {
						margin: auto;
					}
					.border-top {
						border-top: 1px rgb(204, 204, 204) solid;
						border-collapse: separate;
					}
				</style>
			</head>`
}

func header() string {
	return `<tr>
                <th class="center-table">
                    <table style="width: 100%; max-width: 700px;" class="margin-auto">
                        <tr style="background-color: #ff5800">
                            <td class="thead" style="width: 95px">
                                <img src="data:image/png;base64,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" />
                            </td>
                            <td class="thead" style="font-family: SegoeUI, sans-serif;font-size: 14px;color: white;padding-bottom: 10px;">
                                Community
                            </td>
                        </tr>
                    </table>
                </th>
            </tr>`
}

func buildHtmlBody(bodyTemplate string, replacer *strings.Replacer) string {
	replacedBodyTemplate := replacer.Replace(bodyTemplate)
	htmlBody := `<html>`
	htmlBody += headStyle()
	htmlBody += `<body>`
	htmlBody += `<table style="width: 100%">`
	htmlBody += header()
	htmlBody += replacedBodyTemplate
	htmlBody += `</table>`
	htmlBody += `</body>`
	htmlBody += `</html>`
	return htmlBody
}
