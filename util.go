package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

func Handler(handlerFunc http.HandlerFunc) http.HandlerFunc {
	return HttpLogger(func(w http.ResponseWriter, r *http.Request) {
		if !IsAuthenticated(w, r) {
			fmt.Println(r.URL.Path)
			AuthRedirect(w, r, r.URL.Path)
		} else {
			//fmt.Println(runtime.FuncForPC(reflect.ValueOf(handlerFunc).Pointer()).Name())
			handlerFunc(w, r)
		}
	})
}

func HttpLogger(handlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("")
		fmt.Println(time.Now().String() + " | " + r.Host + r.RequestURI)
		handlerFunc(w, r)
	}
}

func AuthRedirect(w http.ResponseWriter, r *http.Request, redir string) {
	prefix := "https://"
	if strings.Contains(Config.Address, "localhost") {
		prefix = "http://"
	}
	http.Redirect(w, r, prefix+r.Host+"/login?redir="+redir, http.StatusFound)
}

func removeResponse(s []Response, index int) []Response {
	return append(s[:index], s[index+1:]...)
}
