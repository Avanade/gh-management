package router

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/unrolled/secure"
)

type muxRouter struct {
	NotFoundHandler http.HandlerFunc
}

func NewMuxRouter() Router {
	return &muxRouter{}
}

var (
	muxDispatcher = mux.NewRouter()
)

func (*muxRouter) GET(uri string, f func(resp http.ResponseWriter, req *http.Request)) {
	muxDispatcher.HandleFunc(uri, f).Methods("GET")
}

func (*muxRouter) POST(uri string, f func(resp http.ResponseWriter, req *http.Request)) {
	muxDispatcher.HandleFunc(uri, f).Methods("POST")
}

func (*muxRouter) PUT(uri string, f func(resp http.ResponseWriter, req *http.Request)) {
	muxDispatcher.HandleFunc(uri, f).Methods("PUT")
}

func (*muxRouter) DELETE(uri string, f func(resp http.ResponseWriter, req *http.Request)) {
	muxDispatcher.HandleFunc(uri, f).Methods("DELETE")
}

func (mr *muxRouter) NOTFOUND(f func(resp http.ResponseWriter, req *http.Request)) {
	muxDispatcher.NotFoundHandler = http.HandlerFunc(f)
}

func (mr *muxRouter) SERVE(port string) {
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

	defaultCors := cors.New(cors.Options{
		AllowedOrigins:   []string{os.Getenv("APPROVAL_SYSTEM_APP_URL")},
		AllowCredentials: true,
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		Debug:            os.Getenv("IS_DEVELOPMENT") == "true",
	}).Handler

	secureMiddleware := secure.New(secureOptions)
	muxDispatcher.Use(
		defaultCors,
		secureMiddleware.Handler,
		commonHeadersMiddleware,
	)
	http.Handle("/", muxDispatcher)

	muxDispatcher.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./public/"))))

	fmt.Printf("Mux HTTP server running on port %v", port)
	http.ListenAndServe(fmt.Sprintf(":%v", port), muxDispatcher)
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
