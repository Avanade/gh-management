package routes

import (
	"log"
	"net/http"
	"os"

	"main/pkg/template"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("APP RUNNING")
	data := map[string]interface{}{
		"OrganizationName": os.Getenv("ORGANIZATION_NAME"),
	}
	template.UseTemplate(&w, r, "home", data)
}

func ToolHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"toolApprovalProcess": os.Getenv("LINK_TOOL_APPROVAL_PROCESS"),
	}
	template.UseTemplate(&w, r, "tool", data)
}
