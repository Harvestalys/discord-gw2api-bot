package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

var userFile string = "userList.json"

type Users struct {
	UserMap map[string]User `json:"UserList"`
}

type User struct {
	ID    string `json:"DiscordID"`
	Token string `json:"GW2APIToken"`
}

func (userList *Users) toJsonFile() {

	bytes, err := json.Marshal(userList)

	if err != nil {

		fmt.Println("Users.toJsonFile(): could not marshal data, " + err.Error())

		return
	}

	fmt.Printf("%s", bytes)

	ioutil.WriteFile(userFile, bytes, 0666)

	fmt.Println("Users.toJsonFile(): userList saved")
}

func (userList *Users) fromJsonFile() {

	raw, err := ioutil.ReadFile(userFile)

	// error on first start expected because there will be no user file
	if err != nil {

		fmt.Println("Users.fromJsonFile(): could not read file, " + err.Error())

		return
	}

	err = json.Unmarshal(raw, &userList)

	if err != nil {

		fmt.Println("Users.fromJsonFile(): could not unmarshal data, " + err.Error())

		return
	}

	fmt.Println("Users.fromJsonFile(): userList loaded")
}

func (userList *Users) setTokenForUser(userID string, token string) {

	if userList.UserMap == nil {

		userList.UserMap = make(map[string]User)
	}

	userList.UserMap[userID] = User{userID, token}
}
