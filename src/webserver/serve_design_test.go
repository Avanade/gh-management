package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWebserverLoads(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	res := httptest.NewRecorder()

	server := &templateHandler{filename: "index.html"}
	server.ServeHTTP(res, req)

	got := res.Code
	want := http.StatusOK

	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

/* 	err := godotenv.Load()
log.Printf("Loaded as %v\n", os.Getenv("ENV"))
if err != nil && os.Getenv("ENV") == "local" {
	log.Fatal("Error loading .env file")
}

*/
