package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type Resources struct {
	Translations map[string]string `json:"Translations"`
}

func (resources *Resources) fromJsonFile() {

	raw, err := ioutil.ReadFile(language + ".json")

	if err != nil {

		fmt.Println("Resource.fromJsonFile(): could not read file:", err.Error())

		if language != "en" {

			// try english as fallback
			fmt.Println("Resource.fromJsonFile(): setting language to \"en\" as fallback")
			language = "en"
			resources.fromJsonFile()
		}

		return
	}

	err = json.Unmarshal(raw, &resources)

	if err != nil {

		fmt.Println("Resource.fromJsonFile(): could not unmarshal data:", err.Error())
		return
	}

	fmt.Println("Resource.fromJsonFile(): resources loaded")
}
