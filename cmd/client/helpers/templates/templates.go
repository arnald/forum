package templates

import (
	"html/template"
	"log"
	"net/http"

	"github.com/arnald/forum/internal/pkg/path"
)

// renderTemplate renders a template with the given data.
func RenderTemplate(w http.ResponseWriter, templateName string, data interface{}) {
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
func NotFoundHandler(w http.ResponseWriter, _ *http.Request, errorMessage string, httpStatus int) {
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
