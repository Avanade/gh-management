package routes

import (
	"log"
	"net/http"
	"os"

	"main/pkg/template"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("APP RUNNING")
	template.UseTemplate(&w, r, "home", nil)
}

func ToolHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"toolApprovalProcess": os.Getenv("LINK_TOOL_APPROVAL_PROCESS"),
	}
	template.UseTemplate(&w, r, "tool", data)
}
