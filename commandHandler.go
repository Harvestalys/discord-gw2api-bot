package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

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
	handleEvent(session, message)
	handleNotify(session, message)
}

func handleNotify(session *discordgo.Session, message *discordgo.MessageCreate) {

	commandHandled := false

	if strings.Contains(message.Content, "notify ") {
		parts := strings.SplitAfter(message.Content, "notify ")

		if strings.HasPrefix(parts[1], "update ") {

			if strings.HasSuffix(parts[1], "on") {
				if configuration.ChannelIDsForGW2Updates == nil {
					fmt.Println("handleNotify(): ChannelIDsForGW2Updates was nil, creating with ChannelID ", message.ChannelID)
					configuration.ChannelIDsForGW2Updates = []string{message.ChannelID}
				} else {
					fmt.Println("handleNotify(): ChannelIDsForGW2Updates already exists, appending ChannelID ", message.ChannelID)
					configuration.ChannelIDsForGW2Updates = append(configuration.ChannelIDsForGW2Updates, message.ChannelID)
				}
				commandHandled = true
				sendLongMessage(session, message.ChannelID, resources.Translations["notifyUpdateOnAck"])

			} else if strings.HasSuffix(parts[1], "off") {

				newChannelList := []string{}
				for _, channelID := range configuration.ChannelIDsForGW2Updates {
					if channelID != message.ChannelID {
						newChannelList = append(newChannelList, channelID)
					}
				}
				configuration.ChannelIDsForGW2Updates = newChannelList

				commandHandled = true
				sendLongMessage(session, message.ChannelID, resources.Translations["notifyUpdateOffAck"])
			}
		}

		if !commandHandled {
			sendLongMessage(session, message.ChannelID, "Nothing happened, please check your command.")
		}
	}
}

func handleEvent(session *discordgo.Session, message *discordgo.MessageCreate) {
	if strings.Contains(strings.ToLower(message.Content), "pof") {

		countdown := ""

		now := time.Now()
		release := time.Date(2017, 9, 22, 18, 0, 0, 0, time.Local)
		diff := release.Sub(now)

		hours := diff.Hours()
		days := math.Floor(hours / 24)
		hours = math.Floor(math.Mod(hours, 24)) // hours -= days * 24
		minutes := math.Ceil(math.Mod(diff.Minutes(), 60))

		if days > 0 {
			countdown += " " + strconv.FormatFloat(days, 'f', 0, 64) + " Tag"
			if days > 1 {
				countdown += "en"
			}
		}

		if hours > 0 {
			countdown += " " + strconv.FormatFloat(hours, 'f', 0, 64) + " Stunde"
			if hours > 1 {
				countdown += "n"
			}
		}

		if minutes > 0 {
			countdown += " " + strconv.FormatFloat(minutes, 'f', 0, 64) + " Minute"
			if minutes > 1 {
				countdown += "n"
			}
		}

		if countdown != "" {
			countdown = "Path of Fire Release in" + countdown
		} else {
			countdown = "Path of Fire is HERE"
		}

		sendLongMessage(session, message.ChannelID, countdown)
	}
}

func handleHelp(session *discordgo.Session, message *discordgo.MessageCreate) {

	if strings.Contains(message.Content, "help") {

		explanation := "Available Commands:\n"

		explanation += "ping - answers with \"pong\"\n"
		explanation += "version - the current build version of the GW2 client\n"
		explanation += "currencies - list of all available currencies\n"
		explanation += "currency X - details on currency with id X\n"
		explanation += "daily [tomorrow] [pve|pvp|wvw|fractals|special] - list of daily achievements for today or tomorrow (if specified), full list or subset of either pve, pvp, wvw, fractals or special\n"
		explanation += "token T - sets T as the current users token for GW2API requests that need authorization\n"
		explanation += "refresh - reloads global data from the API, use this if something is missing (e.g. added to api after last request from the bot)\n"
		explanation += "notify update on|off - (de)activate push notifications on GW2 updates to current channel"

		sendLongMessage(session, message.ChannelID, explanation)
	}
}

func handlePing(session *discordgo.Session, message *discordgo.MessageCreate) {

	if strings.Contains(message.Content, "ping") {

		sendLongMessage(session, message.ChannelID, "Pong!")
	}
}

func handleVersion(session *discordgo.Session, message *discordgo.MessageCreate) {

	if strings.Contains(message.Content, "version") {

		build, err := unauthorizedGW2API.Build()

		if err != nil {

			fmt.Println("GW2: failed to get build, " + err.Error())
			return
		}

		versionMsg := resources.Translations["latestVersion"] + ": " + strconv.Itoa(build)
		sendLongMessage(session, message.ChannelID, versionMsg)
	}
}

func handleCurrencies(session *discordgo.Session, message *discordgo.MessageCreate) {

	if strings.Contains(message.Content, "currencies") {

		currencyIDs, err := unauthorizedGW2API.Currencies()

		if err != nil {

			fmt.Println("GW2: failed to get CurrencyIDs, " + err.Error())
			return
		}

		currencies, err := unauthorizedGW2API.CurrencyIds(configuration.Language, currencyIDs...)

		if err != nil {

			fmt.Println("GW2: failed to get localized Currencies, " + err.Error())
			return
		}

		currenciesReadable := resources.Translations["currencies"] + ":\n"

		for _, currency := range currencies {

			currenciesReadable += strconv.Itoa(currency.ID) + ": " + currency.Name + "\n"
		}

		sendLongMessage(session, message.ChannelID, currenciesReadable)
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

		currencies, err := unauthorizedGW2API.CurrencyIds(configuration.Language, currencyID)

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

		sendLongMessage(session, message.ChannelID, currenciesReadable)
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

			sendLongMessage(session, message.ChannelID, dailies)
		} else {

			dailyStruct, err := unauthorizedGW2API.AchievementsDaily()

			if err != nil {

				fmt.Println("GW2: failed to get dailies for today, " + err.Error())
				return
			}

			dailies := resources.Translations["todaysDailies"] + "\n" + handleDailyIntern(dailyStruct, message)

			sendLongMessage(session, message.ChannelID, dailies)
		}
	}
}

func sendLongMessage(session *discordgo.Session, channelID string, sendMessage string) {

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

		session.ChannelMessageSend(channelID, sendMessage[start:end])
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
