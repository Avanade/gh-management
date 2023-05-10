package main

import (
	"context"
	"fmt"
	"log"
	session "main/pkg/session"
	rtApi "main/routes/api"
	rtAzure "main/routes/login/azure"
	rtGithub "main/routes/login/github"
	"os"

	//rtGithubAPi "main/routes/login/github"
	rtPages "main/routes/pages"
	rtActivities "main/routes/pages/activities"
	rtAdmin "main/routes/pages/admin"
	rtCommunity "main/routes/pages/community"
	rtGuidance "main/routes/pages/guidance"
	rtProjects "main/routes/pages/projects"
	rtSearch "main/routes/pages/search"
	reports "main/routes/timerjobs"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/unrolled/secure"

	ev "main/pkg/envvar"

	"github.com/codegangsta/negroni"
	"github.com/joho/godotenv"
)

func main() {
	secureMiddleware := secure.New(secure.Options{
		SSLRedirect:           true,                                            // Strict-Transport-Security
		SSLHost:               os.Getenv("SSL_HOST"),                           // Strict-Transport-Security
		SSLProxyHeaders:       map[string]string{"X-Forwarded-Proto": "https"}, // Strict-Transport-Security
		FrameDeny:             true,                                            // X-FRAME-OPTIONS
		ContentTypeNosniff:    true,                                            // X-Content-Type-Options
		BrowserXssFilter:      true,
		ReferrerPolicy:        "strict-origin", // Referrer-Policy
		ContentSecurityPolicy: os.Getenv("CONTENT_SECURITY_POLICY"),
		PermissionsPolicy:     "fullscreen=(), geolocation=()", // Permissions-Policy
		STSSeconds:            31536000,                        // Strict-Transport-Security
		STSIncludeSubdomains:  true,                            // Strict-Transport-Security
		IsDevelopment:         true,
	})

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
	mux.Handle("/activities", loadAzGHAuthPage(rtActivities.ActivitiesHandler))
	mux.Handle("/activities/{action:add}", loadAzGHAuthPage(rtActivities.ActivitiesNewHandler))
	mux.Handle("/activities/{action:edit|view}/{id}", loadAzGHAuthPage(rtActivities.ActivitiesNewHandler))
	mux.Handle("/repositories", loadAzGHAuthPage(rtProjects.Projects))
	mux.Handle("/repositories/new", loadAzGHAuthPage(rtProjects.ProjectsNewHandler))
	mux.Handle("/repositories/{id}", loadAzGHAuthPage(rtProjects.ProjectsHandler))
	mux.Handle("/repositories/makepublic/{id}", loadAzGHAuthPage(rtProjects.MakePublic))
	mux.Handle("/search/{offSet}/{rowCount}", loadAzGHAuthPage(rtSearch.GetSearchResults))
	mux.Handle("/search", loadAzGHAuthPage(rtSearch.SearchHandler))

	mux.Handle("/guidance", loadAzGHAuthPage(rtGuidance.GuidanceHandler))
	mux.Handle("/guidance/new", loadAzGHAuthPage(rtGuidance.CategoriesHandler))
	mux.Handle("/guidance/{id}", loadAzGHAuthPage(rtGuidance.CategoryUpdateHandler))
	mux.Handle("/guidance/Article/{id}", loadAzGHAuthPage(rtGuidance.ArticleHandler))
	mux.Handle("/community/new", loadAzGHAuthPage(rtCommunity.CommunityHandler))
	mux.Handle("/community/my", loadAzGHAuthPage(rtCommunity.GetMyCommunitylist))
	mux.Handle("/community/{id}", loadAzGHAuthPage(rtCommunity.CommunityHandler))
	mux.Handle("/community/getcommunity/{id}", loadAzGHAuthPage(rtCommunity.GetUserCommunity))

	mux.Handle("/communities/list", loadAzGHAuthPage(rtCommunity.CommunitylistHandler))
	mux.Handle("/community", loadAzGHAuthPage(rtCommunity.GetUserCommunitylist))
	mux.Handle("/community/{id}/onboarding", loadAzGHAuthPage(rtCommunity.CommunityOnBoarding))
	mux.HandleFunc("/loginredirect", rtPages.LoginRedirectHandler).Methods("GET")
	mux.HandleFunc("/gitredirect", rtPages.GitredirectHandler).Methods("GET")
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
	muxApi.Handle("/activity", loadAzGHAuthPage(rtApi.CreateActivity)).Methods("POST")
	muxApi.Handle("/activity", loadAzGHAuthPage(rtApi.GetActivities)).Methods("GET")
	muxApi.Handle("/activity/{id}", loadAzGHAuthPage(rtApi.GetActivityById)).Methods("GET")
	muxApi.Handle("/community", loadAzGHAuthPage(rtApi.CommunityAPIHandler))
	muxApi.Handle("/communitySponsors", loadAzGHAuthPage(rtApi.CommunitySponsorsAPIHandler))
	muxApi.Handle("/CommunitySponsorsPerCommunityId/{id}", loadAzGHAuthPage(rtApi.CommunitySponsorsPerCommunityId))
	muxApi.Handle("/CommunityTagPerCommunityId/{id}", loadAzGHAuthPage(rtApi.CommunityTagPerCommunityId))
	muxApi.Handle("/community/onboarding/{id}", loadAzGHAuthPage(rtApi.GetCommunityOnBoardingInfo)).Methods("GET", "POST", "DELETE")
	muxApi.Handle("/community/all", loadAzAuthPage(rtApi.GetCommunities)).Methods("GET")
	muxApi.Handle("/community/{id}/members", loadAzAuthPage(rtApi.GetCommunityMembers)).Methods("GET")
	muxApi.Handle("/communitystatus/{id}", loadAzGHAuthPage(rtApi.GetRequestStatusByCommunity))
	muxApi.Handle("/community/getCommunitiesisexternal/{isexternal}", loadAzGHAuthPage(rtApi.GetCommunitiesIsexternal))
	muxApi.Handle("/contributionarea", loadAzGHAuthPage(rtApi.CreateContributionAreas)).Methods("POST")
	muxApi.Handle("/contributionarea", loadAzGHAuthPage(rtApi.GetContributionAreas)).Methods("GET")
	muxApi.Handle("/contributionarea", loadAzGHAuthPage(rtApi.UpdateContributionArea)).Methods("PUT")
	muxApi.Handle("/contributionarea/{id}", loadAzGHAuthPage(rtApi.GetContributionAreaById)).Methods("GET")
	muxApi.Handle("/contributionarea/activity/{id}", loadAzGHAuthPage(rtApi.GetContributionAreasByActivityId)).Methods("GET")
	muxApi.Handle("/Category", loadAzGHAuthPage(rtApi.CategoryAPIHandler))
	muxApi.Handle("/Category/list", loadAzGHAuthPage(rtApi.CategoryListAPIHandler))
	muxApi.Handle("/Category/update", loadAzGHAuthPage(rtApi.CategoryUpdate))
	muxApi.Handle("/Category/{id}", loadAzGHAuthPage(rtApi.GetCategoryByID))
	muxApi.Handle("/importGitHubReposToDatabase", loadAzAuthPage(rtApi.ImportReposToDatabase))

	muxApi.Handle("/relatedcommunityAdd", loadAzAuthPage(rtApi.RelatedCommunitiesInsert))
	muxApi.Handle("/relatedcommunityDelete", loadAzAuthPage(rtApi.RelatedCommunitiesDelete))
	muxApi.Handle("/relatedcommunity/{id}", loadAzAuthPage(rtApi.RelatedCommunitiesSelect)).Methods("GET")

	muxApi.Handle("/CategoryArticlesById/{id}", loadAzGHAuthPage(rtApi.GetCategoryArticlesById))
	muxApi.Handle("/CategoryArticlesByArticlesID/{id}", loadAzGHAuthPage(rtApi.GetCategoryArticlesByArticlesID))
	muxApi.Handle("/CategoryArticlesUpdate", loadAzGHAuthPage(rtApi.CategoryArticlesUpdate))
	muxApi.Handle("/repositories/list", loadAzGHAuthPage(rtApi.GetUserProjects))
	muxApi.Handle("/repositories/{id}", loadAzGHAuthPage(rtApi.GetRequestStatusByProject))
	muxApi.Handle("/repositories/request/public", loadAzGHAuthPage(rtApi.RequestMakePublic))
	muxApi.Handle("/repositories/collaborators/{id}", loadAzGHAuthPage(rtApi.GetRepoCollaboratorsByRepoId))
	muxApi.Handle("/repositories/collaborators/add/{id}/{ghUser}/{permission}", loadAzGHAuthPage(rtApi.AddCollaborator))
	muxApi.Handle("/repositories/collaborators/remove/{id}/{ghUser}/{permission}", loadAzGHAuthPage(rtApi.RemoveCollaborator))
	muxApi.Handle("/repositories/archive/{project}/{projectId}/{state}/{archive}", loadAzGHAuthPage(rtApi.ArchiveProject))
	muxApi.Handle("/repositories/visibility/{project}/{projectId}/{currentState}/{desiredState}", loadAzGHAuthPage(rtApi.SetVisibility))
	muxApi.Handle("/allusers", loadAzAuthPage(rtApi.GetAllUserFromActiveDirectory))
	muxApi.Handle("/alluserswithgithub", loadAzAuthPage(rtApi.GetUsersWithGithub))
	muxApi.Handle("/search/users/{search}", loadAzAuthPage(rtApi.SearchUserFromActiveDirectory))
	muxApi.Handle("/allrepositories", loadAzAuthPage(rtApi.GetAllRepositories))
	muxApi.Handle("/getActiveApprovalTypes", loadAzGHAuthPage(rtApi.GetActiveApprovalTypes))

	//API FOR APPROVAL TYPES
	muxApi.HandleFunc("/approval/type", rtApi.CreateApprovalType).Methods("POST")
	muxApi.HandleFunc("/approval/type/{id}", rtApi.EditApprovalTypeById).Methods("PUT")
	muxApi.HandleFunc("/approval/type/{id}/archived", rtApi.SetIsArchivedApprovalTypeById).Methods("PUT")
	muxApi.HandleFunc("/approval/types", rtApi.GetApprovalTypes).Methods("GET")
	muxApi.HandleFunc("/approval/type/{id}", rtApi.GetApprovalTypeById).Methods("GET")

	// API FOR LOGIC APP
	muxApi.Handle("/init/indexorgrepos", loadGuidAuthApi(rtApi.InitIndexOrgRepos)).Methods("GET")
	muxApi.Handle("/indexorgrepos", loadGuidAuthApi(rtApi.IndexOrgRepos)).Methods("GET")
	muxApi.Handle("/checkAvaInnerSource", loadGuidAuthApi(rtApi.CheckAvaInnerSource)).Methods("GET")
	muxApi.Handle("/checkAvaOpenSource", loadGuidAuthApi(rtApi.CheckAvaOpenSource)).Methods("GET")
	muxApi.Handle("/clearOrgMembers", loadGuidAuthApi(rtApi.ClearOrgMembers)).Methods("GET")
	muxApi.Handle("/RepoOwnerScan", loadGuidAuthApi(rtApi.RepoOwnerScan)).Methods("GET")

	// API FOR ProjectToRepoOwner APP
	muxApi.Handle("/projectToRepoOwner", loadGuidAuthApi(rtApi.ProjectToRepoOwner)).Methods("GET")
	muxApi.Handle("/RepoOwnersCleanup", loadGuidAuthApi(rtApi.RepoOwnersCleanup)).Methods("GET")
	muxAdmin := mux.PathPrefix("/admin").Subrouter()
	muxAdmin.Handle("", loadAdminPage(rtAdmin.AdminIndex))
	muxAdmin.Handle("/members", loadAdminPage(rtAdmin.ListCommunityMembers))
	muxAdmin.Handle("/guidance", loadAdminPage(rtGuidance.GuidanceHandler))
	muxAdmin.Handle("/approvaltypes", loadAdminPage(rtAdmin.ListApprovalTypes))
	muxAdmin.Handle("/communityapprovers", loadAdminPage(rtCommunity.CommunityApproverHandler))
	muxAdmin.Handle("/approvaltype/{action:add}", loadAdminPage(rtAdmin.ApprovalTypeForm))
	muxAdmin.Handle("/approvaltype/{action:view|edit|delete}/{id}", loadAdminPage(rtAdmin.ApprovalTypeForm))
	muxAdmin.Handle("/contributionarea", loadAdminPage(rtAdmin.ListContributionAreas))
	muxAdmin.Handle("/contributionarea/{action:add}", loadAdminPage(rtAdmin.ContributionAreasForm))
	muxAdmin.Handle("/contributionarea/{action:view|edit}/{id}", loadAdminPage(rtAdmin.ContributionAreasForm))

	//EXTERNAL LINKS
	muxAdmin.Handle("/externallinks", loadAdminPage(rtAdmin.ExternalLinksHandler))
	muxAdmin.Handle("/externallinks/", loadAdminPage(rtAdmin.GetExternalLinks))
	muxAdmin.Handle("/externallinks/enabled", loadAdminPage(rtAdmin.GetExternalLinksAllEnabled))
	muxAdmin.Handle("/externallinks/{id}", loadAdminPage(rtAdmin.GetExternalLinksById))
	muxAdmin.Handle("/externallinks/{action:add}/", loadAdminPage(rtAdmin.ExternalLinksForm))
	muxAdmin.Handle("/externallinks/{action:view|edit}/{id}", loadAdminPage(rtAdmin.ExternalLinksForm))
	muxAdmin.Handle("/externallinks/{action:delete}/{id}", loadAdminPage(rtAdmin.ExternalLinksDelete))
	//EXTERNAL LINKS API
	muxApi.HandleFunc("/externallinks/create", rtAdmin.CreateExternalLinks).Methods("POST")
	muxApi.HandleFunc("/externallinks/update/{id}", rtAdmin.UpdateExternalLinks).Methods("PUT")

	muxApi.HandleFunc("/approvals/project/callback", rtProjects.UpdateApprovalStatusProjects).Methods("POST")
	muxApi.HandleFunc("/approvals/project/reassign/callback", rtProjects.UpdateApprovalReassignApprover)
	muxApi.HandleFunc("/approvals/community/callback", rtProjects.UpdateApprovalStatusCommunity).Methods("POST")
	muxApi.HandleFunc("/approvals/community/callback", rtProjects.UpdateApprovalStatusCommunity).Methods("POST")
	muxApi.HandleFunc("/communityapprovers/update", rtCommunity.CommunityApproversListUpdate)
	muxApi.HandleFunc("/communityapprovers/GetCommunityApproversList", rtCommunity.GetCommunityApproversList)
	muxApi.HandleFunc("/communityapprovers/GetAllActiveCommunityApprovers", rtCommunity.GetAllActiveCommunityApprovers)
	muxApi.HandleFunc("/communityapprovers/GetCommunityApproversList/{id}", rtCommunity.GetCommunityApproversById)
	mux.NotFoundHandler = http.HandlerFunc(rtPages.NotFoundHandler)

	o, err := strconv.Atoi(ev.GetEnvVar("SUMMARY_REPORT_TRIGGER", "9"))
	if err != nil {
		fmt.Println(err.Error())
	}
	offset := time.Duration(o) * time.Hour
	ctx := context.Background()
	go reports.ScheduleJob(ctx, offset, reports.DailySummaryReport)
	go checkFailedApprovalRequests()

	mux.Use(secureMiddleware.Handler)
	http.Handle("/", mux)

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

func loadGuidAuthApi(f func(w http.ResponseWriter, r *http.Request)) *negroni.Negroni {
	return negroni.New(
		negroni.HandlerFunc(session.IsGuidAuthenticated),
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

func loadAdminPage(f func(w http.ResponseWriter, r *http.Request)) *negroni.Negroni {
	return negroni.New(
		negroni.HandlerFunc(session.IsAuthenticated),
		negroni.HandlerFunc(session.IsUserAdminMW),
		negroni.Wrap(http.HandlerFunc(f)),
	)
}

func checkFailedApprovalRequests() {
	// TIMER SERVICE
	freq := ev.GetEnvVar("APPROVALREQUESTS_RETRY_FREQ", "15")
	freqInt, _ := strconv.ParseInt(freq, 0, 64)
	if freq > "0" {
		for range time.NewTicker(time.Duration(freqInt) * time.Minute).C {
			go rtApi.ReprocessRequestApproval()
			go rtCommunity.ReprocessRequestCommunityApproval()
		}
	}
}
