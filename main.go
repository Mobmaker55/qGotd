package main

import (
	"errors"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"net/http"
	"os"
	"time"
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
	} else {
		fmt.Println(err)
	}
}

func getIndex(w http.ResponseWriter, r *http.Request) {
	if !IsAuthenticated(w, r) {
		http.Redirect(w, r, "http://"+r.Host+"/login", http.StatusFound)
	}
	user, _ := GetUser(r)
	_, err := io.WriteString(w, "hello "+user.PreferredUsername+"!\n")
	if err != nil {
		return
	}
}

func HttpLogger(handlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("")
		fmt.Println(time.Now().String() + " | " + r.Host + r.RequestURI)
		handlerFunc(w, r)
	}
}

// serves template and handles error
func serveTemplate(w http.ResponseWriter, templateName string, data string) {
	err := HtmlTemplates.ExecuteTemplate(w, templateName, data)
	if err != nil {
		fmt.Println("template failed to execute")
		return
	}
}

// initializes templates from the ./html/ directory
func initTemplates() {
	path, _ := os.Getwd()
	path = path + "/html/"
	var templates []string
	err := fs.WalkDir(os.DirFS(path), ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			fmt.Println(err)
			return err
		}
		fmt.Println(path)
		templates = append(templates, path)

		return nil
	})
	if err != nil {
		return
	}
	HtmlTemplates, err = template.ParseFiles(templates...)
}
