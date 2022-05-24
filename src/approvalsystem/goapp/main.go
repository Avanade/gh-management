package main

import (
	"fmt"
	"log"
	session "main/pkg/session"
	rtAzure "main/routes/login/azure"
	rtPages "main/routes/pages"
	rtApprovals "main/routes/pages/approvals"
	"net/http"
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
	mux.Handle("/response/{appGuid}/{appModuleGuid}/{itemGuid}/{isApproved}", loadAzAuthPage(rtApprovals.ResponseHandler))
	mux.HandleFunc("/login/azure", rtAzure.LoginHandler)
	mux.HandleFunc("/login/azure/callback", rtAzure.CallbackHandler)
	mux.HandleFunc("/logout/azure", rtAzure.LogoutHandler)
	mux.HandleFunc("/request", rtApprovals.ApprovalRequestHandler)
	mux.HandleFunc("/process", rtApprovals.ProcessResponseHandler)
	mux.NotFoundHandler = loadAzAuthPage(rtPages.NotFoundHandler)

	go checkFailedCallbacks()

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

func checkFailedCallbacks() {
	// TIMER SERVICE
	for now := range time.NewTicker(5 * time.Second).C {
		fmt.Println(now, ": Processing failed callbacks.")
		rtApprovals.ProcessFailedCallbacks()
	}
}
