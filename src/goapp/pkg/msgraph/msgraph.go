package msgraph

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	db "main/pkg/ghmgmtdb"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	token tokenInfo
)

type tokenInfo struct {
	AccessToken string
	ExpiresIn   time.Time
}

type TokenResponse struct {
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	ExtExpiresIn int    `json:"ext_expires_in"`
	AccessToken  string `json:"access_token"`
}

type ListUsersResponse struct {
	DataContext  string `json:"@odata.context"`
	DataNextLink string `json:"@odata.nextLink"`
	Value        []User `json:"value"`
}

type User struct {
	UserPrincipalName string   `json:"userPrincipalName"`
	Name              string   `json:"displayName"`
	Email             string   `json:"mail"`
	OtherMails        []string `json:"otherMails"`
	UserType          string   `json:"userType"`
	AccountEnabled    bool     `json:"accountEnabled"`
}

type ADGroupsResponse struct {
	NextLink string    `json:"@odata.nextLink"`
	Value    []ADGroup `json:"value"`
}

type ADGroup struct {
	Id   string `json:"id"`
	Name string `json:"displayName"`
}

type AppRoleAssignmentResponse struct {
	Value []AppRoleAssignment `json:"value"`
}

type AppRoleAssignment struct {
	ResourceDisplayName string `json:"resourceDisplayName"`
}

func GetAzGroupIdByName(groupName string) (string, error) {
	accessToken, err := GetToken()
	if err != nil {
		return "", err
	}

	client := &http.Client{
		Timeout: time.Second * 90,
	}

	urlPath := fmt.Sprintf("https://graph.microsoft.com/v1.0/groups?$search=\"displayName:%s\"", groupName)

	req, err := http.NewRequest("GET", urlPath, nil)
	if err != nil {
		return "", err
	}
	req.Header.Add("Authorization", "Bearer "+accessToken)
	req.Header.Add("ConsistencyLevel", "eventual")
	response, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	var listGroupResponse ADGroupsResponse
	err = json.NewDecoder(response.Body).Decode(&listGroupResponse)
	if err != nil {
		return "", err
	}

	if len(listGroupResponse.Value) == 0 {
		return "", fmt.Errorf("No group found")
	}

	return listGroupResponse.Value[0].Id, nil
}

func SearchUsers(search string) ([]User, error) {
	accessToken, err := GetToken()
	if err != nil {
		return nil, err
	}

	client := &http.Client{
		Timeout: time.Second * 90,
	}

	urlPath := `https://graph.microsoft.com/v1.0/users`
	URL, errURL := url.Parse(urlPath)
	if err != nil {
		return nil, errURL
	}
	query := URL.Query()
	query.Set("$select", "displayName,otherMails,mail,userPrincipalName,userType")
	query.Set("$search", fmt.Sprintf(`"displayName:%s" OR "mail:%s" OR "otherMails:%s"`, search, search, search))
	URL.RawQuery = query.Encode()

	req, err := http.NewRequest("GET", URL.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+accessToken)
	req.Header.Add("ConsistencyLevel", "eventual")
	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var listUsersResponse ListUsersResponse
	err = json.NewDecoder(response.Body).Decode(&listUsersResponse)
	if err != nil {
		return nil, err
	}

	// Remove users without email address
	var users []User
	for _, user := range listUsersResponse.Value {
		if user.UserType == "Member" {
			user.Email = user.UserPrincipalName
		} else {
			if user.Email != "" || len(user.OtherMails) > 0 {
				if user.Email == "" && len(user.OtherMails) > 0 {
					user.Email = user.OtherMails[0]
				}
			}
		}
		users = append(users, user)
	}

	return users, nil
}

func GetAllUsers() ([]User, error) {
	accessToken, err := GetToken()
	if err != nil {
		return nil, err
	}

	client := &http.Client{
		Timeout: time.Second * 90,
	}

	urlPath := "https://graph.microsoft.com/v1.0/users"

	req, err := http.NewRequest("GET", urlPath, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+accessToken)
	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var listUsersResponse ListUsersResponse
	err = json.NewDecoder(response.Body).Decode(&listUsersResponse)
	if err != nil {
		return nil, err
	}

	// Remove users without email address
	var users []User
	for _, user := range listUsersResponse.Value {
		if user.Email != "" {
			users = append(users, user)
		}
	}

	return users, nil
}

func IsDirectMember(user string) (bool, error) {
	accessToken, err := GetToken()
	if err != nil {
		return false, err
	}

	client := &http.Client{
		Timeout: time.Second * 90,
	}

	urlPath := fmt.Sprintf("https://graph.microsoft.com/v1.0/users/%s", user)

	req, err := http.NewRequest("GET", urlPath, nil)
	if err != nil {
		return false, err
	}
	req.Header.Add("Authorization", "Bearer "+accessToken)
	response, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer response.Body.Close()

	var data struct {
		UserPrincipalName string `json:"userPrincipalName"`
	}

	errDecode := json.NewDecoder(response.Body).Decode(&data)
	if errDecode != nil {
		fmt.Print(err)
	}

	// #EXT# It checks if the user principal name is extension or direct.
	return !strings.Contains(data.UserPrincipalName, "#EXT#"), nil
}

func IsGithubEnterpriseMember(user string) (bool, error) {
	accessToken, err := GetToken()
	if err != nil {
		return false, err
	}

	client := &http.Client{
		Timeout: time.Second * 90,
	}

	adGroups, err := db.ADGroup_SelectAll()
	if err != nil {
		return false, err
	}

	urlPath := fmt.Sprintf("https://graph.microsoft.com/v1.0/users/%s/checkMemberGroups", user)
	isMember := false

	for x := 0; x < len(adGroups); x += 20 {
		end := 0
		if x+20 > len(adGroups) {
			end = len(adGroups)
		} else {
			end = x + 20
		}
		postBody, _ := json.Marshal(map[string]interface{}{
			"groupIds": adGroups[x:end],
		})

		reqBody := bytes.NewBuffer(postBody)

		req, err := http.NewRequest("POST", urlPath, reqBody)
		if err != nil {
			isMember = false
			break
		}

		req.Header.Add("Authorization", "Bearer "+accessToken)
		req.Header.Add("Content-Type", "application/json")
		response, err := client.Do(req)
		if err != nil {
			isMember = false
			break
		}
		defer response.Body.Close()

		var data struct {
			Value []string `json:"value"`
		}

		errDecode := json.NewDecoder(response.Body).Decode(&data)
		if errDecode != nil {
			fmt.Print(err)
			isMember = false
			break
		}

		if len(data.Value) > 0 {
			isMember = true
			break
		}
	}

	return isMember, nil
}

func IsUserAdmin(user string) (bool, error) {
	accessToken, err := GetToken()
	if err != nil {
		return false, err
	}

	client := &http.Client{
		Timeout: time.Second * 90,
	}

	urlPath := fmt.Sprintf("https://graph.microsoft.com/v1.0/users/%s/checkMemberGroups", user)

	groupId, err := GetAzGroupIdByName(os.Getenv("GH_AZURE_AD_ADMIN_GROUP"))
	if err != nil {
		return false, err
	}

	groupIds := []string{
		groupId,
	}

	postBody, _ := json.Marshal(map[string]interface{}{
		"groupIds": groupIds,
	})

	reqBody := bytes.NewBuffer(postBody)

	req, err := http.NewRequest("POST", urlPath, reqBody)
	if err != nil {
		return false, err
	}

	req.Header.Add("Authorization", "Bearer "+accessToken)
	req.Header.Add("Content-Type", "application/json")
	response2, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer response2.Body.Close()

	var data struct {
		Value []string `json:"value"`
	}

	err = json.NewDecoder(response2.Body).Decode(&data)
	if err != nil {
		return false, err
	}

	for _, v := range data.Value {
		if v == groupId {
			return true, nil
		}
	}
	return false, nil
}

func GetUserPhoto(user string) (bool, string, error) {
	accessToken, err := GetToken()
	if err != nil {
		return false, "", err
	}

	client := &http.Client{
		Timeout: time.Second * 90,
	}

	urlPath := fmt.Sprintf("https://graph.microsoft.com/v1.0/users/%s/photos/64x64/$value", user)

	req, err := http.NewRequest("GET", urlPath, nil)
	if err != nil {
		return false, "", err
	}

	req.Header.Add("Authorization", "Bearer "+accessToken)
	req.Header.Add("Content-Type", "application/json")
	response, err := client.Do(req)
	if err != nil {
		return false, "", err
	}
	defer response.Body.Close()
	if response.StatusCode == http.StatusNotFound {
		return false, "", nil
	}

	userPhotoBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return false, "", err
	}
	userPhotoBase64 := base64.StdEncoding.EncodeToString(userPhotoBytes)
	return true, userPhotoBase64, nil
}

func GetToken() (string, error) {
	if token.AccessToken != "" {
		if token.ExpiresIn.After(time.Now()) {
			return token.AccessToken, nil
		}
	}

	newToken, err := requestNewToken()
	if err != nil {
		log.Println(err.Error())
		return "", err
	}

	const ALLOWANCE_TIME_BEFORE_EXPIRATION = 99

	duration, _ := time.ParseDuration(fmt.Sprint(newToken.ExpiresIn-ALLOWANCE_TIME_BEFORE_EXPIRATION, "s"))

	expiresin := time.Now().Add(duration)

	token.AccessToken = newToken.AccessToken
	token.ExpiresIn = expiresin

	return token.AccessToken, nil
}

func requestNewToken() (*TokenResponse, error) {

	urlPath := fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/v2.0/token", os.Getenv("TENANT_ID"))
	client := &http.Client{
		Timeout: time.Second * 90,
	}

	data := url.Values{}
	data.Set("client_id", os.Getenv("CLIENT_ID"))
	data.Set("scope", "https://graph.microsoft.com/.default")
	data.Set("client_secret", os.Getenv("CLIENT_SECRET"))
	data.Set("grant_type", "client_credentials")
	encodedData := data.Encode()

	req, err := http.NewRequest("POST", urlPath, strings.NewReader(encodedData))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var tokenResponse TokenResponse
	err = json.NewDecoder(response.Body).Decode(&tokenResponse)
	if err != nil {
		return nil, err
	}

	return &tokenResponse, nil
}

func IsUserExist(userPrincipalName string) (isMember bool, isAccountEnabled bool, err error) {
	accessToken, err := GetToken()
	if err != nil {
		return
	}

	client := &http.Client{
		Timeout: time.Second * 90,
	}

	urlPath := `https://graph.microsoft.com/v1.0/users`
	URL, err := url.Parse(urlPath)
	if err != nil {
		return
	}
	query := URL.Query()
	query.Set("$select", "accountEnabled")
	query.Set("$search", fmt.Sprintf(`"displayName:%[1]s" OR "otherMails:%[1]s" OR "mail:%[1]s" OR "userPrincipalName:%[1]s"`, userPrincipalName))
	URL.RawQuery = query.Encode()

	req, err := http.NewRequest("GET", URL.String(), nil)
	if err != nil {
		return
	}
	req.Header.Add("Authorization", "Bearer "+accessToken)
	req.Header.Add("ConsistencyLevel", "eventual")
	response, err := client.Do(req)
	if err != nil {
		return
	}
	if response.StatusCode != 200 {
		err = errors.New(response.Status)
		return
	}
	defer response.Body.Close()

	var listUsersResponse ListUsersResponse

	err = json.NewDecoder(response.Body).Decode(&listUsersResponse)
	if err != nil {
		return
	}

	isMember = len(listUsersResponse.Value) > 0

	if isMember {
		isAccountEnabled = listUsersResponse.Value[0].AccountEnabled
	}

	return
}

func GetTeamsMembers(ChannelId string, token string) ([]User, error) {
	accessToken, err := GetToken()
	if err != nil {
		return nil, err
	}
	if token == "" {
		token = accessToken
	}
	client := &http.Client{
		Timeout: time.Second * 90,
	}

	urlPath := fmt.Sprintf("https://graph.microsoft.com/v1.0/groups/%s/members", ChannelId)

	req, err := http.NewRequest("GET", urlPath, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+token)
	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var listUsersResponse ListUsersResponse
	err = json.NewDecoder(response.Body).Decode(&listUsersResponse)
	if err != nil {
		return nil, err
	}

	var users []User
	for _, user := range listUsersResponse.Value {
		if user.Email != "" {
			users = append(users, user)
		}
	}

	return users, nil
}

func GetADGroups() ([]ADGroup, error) {
	accessToken, err := GetToken()
	if err != nil {
		return nil, err
	}

	client := &http.Client{
		Timeout: time.Second * 90,
	}

	urlPath := "https://graph.microsoft.com/v1.0/groups"

	req, err := http.NewRequest("GET", urlPath, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+accessToken)
	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var listGroupResponse ADGroupsResponse
	var finalList []ADGroup
	err = json.NewDecoder(response.Body).Decode(&listGroupResponse)
	if err != nil {
		return nil, err
	}
	finalList = append(finalList, listGroupResponse.Value...)
	nextLink := listGroupResponse.NextLink

	for nextLink != "" {
		req, err = http.NewRequest("GET", nextLink, nil)
		if err != nil {
			return nil, err
		}
		req.Header.Add("Authorization", "Bearer "+accessToken)
		response, err = client.Do(req)
		if err != nil {
			return nil, err
		}

		err = json.NewDecoder(response.Body).Decode(&listGroupResponse)
		if err != nil {
			return nil, err
		}
		finalList = append(finalList, listGroupResponse.Value...)
		if nextLink != listGroupResponse.NextLink {
			nextLink = listGroupResponse.NextLink
		} else {
			nextLink = ""
		}
	}

	return finalList, nil
}

func HasGitHubAccess(objectId string) (bool, error) {
	appReg := os.Getenv("GH_AZURE_AD_GROUP")
	accessToken, err := GetToken()
	if err != nil {
		return false, err
	}

	client := &http.Client{
		Timeout: time.Second * 90,
	}

	urlPath := fmt.Sprintf("https://graph.microsoft.com/v1.0/groups/%s/appRoleAssignments", objectId)

	req, err := http.NewRequest("GET", urlPath, nil)
	if err != nil {
		return false, err
	}
	req.Header.Add("Authorization", "Bearer "+accessToken)
	response, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer response.Body.Close()

	var appRoleAssignmentResponse AppRoleAssignmentResponse
	err = json.NewDecoder(response.Body).Decode(&appRoleAssignmentResponse)
	if err != nil {
		return false, err
	}

	for _, appRoleAssignment := range appRoleAssignmentResponse.Value {
		if appRoleAssignment.ResourceDisplayName == appReg {
			return true, nil
		}
	}
	return false, nil
}

func GetUsersByGroupId(groupId string) ([]User, error) {
	accessToken, err := GetToken()
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	urlPath := fmt.Sprintf(`https://graph.microsoft.com/v1.0/groups/%s/members`, groupId)
	URL, errURL := url.Parse(urlPath)
	if errURL != nil {
		return nil, errURL
	}

	query := URL.Query()
	query.Set("$select", "displayName,otherMails,mail,userPrincipalName,userType")
	URL.RawQuery = query.Encode()

	var users []User
	for {
		req, err := http.NewRequest("GET", URL.String(), nil)
		if err != nil {
			return nil, err
		}
		req.Header.Add("Authorization", "Bearer "+accessToken)
		req.Header.Add("ConsistencyLevel", "eventual")

		response, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		defer response.Body.Close()

		var listUsersResponse ListUsersResponse
		err = json.NewDecoder(response.Body).Decode(&listUsersResponse)
		if err != nil {
			return nil, err
		}

		// Append users to the list
		for _, user := range listUsersResponse.Value {
			if user.Email != "" || len(user.OtherMails) > 0 {
				users = append(users, user)
			}
		}

		// Check if there is a next page
		if listUsersResponse.DataNextLink == "" {
			break
		}
		URL, err = url.Parse(listUsersResponse.DataNextLink)
		if err != nil {
			return nil, err
		}
	}

	return users, nil
}
