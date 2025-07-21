package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"text/template"
)

const (
	notFoundMessage = "Oops! The page you're looking for has vanished into the digital void."
)

type Topic struct {
	Title string `json:"title"`
	ID    int    `json:"id"`
}

type Logo struct {
	URL    string `json:"url"`
	ID     int    `json:"id"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

type Category struct {
	Name        string  `json:"name"`
	Color       string  `json:"color"`
	Slug        string  `json:"slug"`
	Description string  `json:"description"`
	Topics      []Topic `json:"topics"`
	Logo        Logo    `json:"logo"`
	ID          int     `json:"id"`
}

type CategoryData struct {
	Data struct {
		Categories []Category `json:"categories"`
	} `json:"data"`
}

// HomePage Handler
func HomePage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		notFoundHandler(w, r, notFoundMessage, http.StatusNotFound)

		return
	}

	basePath, err := os.Getwd()
	if err != nil {
		log.Println("Error getting working directory:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Load categories.json
	jsonPath := filepath.Join(basePath, "cmd", "client", "data", "categories.json")
	jsonPath = filepath.Clean(jsonPath)

	file, err := os.Open(jsonPath)
	if err != nil {
		log.Println("Error opening categories.json:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	var data CategoryData
	err = json.NewDecoder(file).Decode(&data)
	if err != nil {
		log.Println("Error decoding JSON:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// templatePath := filepath.Join(basePath, "frontend", "html", "pages", "home.html")
	// tmpl, err := template.ParseFiles(templatePath)
	tmpl, err := template.ParseGlob(filepath.Join(basePath, "frontend/html/**/*.html"))
	if err != nil {
		log.Println("Error loading home.html:", err)
		notFoundHandler(w, r, "Failed to load page", http.StatusInternalServerError)

		return
	}
	err = tmpl.ExecuteTemplate(w, "base", data.Data.Categories)
	if err != nil {
		log.Println("Error executing template:", err)
		http.Error(w, "Failed to render page", http.StatusInternalServerError)
	}
}

// Register Handler
func RegisterPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	basePath, err := os.Getwd()
	if err != nil {
		log.Println("Error getting working directory:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	tmplPath := filepath.Join(basePath, "frontend", "html", "pages", "register.html")
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		log.Println("Error parsing register.html:", err)
		http.Error(w, "Failed to load page", http.StatusInternalServerError)
		return
	}

	err = tmpl.ExecuteTemplate(w, "register", nil)
	if err != nil {
		log.Println("Error executing register.html:", err)
		http.Error(w, "Failed to render page", http.StatusInternalServerError)
	}
}

// Login Handler
func LoginPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	basePath, err := os.Getwd()
	if err != nil {
		log.Println("Error getting working directory:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	tmplPath := filepath.Join(basePath, "frontend", "html", "pages", "login.html")
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		log.Println("Error parsing login.html:", err)
		http.Error(w, "Failed to load page", http.StatusInternalServerError)
		return
	}

	err = tmpl.ExecuteTemplate(w, "login", nil)
	if err != nil {
		log.Println("Error executing login.html:", err)
		http.Error(w, "Failed to render page", http.StatusInternalServerError)
	}
}

func notFoundHandler(w http.ResponseWriter, _ *http.Request, errorMessage string, httpStatus int) {
	basePath, err := os.Getwd()
	if err != nil {
		log.Println("Error getting working directory:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	templatePath := filepath.Join(basePath, "frontend", "html", "pages", "not_found.html")
	templatePath = filepath.Clean(templatePath)

	tmpl, err := template.ParseFiles(templatePath)
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
