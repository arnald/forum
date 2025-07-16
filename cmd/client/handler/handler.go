package handler

import (
	"log"
	"net/http"
	"strconv"
	"text/template"
)

const (
	notFoundMessage = "Oops! The page you're looking for has vanished into the digital void."
)

func HomePage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		notFoundHandler(w, r, notFoundMessage, http.StatusNotFound)

		return
	}
	tmpl, err := template.ParseFiles("frontend/html/pages/home.html")
	if err != nil {
		log.Println("Error loading home.html:", err)
		notFoundHandler(w, r, "Failed to load page", http.StatusInternalServerError)

		return
	}
	tmpl.Execute(w, nil)
}

func notFoundHandler(w http.ResponseWriter, _ *http.Request, errorMessage string, httpStatus int) {
	tmpl, err := template.ParseFiles("frontend/html/pages/not_found.html")
	if err != nil {
		http.Error(w, errorMessage, httpStatus)
		log.Println("Error loading not_found_page.html:", err)

		return
	}

	w.WriteHeader(httpStatus)
	tmpl.Execute(w, map[string]string{
		"ErrorMessage":   errorMessage,
		"HttpStatusCode": strconv.Itoa(httpStatus),
	})
}
