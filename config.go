package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type config struct {
	ClientID     string
	ClientSecret string
	Issuer       string
	Address      string
}

var Config = config{}

func readConfig() {
	file, _ := os.Open("conf.json")
	defer file.Close()
	decoder := json.NewDecoder(file)
	configuration := config{}
	err := decoder.Decode(&configuration)
	Config = configuration
	if err != nil {
		fmt.Println("error:", err)
	}
}
