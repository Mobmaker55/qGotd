package main

import (
	"net/http"
	"strings"
)

type Question struct {
	QuestionID    int    `json:"questionID"`
	Question      string `json:"question"`
	Author        string `json:"author"`
	DateScheduled string `json:"dateScheduled"`
}

var questionID = 0
var questions map[int]Question

func initQuestions() {
	questions = make(map[int]Question)
	questions[0] = Question{
		QuestionID:    0,
		Question:      "Do you like this app?",
		Author:        "mob",
		DateScheduled: "",
	}
}

func questionHandler(w http.ResponseWriter, r *http.Request) {
	//TODO: Switch to a HTTP method based API
	switch strings.Split(r.URL.Path, "/")[2] {
	case "add":
		NewQuestion(w, r)
		break
	case "get":
		GetQuestions(w, r)
		break
	case "delete":
		DeleteQuestion(w, r)
		break
	}
}

func NewQuestion(w http.ResponseWriter, r *http.Request) {

}

func GetQuestions(w http.ResponseWriter, r *http.Request) {

}

func DeleteQuestion(w http.ResponseWriter, r *http.Request) {

}
