package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	session "webserver/pkg/session"
	routes "webserver/routes"

	"github.com/codegangsta/negroni"
	"github.com/joho/godotenv"
)

var tmpl *template.Template

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
	mux.Handle("/", negroni.New(
		negroni.HandlerFunc(session.IsAuthenticated),
		negroni.Wrap(http.HandlerFunc(routes.IndexHandler)),
	))
	mux.HandleFunc("/login", routes.LoginHandler)
	mux.HandleFunc("/callback", routes.CallbackHandler)
	mux.HandleFunc("/logout", routes.LogoutHandler)

	port := "8080"
	fmt.Printf("Now listening on port %v\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), mux))
}
