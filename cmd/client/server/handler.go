package server

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/arnald/forum/cmd/client/domain"
	h "github.com/arnald/forum/cmd/client/helpers"
	"github.com/arnald/forum/cmd/client/middleware"
	"github.com/arnald/forum/internal/pkg/path"
)

// ============================================================================
// HELPER FUNCTIONS (Package-level, not methods)
// ============================================================================

// renderTemplate renders a template with the given data.
func renderTemplate(w http.ResponseWriter, templateName string, data interface{}) {
	resolver := path.NewResolver()
	tmplPath := resolver.GetPath("frontend/html/pages/" + templateName + ".html")

	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		log.Printf("Error parsing %s: %v", tmplPath, err)
		http.Error(w, "Failed to load page", http.StatusInternalServerError)
		return
	}

	err = tmpl.ExecuteTemplate(w, templateName, data)
	if err != nil {
		log.Printf("Error executing %s template: %v", templateName, err)
		http.Error(w, "Failed to render page", http.StatusInternalServerError)
	}
}

// notFoundHandler renders a 404 error page.
func notFoundHandler(w http.ResponseWriter, _ *http.Request, errorMessage string, httpStatus int) {
	resolver := path.NewResolver()

	tmpl, err := template.ParseFiles(resolver.GetPath("frontend/html/pages/not_found.html"))
	if err != nil {
		http.Error(w, errorMessage, httpStatus)
		log.Println("Error loading not_found_page.html:", err)
		return
	}

	data := struct {
		StatusText   string
		ErrorMessage string
		StatusCode   int
	}{
		StatusText:   http.StatusText(httpStatus),
		ErrorMessage: errorMessage,
		StatusCode:   httpStatus,
	}

	w.WriteHeader(httpStatus)
	err = tmpl.Execute(w, data)
	if err != nil {
		log.Println("Error executing template:", err)
		http.Error(w, errorMessage, httpStatus)
	}
}

// backendError is a custom error type for backend errors.
type backendError string

func (e backendError) Error() string {
	return string(e)
}

// ============================================================================
// PAGE HANDLERS (Methods on ClientServer)
// ============================================================================

// HomePage handles requests to the homepage.
func (cs *ClientServer) HomePage(w http.ResponseWriter, r *http.Request) {
	// Handle / and /categories
	if r.URL.Path != "/" && r.URL.Path != "/categories" {
		notFoundHandler(w, r, notFoundMessage, http.StatusNotFound)
		return
	}

	resolver := path.NewResolver()

	file, err := os.Open(resolver.GetPath("cmd/client/data/categories.json"))
	if err != nil {
		log.Println("Error opening categories.json:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	var categoryData domain.CategoryData
	err = json.NewDecoder(file).Decode(&categoryData)
	if err != nil {
		log.Println("Error decoding JSON:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	h.PrepareCategories(categoryData.Data.Categories)

	// Get user from context (set by middleware)
	user := middleware.GetUserFromContext(r.Context())

	pageData := domain.HomePageData{
		Categories: categoryData.Data.Categories,
		ActivePage: r.URL.Path,
		User:       user,
	}

	tmpl, err := template.ParseFiles(
		"frontend/html/layouts/base.html",
		"frontend/html/pages/home.html",
		"frontend/html/partials/navbar.html",
		"frontend/html/partials/category_details.html",
		"frontend/html/partials/categories.html",
		"frontend/html/partials/footer.html",
	)
	if err != nil {
		log.Println("Error loading home.html:", err)
		notFoundHandler(w, r, "Failed to load page", http.StatusInternalServerError)
		return
	}

	err = tmpl.ExecuteTemplate(w, "base", pageData)
	if err != nil {
		log.Println("Error executing template:", err)
		http.Error(w, "Failed to render page", http.StatusInternalServerError)
	}
}
