package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/data/aztables"
	"github.com/google/go-github/v42/github"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	githuboauth "golang.org/x/oauth2/github"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func inviteUser(id int64, orgName string) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	// Create a new invitation
	role := "direct_member"
	invite := &github.CreateOrgInvitationOptions{
		InviteeID: &id,
		Role:      &role,
		TeamID:    []int64{
			// Team ID for "Developers"
		},
	}

	_, _, err := client.Organizations.CreateOrgInvitation(ctx, orgName, invite)
	if err != nil {
		log.Printf("Error creating invitation: %v", err)
	}
}

func getGithubOrgMembership(username string) (isPublicMember bool, isPrivateMember bool) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	// check to see if the user is a member of the rog
	isPublicMember, _, err := client.Organizations.IsMember(ctx, os.Getenv("GITHUB_PUBLIC_ORG"), username)
	if err != nil {
		log.Printf("Error checking public membership: %v", err)
	}
	isInternalMember, _, err := client.Organizations.IsMember(ctx, os.Getenv("GITHUB_INTERNAL_ORG"), username)
	if err != nil {
		log.Printf("Error checking internal membership: %v", err)
	}
	return isPublicMember, isInternalMember
}

//createGithubRepository("gh_app_test", "ava-innersource", "test")
func createPrivateGithubRepository(name string, owner string, description string) {
	// create github oauth client from token
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	repoRequest := &github.TemplateRepoRequest{
		Name:        &name,
		Owner:       &owner,
		Description: &description,
		Private:     github.Bool(true),
	}

	_, _, err := client.Repositories.CreateFromTemplate(context.Background(), "avanade", "avanade-template", repoRequest)
	if err != nil {
		log.Fatalf("Error happened in getting creating repo from template. Err: %s", err)
	}
}

func handlePrivateRepo(w http.ResponseWriter, r *http.Request) {
	// create github oauth client from token
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	state := r.FormValue("state")

	// set up json response
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")

	resp := make(map[string]interface{})
	resp["processed"] = false
	if (state == "") || (!isValidGuid(state)) {
		resp["message"] = "Your repository was not created. If the app has been opened for more than an hour, reload the page."
	} else {
		userEmail := getUserEmail(state)
		userRecord, err := getUserRecord(userEmail)
		ghUsername := fmt.Sprint(userRecord.Properties["githubUsername"])
		if err != nil {
			log.Println("failed to resolve entity: %w", err)
			resp["message"] = "Unable to identify your GitHub username."
		} else {
			// get the repo name from the json body
			newRepoName := r.FormValue("new_repo_name")
			// get the repo description from the json body
			newRepoDescription := r.FormValue("new_repo_description")

			if newRepoName == "" {
				log.Fatal("new repo name is empty")
			}
			log.Printf("Requested by %s", userEmail)
			// check to see if github repository already exists
			repo, _, err := client.Repositories.Get(ctx, os.Getenv("GITHUB_INTERNAL_ORG"), newRepoName)
			if err != nil && !strings.Contains(fmt.Sprintf("%v", err), "404 Not Found []") {
				log.Printf("Error checking for existing repository: %v", err)
				resp["message"] = "Unable to create your repository. Please try again later."
			} else if repo != nil {
				log.Printf("Repo already exists: %v", repo)
				resp["message"] = "A repository by that name already exists."
			} else {
				log.Printf("Repo does not exist, creating it")
				createPrivateGithubRepository(newRepoName, os.Getenv("GITHUB_INTERNAL_ORG"), newRepoDescription)
				resp["message"] = "Your repository has been created."
				resp["processed"] = true
				// set github repository visibility to private
				_, _, err := client.Repositories.Edit(ctx, os.Getenv("GITHUB_INTERNAL_ORG"), newRepoName, &github.Repository{
					Visibility: github.String("internal"),
				})
				if err != nil {
					log.Printf("Error setting visibility: %v", err)
					resp["message"] = "Repository created, although there was an error setting the proper visibility."
				}
				// add user to github repository in the admin role
				collaboratorOptions := &github.RepositoryAddCollaboratorOptions{
					Permission: "admin",
				}
				_, _, err = client.Repositories.AddCollaborator(ctx, os.Getenv("GITHUB_INTERNAL_ORG"), newRepoName, ghUsername, collaboratorOptions)
				if err != nil {
					log.Printf("Error adding user in admin role: %v", err)
					resp["message"] = "Repository created, although there was an error adding your username to it as the admin."
				}
			}
		}
	}
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	w.Write(jsonResp)
}

func handleTest(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")

	resp := make(map[string]interface{})

	resp["processed"] = true
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	w.Write(jsonResp)
	return
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
		log.Println("Error getting entity getting user email")
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
		log.Println("Error getting entity in GUID check")
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

func getUserRecord(userEmail string) (aztables.EDMEntity, error) {
	var err error
	connStr := os.Getenv("TABLE_CONNECTION_STRING")
	serviceClient, err := aztables.NewServiceClientFromConnectionString(connStr, nil)
	if err != nil {
		log.Println(err)
	}
	table := serviceClient.NewClient("users")
	ctx := context.TODO()
	fmt.Println("Getting user record for: " + userEmail)
	entity, err := table.GetEntity(ctx, "avanade", userEmail, nil)
	if err != nil {
		log.Println("Error getting entity in user record")
	}
	var entityResult aztables.EDMEntity
	err = json.Unmarshal(entity.Value, &entityResult)
	if err != nil {
		log.Println("failed to unmarshal entity: %w", err)
	}
	return entityResult, err
}

// /api/ebonded. Called by PowerApp on load.
// Displays if the user is currently connected to GitHub or not.
func handleApiEbondedStatus(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	state := r.FormValue("state")
	userEmail := getUserEmail(state)
	ebonded := false
	entityResult, err := getUserRecord(userEmail)
	if err != nil {
		log.Println("failed to resolve entity: %w", err)
	} else {
		ebonded = true
	}
	resp := make(map[string]interface{})
	if ebonded {
		resp["ebonded"] = true
		resp["githubId"] = entityResult.Properties["githubId"]
		resp["githubUsername"] = entityResult.Properties["githubUsername"]
		isPublic, isPrivate := getGithubOrgMembership(fmt.Sprint(entityResult.Properties["githubUsername"]))
		resp["isPublicMember"] = isPublic
		resp["isPrivateMember"] = isPrivate
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

func handleInvite(w http.ResponseWriter, r *http.Request, orgName string) {
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	state := r.FormValue("state")
	userEmail := getUserEmail(state)
	entityResult, err := getUserRecord(userEmail)
	if err != nil {
		log.Println("failed to unmarshal entity from user table: %w", err)
	}
	resp := make(map[string]interface{})
	githubId, _ := strconv.ParseInt(fmt.Sprint(entityResult.Properties["githubId"]), 10, 64)

	inviteUser(githubId, orgName)
	resp["processed"] = true
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal with github ID. Err: %s", err)
	}
	w.Write(jsonResp)
	return
}

func handleInvitePublic(w http.ResponseWriter, r *http.Request) {
	handleInvite(w, r, os.Getenv("GITHUB_PUBLIC_ORG"))
	return
}

func handleInviteInternal(w http.ResponseWriter, r *http.Request) {
	handleInvite(w, r, os.Getenv("GITHUB_INTERNAL_ORG"))
	return
}

func main() {
	err := godotenv.Load()
	log.Printf("Loaded as %v\n", os.Getenv("ENV"))
	if err != nil && os.Getenv("ENV") == "local" {
		log.Fatal("Error loading .env file")
	}
	// TODO: Periodically check github username hasn't been changed
	// TODO: Pass state and approval directly from automate so user can't listen to state keys
	http.HandleFunc("/", handleMain)
	http.HandleFunc("/login", handleGitHubLogin)
	http.HandleFunc("/github_oauth_cb", handleGitHubCallback)
	http.HandleFunc("/api/ebonded", handleApiEbondedStatus)
	http.HandleFunc("/api/invite_public", handleInvitePublic)
	http.HandleFunc("/api/invite_private", handleInviteInternal)
	http.HandleFunc("/api/create_private_repo", handlePrivateRepo)
	http.HandleFunc("/api/test", handleTest)
	fmt.Print("Started running on http://127.0.0.1:8000\n")
	fmt.Println(http.ListenAndServe(":8000", nil))
}
