package main

import (
	"context"
	"fmt"
	"log"
	session "main/pkg/session"
	"os"

	//rtGithubAPi "main/routes/login/github"
	rtPages "main/routes/pages"
	reports "main/routes/timerjobs"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/unrolled/secure"

	ev "main/pkg/envvar"

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
		IsDevelopment:         false,
	})

	// Set environment variables
	err := godotenv.Load()
	if err != nil {
		log.Print(err.Error())
	}

	// Create session and GitHubClient
	session.InitializeSession()

	// SETUP ROUTES
	mux := mux.NewRouter()
	setPageRoutes(mux)
	setAdminPageRoutes(mux)
	setApiRoutes(mux)

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
