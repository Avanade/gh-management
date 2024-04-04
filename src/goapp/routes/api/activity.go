package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"main/pkg/appinsights_wrapper"
	"main/pkg/email"
	"main/pkg/envvar"
	db "main/pkg/ghmgmtdb"
	"main/pkg/notification"
	"main/pkg/session"

	"github.com/gorilla/mux"
	"github.com/microsoft/ApplicationInsights-Go/appinsights/contracts"
)

type ActivitiesDto struct {
	Data  interface{} `json:"data"`
	Total int         `json:"total"`
}

type ActivityDto struct {
	Name        string  `json:"name"`
	Url         string  `json:"url"`
	Date        string  `json:"date"`
	Type        ItemDto `json:"type"`
	CommunityId int     `json:"communityid"`
	CreatedBy   string
	ModifiedBy  string

	PrimaryContributionArea     ItemDto   `json:"primarycontributionarea"`
	AdditionalContributionAreas []ItemDto `json:"additionalcontributionareas"`

	Help HelpDto `json:"help"`
}

type HelpDto struct {
	ItemDto
	Details string `json:"details"`
}

type ItemDto struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func GetActivities(w http.ResponseWriter, r *http.Request) {
	sessionaz, _ := session.Store.Get(r, "auth-session")
	iprofile := sessionaz.Values["profile"]
	profile := iprofile.(map[string]interface{})
	username := fmt.Sprint(profile["preferred_username"])

	var result ActivitiesDto

	params := r.URL.Query()

	if params.Has("offset") && params.Has("filter") {
		filter, _ := strconv.Atoi(params["filter"][0])
		offset, _ := strconv.Atoi(params["offset"][0])
		search := params["search"][0]
		orderby := params["orderby"][0]
		ordertype := params["ordertype"][0]
		result = ActivitiesDto{
			Data:  db.CommunitiesActivities_Select_ByOffsetAndFilterAndCreatedBy(offset, filter, orderby, ordertype, search, username),
			Total: db.CommunitiesActivities_TotalCount_ByCreatedBy(username, search),
		}
	} else {
		result = ActivitiesDto{
			Data:  db.CommunitiesActivities_Select(),
			Total: db.CommunitiesActivities_TotalCount(),
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

func CreateActivity(w http.ResponseWriter, r *http.Request) {
	logger := appinsights_wrapper.NewClient()
	defer logger.EndOperation()

	sessionaz, _ := session.Store.Get(r, "auth-session")
	iprofile := sessionaz.Values["profile"]
	profile := iprofile.(map[string]interface{})
	username := fmt.Sprint(profile["preferred_username"])

	var body ActivityDto
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// CHECK ACTIVITY TYPE IF EXIST / INSERT IF NOT EXIST
	if body.Type.Id == 0 {
		id, err := db.ActivityTypes_Insert(body.Type.Name)
		if err != nil {
			logger.LogException(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		body.Type.Id = id
	}

	// COMMUNITY ACTIVITY
	communityActivityId, err := db.CommunitiesActivities_Insert(db.Activity{
		Name:        body.Name,
		Url:         body.Url,
		Date:        body.Date,
		TypeId:      body.Type.Id,
		CommunityId: body.CommunityId,
		CreatedBy:   username,
	})
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if body.Help.Id != 0 {
		if os.Getenv("NOTIFICATION_EMAIL_SUPPORT") == "" {
			logger.LogTrace("Activity does not send successfully because the notification support email environment is not set.", contracts.Error)
		} else {
			recipients := strings.Split(os.Getenv("NOTIFICATION_EMAIL_SUPPORT"), ",")

			scheme := envvar.GetEnvVar("SCHEME", "https")

			activityLink := fmt.Sprint(scheme, "://", r.Host, "/activities/view/", communityActivityId)

			messageBody := notification.ActivityAddedRequestForHelpMessageBody{
				Recipients:   recipients,
				ActivityLink: activityLink,
				UserName:     profile["name"].(string),
			}
			err = messageBody.Send()
			if err != nil {
				logger.LogException(err)
			}

			errHelp := processHelp(communityActivityId, activityLink, username, profile["name"].(string), body.Help)
			if errHelp != nil {
				logger.LogException(err)
			}
		}
	}

	// PRIMARY CONTRIBUTION AREA
	err = insertCommunityActivitiesContributionArea(body.PrimaryContributionArea, db.CommunityActivitiesContributionAreas{
		CommunityActivityId: communityActivityId,
		ContributionAreaId:  body.PrimaryContributionArea.Id,
		IsPrimary:           true,
		CreatedBy:           username,
	})
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// ADDITIONAL CONTRIBUTION AREA
	for _, contributionArea := range body.AdditionalContributionAreas {
		err = insertCommunityActivitiesContributionArea(contributionArea, db.CommunityActivitiesContributionAreas{
			CommunityActivityId: communityActivityId,
			ContributionAreaId:  contributionArea.Id,
			IsPrimary:           false,
			CreatedBy:           username,
		})
		if err != nil {
			logger.LogException(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	fmt.Fprint(w, body)
}

func GetActivityById(w http.ResponseWriter, r *http.Request) {
	logger := appinsights_wrapper.NewClient()
	defer logger.EndOperation()

	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	result, err := db.CommunitiesActivities_Select_ById(id)
	if err != nil {
		logger.LogTrace(err.Error(), contracts.Error)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

func insertCommunityActivitiesContributionArea(ca ItemDto, caca db.CommunityActivitiesContributionAreas) error {
	if ca.Id == 0 {
		id, err := db.ContributionAreas_Insert(ca.Name, caca.CreatedBy)
		if err != nil {
			return err
		}

		caca.ContributionAreaId = id
	}

	_, err := db.CommunityActivitiesContributionAreas_Insert(db.CommunityActivitiesContributionAreas{
		CommunityActivityId: caca.CommunityActivityId,
		ContributionAreaId:  caca.ContributionAreaId,
		IsPrimary:           caca.IsPrimary,
		CreatedBy:           caca.CreatedBy,
	})
	if err != nil {
		return err
	}

	return nil
}

func processHelp(activityId int, activityLink, requestorEmail string, requestorName string, h HelpDto) error {
	// INSERT
	_, err := db.CommunityActivitiesHelpTypes_Insert(activityId, h.Id, h.Details)
	if err != nil {
		return err
	}

	bodyTemplate := `
		<html>
			<head>
				<style>
					table,
					th,
					tr,
					td {
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
			</head>

			<body>
				<table style="width: 100%">
					<tr>
						<th class="center-table">
							<table style="width: 100%; max-width: 700px;" class="margin-auto">
								<tr style="background-color: #ff5800">
									<td class="thead" style="width: 95px">
									<img
										src="data:image/png;base64,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" />
									</td>
									<td class="thead" style="
										font-family: SegoeUI, sans-serif;
										font-size: 14px;
										color: white;
										padding-bottom: 10px;
										">
									Community
									</td>
								</tr>
							</table>
						</th>
					</tr>
					<tr>
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
						<tr class="center-table">
							<table style="width: 100%; max-width: 700px;" class="margin-auto">
								<tr>
									<td style="padding-top: 20px">
										|RequestorName| added an <a href=\"|ActivityLink|\">activity</a> and is requesting for help below are the details of the request.<br>
										DETAILS: |Details|
									</td>
								</tr>
							</table>
						</tr>
					</tr>
				</table>
				<br>
			</body>

		</html>
	`

	replacer := strings.NewReplacer(
		"|RequestorName|", requestorName,
		"|ActivityLink|", activityLink,
		"|Details|", h.Details,
	)
	body := replacer.Replace(bodyTemplate)

	// SEND EMAIL
	m := email.Message{
		Subject: fmt.Sprintf("%s : Request for Help", h.Name),
		Body: email.Body{
			Content: body,
			Type:    email.HtmlMessageType,
		},
		ToRecipients: []email.Recipient{
			{
				Email: os.Getenv("EMAIL_SUPPORT"),
			},
		},
		CcRecipients: []email.Recipient{
			{
				Email: requestorEmail,
			},
		},
	}

	errEmail := email.SendEmail(m, false)
	if errEmail != nil {
		return errEmail
	}

	// NO ERROR
	return nil
}
