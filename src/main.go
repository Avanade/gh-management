package main

import (
	"fmt"
	"html/template"
	"log"
	"main/models"
	session "main/pkg/session"
	routes "main/routes"
	"net/http"

	"github.com/codegangsta/negroni"
	"github.com/joho/godotenv"
)

var tmpl *template.Template

var data models.TypPageData

func main() {
	// Set environment variables
	err := godotenv.Load()
	if err != nil {
		log.Print(err.Error())
	}

	// Create session
	session.InitializeSession()

	// Start server
	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir("./public"))
	mux.Handle("/public/", http.StripPrefix("/public/", fs))
	mux.Handle("/", loadPage(routes.IndexHandler))
	mux.Handle("/github", loadPage(routes.GithubHandler))
	mux.HandleFunc("/login", routes.LoginHandler)
	mux.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		routes.CallbackHandler(w, r, &data)
	})
	mux.HandleFunc("/logout", routes.LogoutHandler)

	port := "8080"
	fmt.Printf("Now listening on port %v\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), mux))
}

// Verifies authentication before loading the page.
func loadPage(f func(w http.ResponseWriter, r *http.Request, data *models.TypPageData)) *negroni.Negroni {
	return negroni.New(
		negroni.HandlerFunc(session.IsAuthenticated),
		negroni.Wrap(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// reset contents on data
			data.Content = nil

			// data pointer is passed on
			f(w, r, &data)
		})))
}
