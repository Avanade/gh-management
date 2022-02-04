package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/data/aztables"
	"github.com/google/go-github/github"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	githuboauth "golang.org/x/oauth2/github"
	"html/template"
	"log"
	"net/http"
	"os"
)

func getGithubUsers() {

}

func inviteUser(id int64, email string) {

}

func getGithubRepositories() {

}

func getUserEmail(stateString string) string {
	connStr := os.Getenv("TABLE_CONNECTION_STRING")
	serviceClient, err := aztables.NewServiceClientFromConnectionString(connStr, nil)
	if err != nil {
		log.Fatal(err)
	}
	table := serviceClient.NewClient("access")
	ctx := context.TODO()
	entity, err := table.GetEntity(ctx, "avanade", stateString, nil)
	if err != nil {
		log.Println("Error getting entity")
		return ""
	}
	var entityResult aztables.EDMEntity
	err = json.Unmarshal(entity.Value, &entityResult)
	if err != nil {
		log.Println("failed to unmarshal entity: %w", err)
		return ""
	}
	userEmail := entityResult.Properties["userEmail"]
	fmt.Printf("User email: %s\n", userEmail)
	return fmt.Sprint(userEmail)
}

func isValidGuid(stateString string) bool {
	connStr := os.Getenv("TABLE_CONNECTION_STRING")
	serviceClient, err := aztables.NewServiceClientFromConnectionString(connStr, nil)
	if err != nil {
		log.Fatal(err)
	}
	table := serviceClient.NewClient("access")
	ctx := context.TODO()
	entity, err := table.GetEntity(ctx, "avanade", stateString, nil)
	if err != nil {
		log.Println("Error getting entity")
		return false
	}
	var entityResult aztables.EDMEntity
	err = json.Unmarshal(entity.Value, &entityResult)
	if err != nil {
		log.Println("failed to unmarshal entity: %w", err)
		return false
	}
	userEmail := entityResult.Properties["userEmail"]
	fmt.Printf("User email: %s\n", userEmail)
	return true
}

func saveGitHubUser(stateString string, userEmail string, userName string, id int64) {
	fmt.Printf("Saving GitHub user: %s (%v)\n", userName, id)
	connStr := os.Getenv("TABLE_CONNECTION_STRING")
	serviceClient, err := aztables.NewServiceClientFromConnectionString(connStr, nil)
	if err != nil {
		log.Fatal(err)
	}
	table := serviceClient.NewClient("users")
	myEntity := aztables.EDMEntity{
		Entity: aztables.Entity{
			PartitionKey: "avanade",
			RowKey:       userEmail,
		},
		Properties: map[string]interface{}{
			"ebonded":        true,
			"githubId":       id,
			"githubUsername": userName,
		},
	}
	marshalled, err := json.Marshal(myEntity)
	if err != nil {
		log.Fatal(err)
	}
	// create an insertentityoptions for mergeentity
	insertEntityOptions := aztables.InsertEntityOptions{
		UpdateMode: aztables.MergeEntity,
	}
	_, err = table.InsertEntity(context.TODO(), marshalled, &insertEntityOptions)
	if err != nil {
		log.Fatal(err)
	}
}

func getOauthConfig() *oauth2.Config {
	oauthConf := &oauth2.Config{
		ClientID:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
		Scopes:       []string{"user:email", "repo"},
		Endpoint:     githuboauth.Endpoint,
	}
	return oauthConf
}

func handleError(w http.ResponseWriter, errorTitle string, errorMessage string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	t, _ := template.ParseFiles("templates/error.html")
	data := map[string]string{"PrivacyUrl": os.Getenv("PRIVACY_URL"), "ErrorTitle": errorTitle, "ErrorMessage": errorMessage}
	t.Execute(w, data)
}

func handleMain(w http.ResponseWriter, r *http.Request) {
	state := r.FormValue("state")
	if (state == "") || (!isValidGuid(state)) {
		handleError(w, "Invalid State", "The state parameter is invalid. Please try again.")
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	t, _ := template.ParseFiles("templates/login.html")
	data := map[string]string{"PrivacyUrl": os.Getenv("PRIVACY_URL"), "State": state}
	t.Execute(w, data)
}

// /login redirects to GitHub’s authorization page:
// GitHub will show authorization page to your user
// If the user authorizes your app, GitHub will re-direct to OAuth callback
func handleGitHubLogin(w http.ResponseWriter, r *http.Request) {
	oauthConf := getOauthConfig()
	state := r.FormValue("state")
	if (state == "") || (!isValidGuid(state)) {
		handleError(w, "Invalid State", "The state parameter is invalid. Please try again.")
		return
	}
	url := oauthConf.AuthCodeURL(state, oauth2.AccessTypeOnline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// /github_oauth_cb. Called by GitHub after authorization is granted
// This handler accesses the GitHub token, and converts the token into a http client
// We then use that http client to list GitHub information about the user
func handleGitHubCallback(w http.ResponseWriter, r *http.Request) {
	state := r.FormValue("state")
	oauthConf := getOauthConfig()
	if !isValidGuid(state) {
		fmt.Printf("oauth string '%s', isn't valid", state)
		handleError(w, "Invalid State", "The state parameter passed back from GitHub is invalid. Please try again.")
	}

	code := r.FormValue("code")
	token, err := oauthConf.Exchange(oauth2.NoContext, code)
	if err != nil {
		fmt.Printf("oauthConf.Exchange() failed with '%s'\n", err)
		handleError(w, "GitHub Error", err.Error())
		return
	}

	oauthClient := oauthConf.Client(oauth2.NoContext, token)
	client := github.NewClient(oauthClient)
	user, _, err := client.Users.Get(oauth2.NoContext, "")
	if err != nil {
		fmt.Printf("client.Users.Get() faled with '%s'\n", err)
		handleError(w, "GitHub Error", err.Error())
		return
	}
	fmt.Printf("Logged in as GitHub user: %s (%v)\n", *user.Login, *user.ID)
	userEmail := getUserEmail(state)
	saveGitHubUser(state, userEmail, *user.Login, *user.ID)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	t, _ := template.ParseFiles("templates/logged-in.html")
	t.Execute(w, map[string]string{"PrivacyUrl": os.Getenv("PRIVACY_URL"), "Username": *user.Login})
}

// /api/ebonded. Called by PowerApp on load.
// Displays if the user is currently connected to GitHub or not.
func handleApiEbondedStatus(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	state := r.FormValue("state")
	userEmail := getUserEmail(state)
	ebonded := false
	connStr := os.Getenv("TABLE_CONNECTION_STRING")
	serviceClient, err := aztables.NewServiceClientFromConnectionString(connStr, nil)
	if err != nil {
		log.Println(err)
	}
	table := serviceClient.NewClient("users")
	ctx := context.TODO()
	entity, err := table.GetEntity(ctx, "avanade", userEmail, nil)
	if err != nil {
		log.Println("Error getting entity")
	}
	var entityResult aztables.EDMEntity
	err = json.Unmarshal(entity.Value, &entityResult)
	if err != nil {
		log.Println("failed to unmarshal entity: %w", err)
	} else {
		ebonded = true
	}
	resp := make(map[string]interface{})
	if ebonded {
		resp["ebonded"] = true
		resp["ghJoined"] = entityResult.Properties["ghJoined"]
		resp["gheJoined"] = entityResult.Properties["gheJoined"]
		resp["githubId"] = entityResult.Properties["githubId"]
		resp["githubUsername"] = entityResult.Properties["githubUsername"]
	} else {
		resp["ebonded"] = false
	}

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	w.Write(jsonResp)
	return
}

func main() {
	err := godotenv.Load()
	log.Printf("Loaded as %v\n", os.Getenv("ENV"))
	if err != nil && os.Getenv("ENV") == "local" {
		log.Fatal("Error loading .env file")
	}

	http.HandleFunc("/", handleMain)
	http.HandleFunc("/login", handleGitHubLogin)
	http.HandleFunc("/github_oauth_cb", handleGitHubCallback)
	http.HandleFunc("/api/ebonded", handleApiEbondedStatus)
	fmt.Print("Started running on http://127.0.0.1:8000\n")
	fmt.Println(http.ListenAndServe(":8000", nil))
}
