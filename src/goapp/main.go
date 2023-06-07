package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	ev "main/pkg/envvar"
	"main/pkg/session"
	rtPages "main/routes/pages"
	reports "main/routes/timerjobs"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/unrolled/secure"
)

func main() {
	// Set environment variables
	err := godotenv.Load()
	if err != nil {
		log.Print(err.Error())
	}

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
		IsDevelopment:         os.Getenv("IS_DEVELOPMENT") == "true",
	})

	// Create session and GitHubClient
	session.InitializeSession()

	// Setup logging format
	log.SetFlags(log.LstdFlags | log.Lshortfile)

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
