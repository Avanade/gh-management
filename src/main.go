package main

import (
	"fmt"
	"log"
	"main/models"
	session "main/pkg/session"
	routes "main/routes"
	"net/http"

	"github.com/codegangsta/negroni"
	"github.com/joho/godotenv"
)

func main() {
	// Set environment variables
	err := godotenv.Load()
	if err != nil {
		log.Print(err.Error())
	}

	// Data on master page
	var menu []models.TypMenu

	menu = append(menu, models.TypMenu{Name: "Home", Url: "/"})
	menu = append(menu, models.TypMenu{Name: "Github", Url: "/github"})

	pageHeaders := models.TypHeaders{Menu: menu}

	// Create session
	session.InitializeSession()

	// Start server
	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir("./public"))
	mux.Handle("/public/", http.StripPrefix("/public/", fs))
	mux.Handle("/", loadPage(routes.IndexHandler, &pageHeaders))
	mux.Handle("/github", loadPage(routes.GithubHandler, &pageHeaders))
	mux.HandleFunc("/login/azure", routes.LoginHandler)
	mux.HandleFunc("/login/azure/callback", routes.CallbackHandler)
	mux.HandleFunc("/logout", routes.LogoutHandler)

	port := "8080"
	fmt.Printf("Now listening on port %v\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), mux))
}

// Verifies authentication before loading the page.
func loadPage(f func(http.ResponseWriter, *http.Request, *models.TypHeaders), pageHeaders *models.TypHeaders) *negroni.Negroni {
	return negroni.New(
		negroni.HandlerFunc(session.IsAuthenticated),
		negroni.Wrap(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			f(w, r, pageHeaders)
		})))
}
