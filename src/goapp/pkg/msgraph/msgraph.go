package msgraph

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

func GetAzGroupIdByName(groupName string) (string, error) {
	accessToken, err := getToken()
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

// Search user by name and mail
func SearchUsers(search string) ([]User, error) {
	accessToken, err := getToken()
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
	query.Set("$search", fmt.Sprintf(`"displayName:%s" OR "mail:%s"`, search, search))
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

	var listUsersResponse ListUSersResponse
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

// Get all users from the active directory
func GetAllUsers() ([]User, error) {
	accessToken, err := getToken()
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

	var listUsersResponse ListUSersResponse
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
	accessToken, err := getToken()
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

// Get all users from the active directory
// string user accepts user id and user principal name
func IsGithubEnterpriseMember(user string) (bool, error) {
	accessToken, err := getToken()
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
	accessToken, err := getToken()
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

// Get Access Token for the Application
func getToken() (string, error) {

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

type TokenResponse struct {
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	ExtExpiresIn int    `json:"ext_expires_in"`
	AccessToken  string `json:"access_token"`
}

type ListUSersResponse struct {
	DataContext string `json:"@odata.context"`
	Value       []User `json:"value"`
}

type User struct {
	Name  string `json:"displayName"`
	Email string `json:"mail"`
}

type ADGroupsResponse struct {
	Value []ADGroup `json:"value"`
}

type ADGroup struct {
	Id   string `json:"id"`
	Name string `json:"displayName"`
}
