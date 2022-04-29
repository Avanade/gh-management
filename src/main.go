package main

import (
	"fmt"
	"log"
	session "main/pkg/session"
	routes "main/routes"
	"net/http"

	"github.com/codegangsta/negroni"
	"github.com/joho/godotenv"

	ev "main/pkg/envvar"
)

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
	mux.HandleFunc("/login/azure", routes.LoginHandler)
	mux.HandleFunc("/login/azure/callback", routes.CallbackHandler)
	mux.HandleFunc("/logout", routes.LogoutHandler)
	mux.HandleFunc("/login/github", routes.GithubLoginHandler)
	mux.HandleFunc("/login/github/callback", routes.GithubCallbackHandler)

	port := ev.GetEnvVar("PORT", "8080")
	fmt.Printf("Now listening on port %v\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), mux))

}

// Verifies authentication before loading the page.
func loadPage(f func(w http.ResponseWriter, r *http.Request)) *negroni.Negroni {
	return negroni.New(
		negroni.HandlerFunc(session.IsAuthenticated),
		negroni.Wrap(http.HandlerFunc(f)),
	)
}
