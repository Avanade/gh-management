package routes

import (
	"bytes"
	"context"
	"fmt"
	"main/models"
	email "main/pkg/email"
	ghmgmt "main/pkg/ghmgmtdb"
	"os"
	"text/template"
	"time"
)

// Executes function `f` offsetted by `o`.
func ScheduleJob(ctx context.Context, o time.Duration, f func()) {

	dayPeriod := 24 * time.Hour
	firstRun := time.Now().Truncate(dayPeriod).Add(o)
	n := time.Now()
	if firstRun.Before(n) {
		firstRun = firstRun.Add(dayPeriod)
	}
	firstRunChannel := time.After(firstRun.Sub(time.Now()))
	t := &time.Ticker{C: nil}

	for {
		select {
		case v := <-firstRunChannel:
			t = time.NewTicker(dayPeriod)
			fmt.Printf("A scheduled job was triggered at %s\n", v)
			f()
		case v := <-t.C:
			fmt.Printf("A scheduled job was triggered at %s\n", v)
			f()
		case <-ctx.Done():
			t.Stop()
			return
		}
	}

}

func DailySummaryReport() {
	e := time.Now()
	s := e.AddDate(0, 0, -1)
	o := os.Getenv("GH_ORG_INNERSOURCE")

	r, err := ghmgmt.GetRequestedReposByDateRange(s, e)
	if err != nil {
		fmt.Println("An error occured while pulling the list of projects by date range.")
		return
	}

	if len(r) == 0 {
		fmt.Printf("No repositories were requested on %s.\n No email summary was sent.\n", e)
		return
	}

	var data models.TypRequestedRepoSummary

	data.Date = e.Format("January 02, 2006")
	data.Organization = o
	data.Repos = r

	t, err := template.ParseFiles("templates/reports/requestedRepoSummary.html")
	if err != nil {
		fmt.Println("An error occured while parsing email template.")
		return
	}

	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		fmt.Println("An error occured while composing daily summary of requested repos report.")
		return
	}

	body := buf.String()
	m := email.TypEmailMessage{
		Subject: "Requested Repositories",
		Body:    body,
		To:      "ismael.r.ibuan@accenture.com",
	}

	email.SendEmail(m)
	fmt.Printf("Summary of requested repositories on %s was sent.", e)
}
