package routes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	models "main/models"
	ghmgmtdb "main/pkg/ghmgmtdb"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func getHttpPostResponseStatus(url string, data interface{}, ch chan *http.Response) {
	jsonReq, err := json.Marshal(data)
	if err != nil {
		log.Println(err.Error())
		ch <- nil
	}
	res, err := http.Post(url, "application/json; charset=utf-8", bytes.NewBuffer(jsonReq))
	if err != nil {
		log.Println(err.Error())
		ch <- nil
	}
	ch <- res
}

func RequestCommunityApproval(id int64) error {
	communityApprovals := ghmgmtdb.PopulateCommunityApproval(id)

	for _, v := range communityApprovals {
		err := ApprovalSystemRequestCommunity(v)
		if err != nil {
			log.Println("ID:" + strconv.FormatInt(v.Id, 10) + " " + err.Error())
			return err
		}
	}
	return nil
}

func ReprocessRequestCommunityApproval() {
	projectApprovals := ghmgmtdb.GetFailedCommunityApprovalRequests()

	for _, v := range projectApprovals {
		go ApprovalSystemRequestCommunity(v)
	}

}

func ApprovalSystemRequestCommunity(data models.TypCommunityApprovals) error {

	url := os.Getenv("APPROVAL_SYSTEM_APP_URL")
	if url != "" {
		url = url + "/request"
		ch := make(chan *http.Response)
		// var res *http.Response

		bodyTemplate := `<p>Hi |ApproverUserPrincipalName|!</p>
		<p>|RequesterName| is requesting for a new |CommunityType| community and is now pending for approval.</p>
		<p>Below are the details:</p>
		<table>
			<tr>
				<td style="font-weight: bold;">Community Name<td>
				<td style="font-size:larger">|CommunityName|<td>
			</tr>
			<tr>
				<td style="font-weight: bold;">Url<td>
				<td style="font-size:larger">|CommunityUrl|<td>
			</tr>
			<tr>
				<td style="font-weight: bold;">Description<td>
				<td style="font-size:larger">|CommunityDescription|<td>
			</tr>
			<tr>
				<td style="font-weight: bold;">Notes<td>
				<td style="font-size:larger">|CommunityNotes|<td>
			</tr>
		</table>
		<p>For more information, send an email to <a href="mailto:|RequesterUserPrincipalName|">|RequesterUserPrincipalName|</a></p>
		`

		replacer := strings.NewReplacer("|ApproverUserPrincipalName|", data.ApproverUserPrincipalName,
			"|RequesterName|", data.RequesterName,
			"|CommunityType|", data.CommunityType,
			"|CommunityName|", data.CommunityName,
			"|CommunityUrl|", data.CommunityUrl,
			"|CommunityDescription|", data.CommunityDescription,
			"|CommunityNotes|", data.CommunityNotes,
			"|RequesterUserPrincipalName|", data.RequesterUserPrincipalName,
		)
		body := replacer.Replace(bodyTemplate)

		postParams := models.TypApprovalSystemPost{
			ApplicationId:       os.Getenv("APPROVAL_SYSTEM_APP_ID"),
			ApplicationModuleId: os.Getenv("APPROVAL_SYSTEM_APP_MODULE_COMMUNITY"),
			Email:               data.ApproverUserPrincipalName,
			Subject:             fmt.Sprintf("[GH-Management] New Community For Approval - %v", data.CommunityName),
			Body:                body,
			RequesterEmail:      data.RequesterUserPrincipalName,
		}

		go getHttpPostResponseStatus(url, postParams, ch)
		r := <-ch
		if r != nil {
			var res models.TypApprovalSystemPostResponse
			err := json.NewDecoder(r.Body).Decode(&res)
			if err != nil {
				log.Println(err.Error())
				return err
			}

			ghmgmtdb.CommunityApprovalUpdateGUID(data.Id, res.ItemId)
		}
	}
	return nil
}
