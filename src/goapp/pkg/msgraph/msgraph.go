package msgraph

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

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

	result, err := GetAllGroups()

	if result != nil {
		return nil, err
	}

	return users, nil
}

// Get all users from the active directory
func GetAllGroups() ([]User, error) {
	accessToken, err := getToken()
	if err != nil {
		return nil, err
	}

	client := &http.Client{
		Timeout: time.Second * 10,
	}

	urlPath := "https://graph.microsoft.com/v1.0/users/38ac8fe1-3b2a-41b8-bcf4-dfff7afb7fa7/checkMemberGroups"

	groupIds := []string{
		"a50dadf6-0ff8-4bc8-8dd0-c26e84860e95",
	}

	postBody, _ := json.Marshal(map[string]interface{}{
		"groupIds": groupIds,
	})

	reqBody := bytes.NewBuffer(postBody)

	req, err := http.NewRequest("POST", urlPath, reqBody)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer "+accessToken)
	req.Header.Add("Content-Type", "application/json")
	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Print(err)
	}
	bodyString := string(bodyBytes)
	fmt.Print(bodyString)
	return nil, nil
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
