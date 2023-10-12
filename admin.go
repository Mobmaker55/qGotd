package main

import (
	"net/http"
	"strings"
)

func adminHandler(w http.ResponseWriter, r *http.Request) {
	switch strings.Split(r.URL.Path, "/")[2] {
	case "questions":
		adminQuestions(w, r)
		break
	case "reports":
		adminReports(w, r)
		break
	}
}

func adminQuestions(w http.ResponseWriter, r *http.Request) {
	serveTemplate(w, r, "questions.html", map[string]interface{}{})

}

func adminReports(w http.ResponseWriter, r *http.Request) {

}
