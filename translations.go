package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

var resourceFile string = "en.json" //TODO: get startup argument for language (or use en)

type Resources struct {
	Translations map[string]string `json:"Translations"`
}

func (resources *Resources) fromJsonFile() {

	raw, err := ioutil.ReadFile(resourceFile)

	if err != nil {

		fmt.Println("Resources.fromJsonFile(): could not read file, " + err.Error())

		return
	}

	err = json.Unmarshal(raw, &userList)

	if err != nil {

		fmt.Println("Resources.fromJsonFile(): could not unmarshal data, " + err.Error())

		return
	}

	fmt.Println("Resources.fromJsonFile(): resourceFile loaded")
}
