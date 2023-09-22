package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

var Mux *http.ServeMux

func main() {
	readConfig()
	Mux = http.NewServeMux()
	initAuth()

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
