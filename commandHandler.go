package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/harvestalys/discordgo"
	"github.com/harvestalys/gw2api"
)

func handleCommand(session *discordgo.Session, message *discordgo.MessageCreate) {

	fmt.Printf("handleCommand(): %s\n", message.Content)

	handleHelp(session, message)
	handlePing(session, message)
	handleVersion(session, message)
	handleCurrencies(session, message)
	handleCurrency(session, message)
	handleDaily(session, message)
	handleToken(session, message)
	handleRefresh(session, message)
}

func handleHelp(session *discordgo.Session, message *discordgo.MessageCreate) {

	if strings.Contains(message.Content, "help") {

		explanation := "Available Commands:\n"

		explanation += "ping - answers with \"pong\"\n"
		explanation += "version - the current build version of the GW2 client\n"
		explanation += "currencies - list of all available currencies\n"
		explanation += "currency X - details on currency with id X\n"
		explanation += "daily [tomorrow] [pve|pvp|wvw|fractals|special] - list of daily achievements for today or tomorrow (if specified), full list or subset of either pve, pvp, wvw, fractals or special\n"
		explanation += "token T - sets T as the current users token for GW2API requests that need authorization"
		explanation += "refresh - reloads global data from the API, use this if something is missing (e.g. added to api after last request from the bot)"

		session.ChannelMessageSend(message.ChannelID, explanation)
	}
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

func getEventTitleByID(eventID int) string {

	eventName := achievements[eventID].Name

	if eventName == "" {

		eventName = "ID: " + strconv.Itoa(eventID)
	}

	return eventName
}

func listDailiesFromRangeWithHeadline(headline string, category []gw2api.DailyAchievement) string {

	dailies := ""

	for _, daily := range category {

		dailies += getEventTitleByID(daily.ID)

		if daily.Level.Min != 1 || daily.Level.Max != 80 {
			dailies += ", Level: " + strconv.Itoa(daily.Level.Min)

			if daily.Level.Min != daily.Level.Max {
				dailies += " bis " + strconv.Itoa(daily.Level.Max)
			}
		}

		dailies += " ("

		for i, requirement := range daily.Requirement {

			if i > 0 {
				dailies += ", "
			}

			if requirement == "GuildWars2" {
				dailies += "GW2"
			} else if requirement == "HeartOfThorns" {
				dailies += "HoT"
			} else {
				dailies += requirement
			}
		}

		dailies += ")\n"
	}

	if dailies != "" {

		dailies = headline + "\n" + dailies
	}

	return dailies
}

func listAllDailies(collection gw2api.DailyAchievements) string {

	dailies := listDailiesFromRangeWithHeadline("PvE:", collection.PvE)
	dailies += listDailiesFromRangeWithHeadline("PvP:", collection.PvP)
	dailies += listDailiesFromRangeWithHeadline("WvW:", collection.WvW)
	dailies += listDailiesFromRangeWithHeadline("Fractals:", collection.Fractals)
	dailies += listDailiesFromRangeWithHeadline("Special:", collection.Special)

	return dailies
}

func handleDailyIntern(dailyStruct gw2api.DailyAchievements, message *discordgo.MessageCreate) string {

	dailies := ""

	if strings.Contains(message.Content, " pve") {

		dailies += listDailiesFromRangeWithHeadline("PvE:", dailyStruct.PvE)

	} else if strings.Contains(message.Content, " pvp") {

		dailies += listDailiesFromRangeWithHeadline("PvP:", dailyStruct.PvP)

	} else if strings.Contains(message.Content, " wvw") {

		dailies += listDailiesFromRangeWithHeadline("WvW:", dailyStruct.WvW)

	} else if strings.Contains(message.Content, " fractals") {

		dailies += listDailiesFromRangeWithHeadline("Fractals:", dailyStruct.Fractals)

	} else if strings.Contains(message.Content, " special") {

		dailies += listDailiesFromRangeWithHeadline("Special:", dailyStruct.Special)

	} else {

		dailies += listAllDailies(dailyStruct)
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

			dailies := resources.Translations["tomorrowsDailies"] + "\n" + handleDailyIntern(dailyStruct, message)

			sendLongMessage(session, message, dailies)
		} else {

			dailyStruct, err := unauthorizedGW2API.AchievementsDaily()

			if err != nil {

				fmt.Println("GW2: failed to get dailies for today, " + err.Error())
				return
			}

			dailies := resources.Translations["todaysDailies"] + "\n" + handleDailyIntern(dailyStruct, message)

			sendLongMessage(session, message, dailies)
		}
	}
}

func sendLongMessage(session *discordgo.Session, message *discordgo.MessageCreate, sendMessage string) {

	const tolerance = 150

	i := 0
	maxEnd := len(sendMessage)

	for i < maxEnd {

		start := i
		end := i + maxRunesPerMessage

		if end > maxEnd {

			// send remaining message if rest isn't longer than maxRunesPerMessage
			end = maxEnd
			i = end

		} else {

			// sending only a part of the message, check for last line break before limit
			lb := strings.LastIndex(sendMessage[(end-tolerance):end], "\n")
			if lb != -1 {
				end = end - tolerance + lb
				i = end + 1
			} else {
				i = end
			}

		}

		session.ChannelMessageSend(message.ChannelID, sendMessage[start:end])
	}
}

func handleToken(session *discordgo.Session, message *discordgo.MessageCreate) {

	if strings.Contains(message.Content, "token") {

		substrings := strings.SplitAfter(message.Content, "token ")

		userList.setTokenForUser(message.Author.ID, substrings[1])
	}
}

func handleRefresh(session *discordgo.Session, message *discordgo.MessageCreate) {

	if strings.Contains(message.Content, "refresh") {

		requestAchievements()
	}
}
