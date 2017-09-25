package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

var configurationFile string = "configuration.json"

type Configuration struct {
	Language                string   `json:"Language"`
	LatestGW2Version        int      `json:"LatestGW2Version"`
	UpdateCheckMinutes      int      `json:"UpdateCheckMinutes"`
	ChannelIDsForGW2Updates []string `json:"ChannelIDsForGW2Updates`
}

func (configuration *Configuration) toJsonFile() {

	bytes, err := json.Marshal(configuration)

	if err != nil {

		fmt.Println("Configuration.toJsonFile(): could not marshal data:", err.Error())

		return
	}

	fmt.Printf("%s", bytes)

	ioutil.WriteFile(configurationFile, bytes, 0666)

	fmt.Println("Configuration.toJsonFile(): configuration saved")
}

func (configuration *Configuration) fromJsonFile() {

	raw, err := ioutil.ReadFile(configurationFile)

	// error on first start expected because there will be no user file
	if err != nil {

		fmt.Println("Configuration.fromJsonFile(): could not read file:", err.Error())
		configuration.initialize()
		return
	}

	err = json.Unmarshal(raw, &configuration)

	if err != nil {

		fmt.Println("Configuration.fromJsonFile(): could not unmarshal data:", err.Error())
		configuration.initialize()
		return
	}

	fmt.Println("Configuration.fromJsonFile(): configuration loaded")
}

func (configuration *Configuration) initialize() {
	configuration.Language = "en"
	configuration.LatestGW2Version = 0
	configuration.UpdateCheckMinutes = 2
}
