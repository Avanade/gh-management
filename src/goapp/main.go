package main

import (
	"context"
	"fmt"
	"log"
	router "main/http"
	"main/pkg/appinsights_wrapper"
	ev "main/pkg/envvar"
	"main/pkg/session"
	reports "main/routes/timerjobs"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

var (
	httpRouter router.Router = router.NewMuxRouter()
)

func main() {
	// Set environment variables
	err := godotenv.Load()
	if err != nil {
		log.Print(err.Error())
	}

	// Create session and GitHubClient
	session.InitializeSession()

	// Setup logging format
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Initialize azure application insights
	appinsights_wrapper.Init(os.Getenv("APPINSIGHTS_INSTRUMENTATIONKEY"))

	o, err := strconv.Atoi(ev.GetEnvVar("SUMMARY_REPORT_TRIGGER", "9"))
	if err != nil {
		fmt.Println(err.Error())
	}
	offset := time.Duration(o) * time.Hour
	ctx := context.Background()
	go reports.ScheduleJob(ctx, offset, reports.DailySummaryReport)
	go checkFailedApprovalRequests()

	// SETUP ROUTES
	setPageRoutes()
	setAdminPageRoutes()
	setApiRoutes()
	setUtilityRoutes()

	serve()
}
