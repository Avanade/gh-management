package routes

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"text/template"
	"time"

	"main/pkg/email"
	db "main/pkg/ghmgmtdb"
)

type RequestedRepositorySummary struct {
	Date         string
	Organization string
	Repos        []db.Repository
}

// Executes function `f` offsetted by `o`.
func ScheduleJob(ctx context.Context, o time.Duration, f func()) {

	dayPeriod := 24 * time.Hour
	firstRun := time.Now().Truncate(dayPeriod).Add(o)
	n := time.Now()
	if firstRun.Before(n) {
		firstRun = firstRun.Add(dayPeriod)
	}
	firstRunChannel := time.After(time.Until(firstRun))
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
	recipient := os.Getenv("EMAIL_SUMMARY_REPORT")

	r, err := db.GetRequestedReposByDateRange(s, e)
	if err != nil {
		fmt.Println("An error occured while pulling the list of projects by date range.")
		return
	}

	if len(r) == 0 {
		fmt.Printf("No repositories were requested on %s.\n No email summary was sent.\n", e)
		return
	}

	var data RequestedRepositorySummary

	data.Date = e.Format("January 02, 2006")
	data.Organization = o
	data.Repos = r

	t, err := template.ParseFiles("templates/reports/requestedreposummary.html")
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
	m := email.Message{
		Subject: "Requested Repositories",
		Body: email.Body{
			Content: body,
			Type:    email.HtmlMessageType,
		},
		ToRecipients: []email.Recipient{
			{
				Email: recipient,
			},
		},
	}

	email.SendEmail(m)
	fmt.Printf("Summary of requested repositories on %s was sent.", e)
}
