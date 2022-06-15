package main

import (
	"fmt"
	"log"
	githubAPI "main/pkg/github"
	session "main/pkg/session"
	rtAzure "main/routes/login/azure"
	rtGithub "main/routes/login/github"
	rtPages "main/routes/pages"
	rtProjects "main/routes/pages/projects"
	rtSearch "main/routes/pages/search"
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
	mux.Handle("/projects/my", loadAzGHAuthPage(rtProjects.MyProjects))
	mux.Handle("/projects", loadAzGHAuthPage(rtProjects.GetUserProjects))
	mux.Handle("/projects/{id}", loadAzGHAuthPage(rtProjects.GetRequestStatusByProject))
	mux.Handle("/search/{searchText}", loadAzGHAuthPage(rtSearch.GetSearchResults))
	mux.Handle("/search", loadAzGHAuthPage(rtSearch.SearchHandler))
	mux.Handle("/search/all/", loadAzGHAuthPage(rtSearch.GetAllResults))
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
