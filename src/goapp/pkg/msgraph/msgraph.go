package msgraph

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

type TokenResponse struct {
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	ExtExpiresIn int    `json:"ext_expires_in"`
	AccessToken  string `json:"access_token"`
}

type ListUsersResponse struct {
	DataContext string `json:"@odata.context"`
	Value       []User `json:"value"`
}

type User struct {
	Name       string   `json:"displayName"`
	Email      string   `json:"mail"`
	OtherMails []string `json:"otherMails"`
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
		Timeout: time.Second * 10,
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
		Timeout: time.Second * 10,
	}

	urlPath := `https://graph.microsoft.com/v1.0/users`
	URL, errURL := url.Parse(urlPath)
	if err != nil {
		return nil, errURL
	}
	query := URL.Query()
	query.Set("$select", "displayName,otherMails,mail")
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
		if user.Email != "" || len(user.OtherMails) > 0 {
			if user.Email == "" && len(user.OtherMails) > 0 {
				user.Email = user.OtherMails[0]
			}
			users = append(users, user)
		}
	}

	return users, nil
}

func GetAllUsers() ([]User, error) {
	accessToken, err := GetToken()
	if err != nil {
		return nil, err
	}

	client := &http.Client{
		Timeout: time.Second * 10,
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
		Timeout: time.Second * 10,
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
		Timeout: time.Second * 10,
	}

	urlPath := fmt.Sprintf("https://graph.microsoft.com/v1.0/users/%s/checkMemberGroups", user)

	groupId, err := GetAzGroupIdByName(os.Getenv("GH_AZURE_AD_GROUP"))
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
	response, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer response.Body.Close()

	var data struct {
		Value []string `json:"value"`
	}

	errDecode := json.NewDecoder(response.Body).Decode(&data)
	if errDecode != nil {
		fmt.Print(err)
	}

	for _, v := range data.Value {
		if v == groupId {
			return true, nil
		}
	}
	return false, nil
}

func IsUserAdmin(user string) (bool, error) {
	accessToken, err := GetToken()
	if err != nil {
		return false, err
	}

	client := &http.Client{
		Timeout: time.Second * 10,
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
		Timeout: time.Second * 10,
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

	urlPath := fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/v2.0/token", os.Getenv("TENANT_ID"))
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	data := url.Values{}
	data.Set("client_id", os.Getenv("CLIENT_ID"))
	data.Set("scope", "https://graph.microsoft.com/.default")
	data.Set("client_secret", os.Getenv("CLIENT_SECRET"))
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

func IsUserExist(userPrincipalName string) (bool, error) {
	accessToken, err := GetToken()
	if err != nil {
		return false, err
	}

	client := &http.Client{
		Timeout: time.Second * 10,
	}

	urlPath := `https://graph.microsoft.com/v1.0/users`
	URL, errURL := url.Parse(urlPath)
	if err != nil {
		return false, errURL
	}
	query := URL.Query()
	query.Set("$select", "displayName")
	query.Set("$search", fmt.Sprintf(`"displayName:%[1]s" OR "otherMails:%[1]s" OR "mail:%[1]s" OR "userPrincipalName:%[1]s"`, userPrincipalName))
	URL.RawQuery = query.Encode()

	req, err := http.NewRequest("GET", URL.String(), nil)
	if err != nil {
		return false, err
	}
	req.Header.Add("Authorization", "Bearer "+accessToken)
	req.Header.Add("ConsistencyLevel", "eventual")
	response, err := client.Do(req)
	if err != nil {
		return false, err
	}
	if response.StatusCode != 200 {
		return false, errors.New(response.Status)
	}
	defer response.Body.Close()

	var listUsersResponse ListUsersResponse

	err = json.NewDecoder(response.Body).Decode(&listUsersResponse)
	if err != nil {
		return false, err
	}

	isMember := len(listUsersResponse.Value) > 0

	return isMember, nil
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
		Timeout: time.Second * 10,
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
		Timeout: time.Second * 10,
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
		Timeout: time.Second * 10,
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
