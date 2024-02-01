package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"main/pkg/appinsights_wrapper"
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

	secureOptions := secure.Options{
		SSLRedirect:           true,                                            // Strict-Transport-Security
		SSLHost:               os.Getenv("SSL_HOST"),                           // Strict-Transport-Security
		SSLProxyHeaders:       map[string]string{"X-Forwarded-Proto": "https"}, // Strict-Transport-Security
		FrameDeny:             false,
		ContentTypeNosniff:    true, // X-Content-Type-Options
		BrowserXssFilter:      true,
		ReferrerPolicy:        "strict-origin", // Referrer-Policy
		ContentSecurityPolicy: os.Getenv("CONTENT_SECURITY_POLICY"),
		PermissionsPolicy:     "fullscreen=(), geolocation=()", // Permissions-Policy
		STSSeconds:            31536000,                        // Strict-Transport-Security
		STSIncludeSubdomains:  true,                            // Strict-Transport-Security
		IsDevelopment:         os.Getenv("IS_DEVELOPMENT") == "true",
	}

	if os.Getenv("FRAME_EMBEDDOR") != "" {
		secureOptions.CustomFrameOptionsValue = fmt.Sprint("ALLOW-FROM ", os.Getenv("FRAME_EMBEDDOR"))
	}

	secureMiddleware := secure.New(secureOptions)

	// Create session and GitHubClient
	session.InitializeSession()

	// Setup logging format
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Initialize azure application insights
	appinsights_wrapper.Init(os.Getenv("APPINSIGHTS_INSTRUMENTATIONKEY"))

	// SETUP ROUTES
	mux := mux.NewRouter()
	setPageRoutes(mux)
	setAdminPageRoutes(mux)
	setApiRoutes(mux)
	setUtilityRoutes(mux)

	mux.NotFoundHandler = http.HandlerFunc(rtPages.NotFoundHandler)

	o, err := strconv.Atoi(ev.GetEnvVar("SUMMARY_REPORT_TRIGGER", "9"))
	if err != nil {
		fmt.Println(err.Error())
	}
	offset := time.Duration(o) * time.Hour
	ctx := context.Background()
	go reports.ScheduleJob(ctx, offset, reports.DailySummaryReport)
	go checkFailedApprovalRequests()

	mux.Use(
		secureMiddleware.Handler,
		commonHeadersMiddleware,
	)
	http.Handle("/", mux)

	port := ev.GetEnvVar("PORT", "8080")
	fmt.Printf("Now listening on port %v\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), mux))
}

// commonHeadersMiddleware is the middleware function to set common headers
func commonHeadersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set common headers for all requests
		w.Header().Set("Cache-Control", "no-store")

		// Call the next handler in the chain
		next.ServeHTTP(w, r)
	})
}
