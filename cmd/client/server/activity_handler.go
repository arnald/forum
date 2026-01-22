package server

import (
	"context"
	"log"
	"net/http"
	"text/template"

	"github.com/arnald/forum/cmd/client/domain"
	"github.com/arnald/forum/cmd/client/helpers"
	"github.com/arnald/forum/cmd/client/helpers/templates"
	"github.com/arnald/forum/cmd/client/middleware"
)

// ActivityPage handles requests to the activity page.
func (cs *ClientServer) ActivityPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), requestTimeout)
	defer cancel()

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, backedGetUserActivity(), nil)
	if err != nil {
		http.Error(w, "Error creating request", http.StatusInternalServerError)
		return
	}

	ip := middleware.GetIPFromContext(r)
	if ip == "" {
		http.Error(w, "Error no IP found in request", http.StatusInternalServerError)
	}

	helpers.SetIPHeaders(httpReq, ip)

	httpReq.Header.Set("Content-Type", "application/json")

	for _, cookie := range r.Cookies() {
		log.Printf("Forwarding cookie: %s=%s", cookie.Name, cookie.Value)
		httpReq.AddCookie(cookie)
	}

	backendResp, err := cs.HTTPClient.Do(httpReq)
	if err != nil {
		log.Printf("Error making request to backend: %v", err)
		http.Error(w, "Error communicating with backend", http.StatusInternalServerError)
		return
	}
	defer backendResp.Body.Close()

	var activityData domain.ActivityData
	err = helpers.DecodeBackendResponse(backendResp, &activityData)
	if err != nil {
		http.Error(w, "Error decoding the response to json", http.StatusInternalServerError)
		return
	}

	user := middleware.GetUserFromContext(r.Context())

	activityData.User = user

	tmpl, err := template.ParseFiles(
		"frontend/html/layouts/base.html",
		"frontend/html/pages/activity.html",
		"frontend/html/partials/navbar.html",
		"frontend/html/partials/footer.html",
	)
	if err != nil {
		templates.NotFoundHandler(w, r, "Failed to load page", http.StatusInternalServerError)
		return
	}

	err = tmpl.ExecuteTemplate(w, "base", activityData)
	if err != nil {
		log.Println("Error executing template:", err)
		http.Error(w, "Failed to render page", http.StatusInternalServerError)
	}
}
