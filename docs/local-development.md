## Local Development

1. Set up the environment variables.

**Environment Variables**


| Variable                          | Description                                       |
| -----------                       | -----------                                       |
| PORT                              | Port you like to use. Default: 8080               |
| TENANT_ID                         | Azure tenant ID for authentication                |
| CLIENT_ID                         | Azure application ID for authentication           |
| CLIENT_SECRET                     | Azure app reg secret for authentication           |
| GH_CLIENT_ID                      | GitHub OAuth app ID for SSO                       |
| GH_CLIENT_SECRET                  | GitHub OAuth app secret for SSO                   |
| GH_TOKEN                          | Token of GH user that is an owner of innersource and public organizations|
| GH_ORG_INNERSOURCE                | Name of GitHub innerource organization            |
| GH_ORG_OPENSOURCE                 | Name of GitHub public organization                |
| GH_REPO_TEMPLATE_NAME             | Repo Template where new repos will be based on    |
| GH_REPO_TEMPLATE                  | Owner of Repo Template                            |
| GHMGMTDB_CONNECTION_STRING        | Database connection string                        |
| APPROVAL_SYSTEM_APP_URL           | Domain of approval system                         |
| APPROVAL_SYSTEM_APP_ID            | App GUID from approval system                     |
| APPROVAL_SYSTEM_APP_MODULE_PROJECTS     | Module GUID from approval system            |
| APPROVAL_SYSTEM_APP_MODULE_COMMUNITY    | Module GUID from approval system            |
| APPROVALREQUESTS_RETRY_FREQ       | Delay between each retrys to create approval request. Default: 15 (min)  |
| APPROVALREQUESTS_RETRY_FREQ       | Delay between each retrys to create approval request. Default: 15 (min)  |
| GH_AZURE_AD_GROUP                 | Azure AD group for users with Visual Studio Subscription  |
| GH_AZURE_AD_ADMIN_GROUP           | Azure AD group for community portal admins        |
| EMAIL_SUPPORT                     | OSPO email address                                |
| EMAIL_ENDPOINT                    | Endpoint of the service used to send an email     |
| SUMMARY_REPORT_TRIGGER            | UTC Hour to trigger report of requested repos eg. 8 for 8am, 13 for 1pm  |
| EMAIL_SUMMARY_REPORT              | Recipient of email summary report                 |
| CONTENT_SECURITY_POLICY           | Content security policy for http requests         |
| SCHEME                            | http or https. Defaults to https, use http for local development  |
| IS_DEVELOPMENT                    | Defaults to false, use true for local development |
| APP_TITLE                         | Title to appear on browser's  title bar           |
| APP_LOGO_PATH                     | Path to favicon                                   |

2. On a terminal, navigate to /goapp/tailwind and run 
    `npx tailwindcss -i ./input.css -o ../public/css/output.css --watch`.
3. On another terminal, navigate to /goapp and run `go mod download`.
4. On the same terminal, run `go run .`.