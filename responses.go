package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"slices"
	"strings"
)

type Response struct {
	ResponseID int    `json:"responseID"`
	QuestionID int    `json:"questionID"`
	Response   string `json:"response"`
	Username   string `json:"username"`
	Author     string `json:"author"`
	Upvotes    int    `json:"upvotes"`
	IsHidden   bool   `json:"isHidden"`
}

var responseID = 0
var responses map[int]Response
var upvotes map[string][]int

func initResponses() {
	responses = make(map[int]Response)
	upvotes = make(map[string][]int)

}

func responseHandler(w http.ResponseWriter, r *http.Request) {
	//post := r.Method == "POST"
	//TODO: Switch to HTTP method-based API
	switch strings.Split(r.URL.Path, "/")[2] {
	case "add":
		NewResponse(w, r)
		break
	case "get":
		GetResponses(w, r)
		break
	case "upvote":
		UpvoteResponse(w, r)
		break
	case "report":
		ReportResponse(w, r)
		break
	case "delete":
		DeleteResponse(w, r)
		break

	}

}

func NewResponse(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Only POST supported", http.StatusMethodNotAllowed)
	}
	content, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	var result map[string]interface{}

	err = json.Unmarshal(content, &result)
	if err != nil {
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user := GetUser(w, r)
	realUsername := user.PreferredUsername
	if result["username"] != realUsername {
		fmt.Println("User spoof detected for user" + realUsername)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	fullName := user.GivenName + " " + user.FamilyName

	value := fmt.Sprint(result["response"])
	if value == "" {
		return
	}

	newResp := Response{
		ResponseID: responseID,
		QuestionID: 0,
		Response:   value,
		Username:   realUsername,
		Author:     fullName,
		Upvotes:    0,
		IsHidden:   false,
	}
	responses[responseID] = newResp
	responseID = responseID + 1

	var toSend []interface{}
	toSend = append(toSend, []Response{newResp})
	toSend = append(toSend, GetUser(w, r).isEvil)
	toSend = append(toSend, GetUser(w, r).PreferredUsername)

	send, _ := json.Marshal(toSend)
	w.Write(send)

}

func GetResponses(w http.ResponseWriter, r *http.Request) {
	var sendable []Response

	for i := range responses { //don't send it if it's hidden
		if responses[i].IsHidden && !GetUser(w, r).isEvil {
			break
		}
		sendable = append(sendable, responses[i])
	}

	slices.SortFunc(sendable, func(a, b Response) int { //sort by upvotes so the most voted up is the winner
		if a.Upvotes > b.Upvotes {
			return 1
		}
		if a.Upvotes < b.Upvotes {
			return -1
		}
		return 0
	})

	var toSend []interface{} //combine responses, EBoard, and Username into one form
	toSend = append(toSend, sendable)
	toSend = append(toSend, GetUser(w, r).isEvil)
	toSend = append(toSend, GetUser(w, r).PreferredUsername)

	body, err := json.Marshal(toSend) //make it json
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(body) //send JSON

}

func UpvoteResponse(w http.ResponseWriter, r *http.Request) {
	content, err := io.ReadAll(r.Body) //get data
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	var c map[string]interface{}
	err = json.Unmarshal(content, &c)
	id := c["responseID"].(int) //find the ID
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	response := responses[id] //get response object
	w.WriteHeader(http.StatusOK)

	username := GetUser(w, r).PreferredUsername
	upvoteSlice := upvotes[username] //check if the user has already upvoted
	if username == response.Username {
		return
	}
	if slices.Contains(upvoteSlice, id) {
		return
	}
	response.Upvotes++ //upvote the user
	upvotes[username] = append(upvoteSlice, id)
	responses[id] = response

}

func ReportResponse(w http.ResponseWriter, r *http.Request) {
	content, err := io.ReadAll(r.Body) //get data
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	var c map[string]interface{}
	err = json.Unmarshal(content, &c) //de-json it
	id := c["responseID"]             //find response ID
	if err != nil {
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	fmt.Println("Report sent by: "+GetUser(w, r).PreferredUsername+" on response", id) // TODO: add report functionality
	w.WriteHeader(http.StatusOK)
}

func DeleteResponse(w http.ResponseWriter, r *http.Request) {
	content, err := io.ReadAll(r.Body) //get data
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	var c map[string]float64
	err = json.Unmarshal(content, &c)
	id := int(c["responseID"]) //parse out to int
	if err != nil {
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	user := GetUser(w, r)
	response := responses[id]
	if !(user.PreferredUsername == response.Username || user.isEvil) {
		w.WriteHeader(http.StatusUnauthorized) //nullify delete if user spoofing is detected
		return
	}
	w.WriteHeader(http.StatusOK)
	delete(responses, id) //remove from map
}
