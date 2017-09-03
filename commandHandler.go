package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/harvestalys/discordgo"
	"github.com/harvestalys/gw2api"
)

var unauthorizedGW2API *gw2api.GW2Api
var language string = "de"

/*
Json lokal Speichern:

[Account unabhängige daten + datum (beim start lesen, fallsälter als 7 tage neu laden, beim neu laden nach neuen ids (größte nummer) ausschau halten - push Nachricht, sonst timer für 7 tage)]

Je Chat User {
User.ID,
API Key,
Push channel id,
Watchlist
}

Beim start Account unabhängige Daten anfragen

Event Benutzer online/offline beachten
Solange Benutzer online alle 5 Minuten api request machen (authenticated gw2 api objekt je user. Falls err! = nil beim nächsten request neu versuchen, sonst halten bis benutzer offline geht?)

commands:
token <account api token> - setzt das token für den user
[need <menge> <item> - fügt Watchlist Eintrag hinzu
done <item> - entfernt Watchlist Eintrag]
refresh - manuelles neu laden der nicht account gebundenen daten

auflisten vorhandener items/Währungen/achievements/rezepte

anfragen von geldbörse [optional eine Währung]/Material/item vorrat
(in bank, gemeinsames Inventar, Materiallager, inventar je charakter, möglich durch rezepte...)
*/

func handleCommand(session *discordgo.Session, message *discordgo.MessageCreate) {

	fmt.Printf("handleCommand(): %s\n", message.Content)

	unauthorizedGW2API = gw2api.NewGW2Api()

	//TODO: multiplexer, command registrieren mit beschreibung, automatisch bei help anzeigen, nur an "den" richtigen handler weitergeben
	handlePing(session, message)
	handleVersion(session, message)
	handleCurrencies(session, message)
	handleCurrency(session, message)
	handleDaily(session, message)
	handleToken(session, message)

	//TODO: SendMessage funktion die bei 2000 zeichen splittet

	//TODO: MessageEmbedImage z.B. icons bei currencies

	//TODO: onUserOnline/Offline wg. watchlist push

	//TODO: wegmarken für Pakt-Vorratsnetz-Agenten (Chatcodes)

	//TODO: implement commands
	//TODO: wenn kein handler zum command passt: sorry message senden mit help cmd
	//TODO: auch auf editierte nachrichten nochmal reagieren (ggf. eigene antwort editieren?)

	//TODO: wie oft welchen boss im raid gelegt (falls keine overall abfrage möglich jede woche vor reset (per timer) abfragen und aufsummieren "seit datum")

	//TODO: watchliste schreiben, stand abfragen, pushen

	//TODO: dailies automatisch morgens abfragen und puschen (ggf. nur mystische schmiede tomorrow), in konfigurierten channel

	//TODO: raidplaner (wer kann wann, welche klassen & rollen (bevorzugt/ausweichklasse) mit klassenerfahrung -> gruppenaufstellung (abhängig vom boss und ob training oder "clear" [z.B. Mo Clear VG, Do Training Gorse])
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

	//TODO: get list/map of all (daily) events if not yet available and look up this eventID
	//idList []int
	//for _,
	//achievements, err := unauthorizedGW2API.AchievementIds()

	return "ID: " + strconv.Itoa(eventID)
}

func listDailiesFromRangeWithHeadline(headline string, category []gw2api.DailyAchievement) string {

	dailies := ""

	for _, daily := range category {

		dailies += getEventTitleByID(daily.ID) + ", Level: " + strconv.Itoa(daily.Level.Min) + " bis " + strconv.Itoa(daily.Level.Max) + " ("

		for i, requirement := range daily.Requirement {

			if i > 0 {
				dailies += ", "
			}

			dailies += requirement
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

			session.ChannelMessageSend(message.ChannelID, dailies)
		} else {

			dailyStruct, err := unauthorizedGW2API.AchievementsDaily()

			if err != nil {

				fmt.Println("GW2: failed to get dailies for today, " + err.Error())
				return
			}

			dailies := resources.Translations["todaysDailies"] + "\n" + handleDailyIntern(dailyStruct, message)

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
