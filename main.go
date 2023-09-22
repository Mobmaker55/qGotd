package main

import (
	"errors"
	"fmt"
	"html/template"
	"io/fs"
	"maps"
	"net/http"
	"path/filepath"
	"slices"
)

var Mux *http.ServeMux
var HtmlTemplates *template.Template

func main() {
	readConfig()
	Mux = http.NewServeMux()
	initAuth()
	initTemplates()

	//register main handlers
	Mux.HandleFunc("/", HttpLogger(getIndex))

	//start the server
	err := http.ListenAndServe(Config.Address, Mux)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server one closed\n")
	} else if err.Error() == "" {
		fmt.Println("HTTP server started")
	} else {
		fmt.Println(err)
	}

}

func getIndex(w http.ResponseWriter, r *http.Request) {
	if !IsAuthenticated(w, r) {
		http.Redirect(w, r, "http://"+r.Host+"/login", http.StatusFound)
		return
	}
	serveTemplate(w, r, "index.html", map[string]interface{}{})
}

// serves template and handles error
func serveTemplate(w http.ResponseWriter, r *http.Request, templateName string, data map[string]interface{}) {
	user, err := GetUser(r)
	isEvil := false
	if slices.Contains(user.Groups, "eboard") || slices.Contains(user.Groups, "active_rtp") || user.PreferredUsername == "mob" {
		isEvil = true
	}

	//internally parsed information
	vars := map[string]interface{}{
		"name":     user.GivenName,
		"username": user.PreferredUsername,
		"isEvil":   isEvil,
		"template": "index.html",
	}

	maps.Copy(vars, data)
	err = HtmlTemplates.ExecuteTemplate(w, templateName, vars)
	if err != nil {
		fmt.Println("template failed to execute, ", err.Error())
		return
	}
}

// initializes templates from the ./html/ directory
func initTemplates() {
	//path, _ := os.Getwd()
	var templates []string
	err := filepath.Walk("./html/", func(path string, d fs.FileInfo, err error) error {
		if err != nil {
			fmt.Println(err)
			return err
		}
		if path == "./html/" {
			return nil
		}
		templates = append(templates, path)
		return nil
	})
	if err != nil {
		return
	}
	HtmlTemplates, err = template.ParseFiles(templates...)
}
