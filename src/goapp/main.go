package main

import (
	"fmt"
	"log"
	session "main/pkg/session"
	rtApi "main/routes/api"
	rtAzure "main/routes/login/azure"
	rtGithub "main/routes/login/github"
	rtPages "main/routes/pages"
	rtAdmin "main/routes/pages/admin"
	rtCommunity "main/routes/pages/community"
	rtProjects "main/routes/pages/projects"
	rtSearch "main/routes/pages/search"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"

	ev "main/pkg/envvar"

	"github.com/codegangsta/negroni"
	"github.com/joho/godotenv"
)

func main() {
	// Set environment variables
	err := godotenv.Load()
	if err != nil {
		log.Print(err.Error())
	}

	// Create session and GitHubClient
	session.InitializeSession()

	mux := mux.NewRouter()
	mux.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./public/"))))
	mux.Handle("/", loadAzAuthPage(rtPages.HomeHandler))
	mux.Handle("/error/ghlogin", loadAzAuthPage(rtPages.GHLoginRequire))
	mux.Handle("/projects/new", loadAzGHAuthPage(rtProjects.ProjectsNewHandler))

	mux.Handle("/search/{searchText}", loadAzGHAuthPage(rtSearch.GetSearchResults))
	mux.Handle("/search", loadAzGHAuthPage(rtSearch.SearchHandler))
	mux.Handle("/search/all/", loadAzGHAuthPage(rtSearch.GetAllResults))
	mux.Handle("/search/name/", loadAzGHAuthPage(rtSearch.GetResultsByName))
	mux.Handle("/search/description/", loadAzGHAuthPage(rtSearch.GetResultsByDescription))
	mux.Handle("/community/new", loadAzGHAuthPage(rtCommunity.CommunityHandler))
	mux.Handle("/community/{id}", loadAzGHAuthPage(rtCommunity.CommunityHandler))
	mux.Handle("/community/getcommunity/{id}", loadAzGHAuthPage(rtCommunity.GetUserCommunity))
	mux.Handle("/communities/list", loadAzGHAuthPage(rtCommunity.CommunitylistHandler))
	mux.Handle("/community", loadAzGHAuthPage(rtCommunity.GetUserCommunitylist))

	mux.Handle("/projects", loadAzGHAuthPage(rtProjects.Projects))
	mux.Handle("/community/{id}/onboarding", loadAzGHAuthPage(rtCommunity.CommunityOnBoarding))
	mux.HandleFunc("/login/azure", rtAzure.LoginHandler)
	mux.HandleFunc("/login/azure/callback", rtAzure.CallbackHandler)
	mux.HandleFunc("/logout/azure", rtAzure.LogoutHandler)
	mux.HandleFunc("/login/github", rtGithub.GithubLoginHandler)
	mux.HandleFunc("/login/github/callback", rtGithub.GithubCallbackHandler)
	mux.HandleFunc("/login/github/force", rtGithub.GithubForceSaveHandler)
	mux.HandleFunc("/logout/github", rtGithub.GitHubLogoutHandler)

	muxApi := mux.PathPrefix("/api").Subrouter()
	mux.Handle("/allusers", loadAzAuthPage(rtApi.GetAllUserFromActiveDirectory))
	muxApi.Handle("/activity/type", loadAzGHAuthPage(rtApi.GetActivityTypes)).Methods("GET")
	muxApi.Handle("/activity/type", loadAzGHAuthPage(rtApi.CreateActivityType)).Methods("POST")
	muxApi.Handle("/community", loadAzGHAuthPage(rtApi.CommunityAPIHandler))
	muxApi.Handle("/communitySponsors", loadAzGHAuthPage(rtApi.CommunitySponsorsAPIHandler))
	muxApi.Handle("/CommunitySponsorsPerCommunityId/{id}", loadAzGHAuthPage(rtApi.CommunitySponsorsPerCommunityId))
	muxApi.Handle("/CommunityTagPerCommunityId/{id}", loadAzGHAuthPage(rtApi.CommunityTagPerCommunityId))
	muxApi.Handle("/community/onboarding/{id}", loadAzGHAuthPage(rtApi.GetCommunityOnBoardingInfo)).Methods("GET", "POST", "DELETE")
	muxApi.Handle("/community/all", loadAzAuthPage(rtApi.GetCommunities)).Methods("GET")
	muxApi.Handle("/community/{id}/members", loadAzAuthPage(rtApi.GetCommunityMembers)).Methods("GET")
	muxApi.Handle("/contributionarea", loadAzGHAuthPage(rtApi.CreateContributionAreas)).Methods("POST")
	muxApi.Handle("/contributionarea", loadAzGHAuthPage(rtApi.GetContributionAreas)).Methods("GET")
	muxApi.Handle("/projects/list", loadAzGHAuthPage(rtApi.GetUserProjects))
	muxApi.Handle("/projects/{id}", loadAzGHAuthPage(rtApi.GetRequestStatusByProject))
	muxApi.Handle("/projects/{project}/archive/{archive}/private/{private}", loadAzGHAuthPage(rtApi.ArchiveProject))
	muxApi.Handle("/projects/{project}/private/{private}/archive/{archive}", loadAzGHAuthPage(rtApi.SetVisibility))
	muxApi.Handle("/allusers", loadAzAuthPage(rtApi.GetAllUserFromActiveDirectory))
	muxApi.Handle("/allavanadeprojects", loadAzGHAuthPage(rtApi.GetAvanadeProjects))

	//API FOR APPROVAL TYPES
	muxApi.HandleFunc("/approval/types", rtApi.CreateApprovalType).Methods("POST")
	muxApi.HandleFunc("/approval/types/{id}", rtApi.EditApprovalTypeById).Methods("PUT")
	muxApi.HandleFunc("/approval/types", rtApi.GetApprovalTypes).Methods("GET")
	muxApi.HandleFunc("/approval/types/{id}", rtApi.GetApprovalTypeById).Methods("GET")

	muxAdmin := mux.PathPrefix("/admin").Subrouter()
	muxAdmin.Handle("/members", loadAzGHAuthPage(rtAdmin.ListCommunityMembers))

	mux.HandleFunc("/approvals/project/callback", rtProjects.UpdateApprovalStatusProjects).Methods("POST")
	mux.HandleFunc("/approvals/community/callback", rtProjects.UpdateApprovalStatusCommunity).Methods("POST")
	mux.NotFoundHandler = http.HandlerFunc(rtPages.NotFoundHandler)

	go checkFailedApprovalRequests()

	port := ev.GetEnvVar("PORT", "8080")
	fmt.Printf("Now listening on port %v\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), mux))

}

// Verifies authentication before loading the page.
func loadAzAuthPage(f func(w http.ResponseWriter, r *http.Request)) *negroni.Negroni {
	return negroni.New(
		negroni.HandlerFunc(session.IsAuthenticated),
		negroni.Wrap(http.HandlerFunc(f)),
	)
}

func loadAzGHAuthPage(f func(w http.ResponseWriter, r *http.Request)) *negroni.Negroni {
	return negroni.New(
		negroni.HandlerFunc(session.IsAuthenticated),
		negroni.HandlerFunc(session.IsGHAuthenticated),
		negroni.Wrap(http.HandlerFunc(f)),
	)
}

func checkFailedApprovalRequests() {
	// TIMER SERVICE
	freq := ev.GetEnvVar("APPROVALREQUESTS_RETRY_FREQ", "15")
	freqInt, _ := strconv.ParseInt(freq, 0, 64)
	if freq > "0" {
		for range time.NewTicker(time.Duration(freqInt) * time.Minute).C {
			go rtProjects.ReprocessRequestApproval()
			go rtCommunity.ReprocessRequestCommunityApproval()
		}
	}
}
