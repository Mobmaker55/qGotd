package main

import (
	"fmt"
	"net/http"
	"time"
)

func HttpLogger(handlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("")
		fmt.Println(time.Now().String() + " | " + r.Host + r.RequestURI)
		handlerFunc(w, r)
	}
}
