package main

import (
	"errors"
	"fmt"
	"html/template"
	"io/fs"
	"maps"
	"net/http"
	"path/filepath"
)

var Mux *http.ServeMux
var HtmlTemplates *template.Template

func main() {
	readConfig()
	Mux = http.NewServeMux()
	initAuth()
	initTemplates()

	//register main handlers
	Mux.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("js"))))
	Mux.HandleFunc("/response/", Handler(responseHandler))
	Mux.HandleFunc("/", Handler(getIndex))

	//start the server
	fmt.Println("Pre-init complete, starting server")
	err := http.ListenAndServe(Config.Address, Mux)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server one closed\n")
	} else if err.Error() == "" {
		fmt.Println("HTTP server started")
	} else {
		fmt.Println(err)
	}

}

func testSite(w http.ResponseWriter, r *http.Request) {
	fmt.Println("tested")
	w.Write([]byte("Hi"))
}

func getIndex(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"question": "Do you like this app?",
	}

	serveTemplate(w, r, "index.html", data)
}

// serves template and handles error
func serveTemplate(w http.ResponseWriter, r *http.Request, templateName string, data map[string]interface{}) {
	user := GetUser(w, r)

	//internally parsed information
	vars := map[string]interface{}{
		"name":     user.GivenName,
		"fullName": user.GivenName + " " + user.FamilyName,
		"username": user.PreferredUsername,
		"isEvil":   user.isEvil,
	}

	maps.Copy(vars, data)
	err := HtmlTemplates.ExecuteTemplate(w, templateName, vars)
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
