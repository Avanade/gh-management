package main

import (
	"fmt"
	"log"
	githubAPI "main/pkg/github"
	session "main/pkg/session"
	rtAzure "main/routes/login/azure"
	rtGithub "main/routes/login/github"
	rtPages "main/routes/pages"

	rtApis "main/routes/pages/api"
	rtCommunity "main/routes/pages/community"
	rtProjects "main/routes/pages/projects"
	"net/http"

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
	githubAPI.CreateClient()

	mux := mux.NewRouter()
	mux.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./public/"))))
	mux.Handle("/", loadAzAuthPage(rtPages.HomeHandler))
	mux.Handle("/error/ghlogin", loadAzAuthPage(rtPages.GHLoginRequire))
	mux.Handle("/projects/new", loadAzGHAuthPage(rtProjects.ProjectsNewHandler))
	mux.Handle("/community/new", loadAzGHAuthPage(rtCommunity.CommunityHandler))
	mux.Handle("/community/{id}", loadAzGHAuthPage(rtCommunity.CommunityHandler))
	mux.Handle("/community/getcommunity/{id}", loadAzGHAuthPage(rtCommunity.GetUserCommunity))
	mux.Handle("/communities/list", loadAzGHAuthPage(rtCommunity.CommunitylistHandler))
	mux.Handle("/community", loadAzGHAuthPage(rtCommunity.GetUserCommunitylist))
	//mux.HandleFunc("/api/community", rtApis.CommunityAPIHandler)
	mux.Handle("/api/community", loadAzGHAuthPage(rtApis.CommunityAPIHandler))
	mux.HandleFunc("/api/communitySponsors", rtApis.CommunitySponsorsAPIHandler)
	mux.HandleFunc("/api/CommunitySponsorsPerCommunityId/{id}", rtApis.CommunitySponsorsPerCommunityId)
	mux.Handle("/projects/my", loadAzGHAuthPage(rtProjects.MyProjects))
	mux.Handle("/projects", loadAzGHAuthPage(rtProjects.GetUserProjects))
	mux.Handle("/projects/{id}", loadAzGHAuthPage(rtProjects.GetRequestStatusByProject))
	mux.Handle("/api/allusers", loadAzAuthPage(rtApis.GetAllUserFromActiveDirectory))
	mux.HandleFunc("/login/azure", rtAzure.LoginHandler)
	mux.HandleFunc("/login/azure/callback", rtAzure.CallbackHandler)
	mux.HandleFunc("/logout/azure", rtAzure.LogoutHandler)
	mux.HandleFunc("/login/github", rtGithub.GithubLoginHandler)
	mux.HandleFunc("/login/github/callback", rtGithub.GithubCallbackHandler)
	mux.HandleFunc("/login/github/force", rtGithub.GithubForceSaveHandler)
	mux.HandleFunc("/logout/github", rtGithub.GitHubLogoutHandler)
	mux.HandleFunc("/approvals/callback", rtProjects.UpdateApprovalStatus)
	mux.NotFoundHandler = http.HandlerFunc(rtPages.NotFoundHandler)

	//loadAzAuthPage()

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
