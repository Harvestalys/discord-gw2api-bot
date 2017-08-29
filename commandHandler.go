package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/harvestalys/discordgo"
	"github.com/harvestalys/gw2api"
)

var unauthorizedGW2API *gw2api.GW2Api
var language string = "en"

func handleCommand(session *discordgo.Session, message *discordgo.MessageCreate) {

	fmt.Printf("handleCommand(): %s\n", message.Content)

	unauthorizedGW2API = gw2api.NewGW2Api()

	handlePing(session, message)
	handleVersion(session, message)
	handleCurrencies(session, message)
	handleCurrency(session, message)
	handleDaily(session, message)
	handleToken(session, message)
}

func handlePing(session *discordgo.Session, message *discordgo.MessageCreate) {

	if strings.Contains(message.Content, "ping") {

		session.ChannelMessageSend(message.ChannelID, "Pong!")
	}
}

func handleVersion(session *discordgo.Session, message *discordgo.MessageCreate) {

	if strings.Contains(message.Content, "version") {

		build, err := unauthorizedGW2API.Build()

		if err != nil {

			fmt.Println("GW2: failed to get build, " + err.Error())
			return
		}

		session.ChannelMessageSend(message.ChannelID, "Version: "+strconv.Itoa(build))
	}
}

func handleCurrencies(session *discordgo.Session, message *discordgo.MessageCreate) {

	if strings.Contains(message.Content, "currencies") {

		currencyIDs, err := unauthorizedGW2API.Currencies()

		if err != nil {

			fmt.Println("GW2: failed to get CurrencyIDs, " + err.Error())
			return
		}

		currencies, err := unauthorizedGW2API.CurrencyIds(language, currencyIDs...)

		if err != nil {

			fmt.Println("GW2: failed to get localized Currencies, " + err.Error())
			return
		}

		currenciesReadable := "Currencies:\n"

		for _, currency := range currencies {

			currenciesReadable += strconv.Itoa(currency.ID) + ": " + currency.Name + "\n"
		}

		session.ChannelMessageSend(message.ChannelID, currenciesReadable)
	}
}

func handleCurrency(session *discordgo.Session, message *discordgo.MessageCreate) {

	if strings.Contains(message.Content, "currency ") {

		splitCmd := strings.SplitAfter(message.Content, "currency ")

		currencyID, err := strconv.Atoi(splitCmd[1])

		if err != nil {

			fmt.Println("GW2: failed to get CurrencyID from command, " + err.Error())
			return
		}

		currencies, err := unauthorizedGW2API.CurrencyIds(language, currencyID)

		if err != nil {

			fmt.Println("GW2: failed to get localized Currency, " + err.Error())
			return
		}

		var currenciesReadable string

		// should only contain one element
		for _, currency := range currencies {

			currenciesReadable += strconv.Itoa(currency.ID) + ": " + currency.Name + "\n"
			currenciesReadable += currency.Description + "\n"
			currenciesReadable += "Icon URL: " + currency.Icon + "\n"
		}

		session.ChannelMessageSend(message.ChannelID, currenciesReadable)
	}
}

func listDailiesFromRange(category []gw2api.DailyAchievement) string {

	dailies := ""

	for _, daily := range category {

		dailies += "ID: " + strconv.Itoa(daily.ID) + ", Level: " + strconv.Itoa(daily.Level.Min) + " to " + strconv.Itoa(daily.Level.Max) + " ("

		for i, requirement := range daily.Requirement {

			if i > 0 {
				dailies += ", "
			}

			dailies += requirement
		}

		dailies += ")\n"
	}

	return dailies
}

func listDailies(collection gw2api.DailyAchievements) string {

	dailies := "PvE:\n" + listDailiesFromRange(collection.PvE)
	dailies += "PvP:\n" + listDailiesFromRange(collection.PvP)
	dailies += "WvW:\n" + listDailiesFromRange(collection.WvW)
	dailies += "Fractals:\n" + listDailiesFromRange(collection.Fractals)
	dailies += "Special:\n" + listDailiesFromRange(collection.Special)

	return dailies
}

func handleDailyIntern(dailyStruct gw2api.DailyAchievements, message *discordgo.MessageCreate) string {

	dailies := ""
	listAll := true

	if strings.Contains(message.Content, " pve") {

		dailies += "PvE:\n" + listDailiesFromRange(dailyStruct.PvE)
		listAll = false
	}

	if strings.Contains(message.Content, " pvp") {

		dailies += "PvP:\n" + listDailiesFromRange(dailyStruct.PvP)
		listAll = false
	}

	if strings.Contains(message.Content, " wvw") {

		dailies += "WvW:\n" + listDailiesFromRange(dailyStruct.WvW)
		listAll = false
	}

	if strings.Contains(message.Content, " fractals") {

		dailies += "Fractals:\n" + listDailiesFromRange(dailyStruct.Fractals)
		listAll = false
	}

	if strings.Contains(message.Content, " special") {

		dailies += "Special:\n" + listDailiesFromRange(dailyStruct.Special)
		listAll = false
	}

	if listAll {
		dailies += listDailies(dailyStruct)
	}

	return dailies
}

func handleDaily(session *discordgo.Session, message *discordgo.MessageCreate) {

	if strings.Contains(message.Content, "daily") {

		if strings.Contains(message.Content, "daily tomorrow") {

			dailyStruct, err := unauthorizedGW2API.AchievementsDailyTomorrow()

			if err != nil {

				fmt.Println("GW2: failed to get dailies for tomorrow, " + err.Error())
				return
			}

			dailies := "Tomorrows achievements:\n" + handleDailyIntern(dailyStruct, message)

			session.ChannelMessageSend(message.ChannelID, dailies)
		} else {

			dailyStruct, err := unauthorizedGW2API.AchievementsDaily()

			if err != nil {

				fmt.Println("GW2: failed to get dailies for today, " + err.Error())
				return
			}

			dailies := "Todays achievements:\n" + handleDailyIntern(dailyStruct, message)

			session.ChannelMessageSend(message.ChannelID, dailies)
		}
	}
}

func handleToken(session *discordgo.Session, message *discordgo.MessageCreate) {

	if strings.Contains(message.Content, "token") {

		substrings := strings.SplitAfter(message.Content, "token ")

		userList.setTokenForUser(message.Author.ID, substrings[1])
	}
}
