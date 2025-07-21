package handler

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/arnald/forum/internal/pkg/path"
)

const (
	notFoundMessage = "Oops! The page you're looking for has vanished into the digital void."
)

type CategoryData struct {
	Data struct {
		Categories []Category `json:"categories"`
	} `json:"data"`
}

type Category struct {
	Name        string  `json:"name"`
	Color       string  `json:"color"`
	Slug        string  `json:"slug,omitzero"`
	Description string  `json:"description,omitzero"`
	Topics      []Topic `json:"topics,omitzero"`
	Logo        Logo    `json:"logo"`
	ID          int     `json:"id"`
}

type Logo struct {
	URL    string `json:"url"`
	ID     int    `json:"id"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

type Topic struct {
	Title string `json:"title"`
	ID    int    `json:"id"`
}

func HomePage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
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

	var categoryData CategoryData
	err = json.NewDecoder(file).Decode(&categoryData)
	if err != nil {
		log.Println("Error decoding JSON:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseGlob(resolver.GetPath("frontend/html/**/*.html"))
	if err != nil {
		log.Println("Error loading home.html:", err)
		notFoundHandler(w, r, "Failed to load page", http.StatusInternalServerError)

		return
	}

	err = tmpl.ExecuteTemplate(w, "base", categoryData.Data.Categories)
	if err != nil {
		log.Println("Error executing template:", err)
		http.Error(w, "Failed to render page", http.StatusInternalServerError)
	}
}

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
