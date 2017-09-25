package main

import (
	"fmt"
	"strconv"

	"github.com/mmcdole/gofeed"
)

func readFeed(newVersion int, language string) {

	parser := gofeed.NewParser()
	feed, err := parser.ParseURL("https://forum-" + language + ".guildwars2.com/forum/info/updates.rss")

	if err != nil {
		fmt.Println("Failed to reed feed from guildwars2.com in ", language)

		if language != "en" {
			// try english as fallback
			readFeed(newVersion, "en")
		}

		return
	}

	item := feed.Items[0]

	releaseNotes := resources.Translations["newVersionReleased"] + strconv.Itoa(newVersion) + "\n"
	releaseNotes += item.Title + "\n"
	releaseNotes += item.Published + "\n"
	releaseNotes += item.Description + "\n"
	releaseNotes += item.Link

	fmt.Println(releaseNotes + "\n")

	for _, channelID := range configuration.ChannelIDsForGW2Updates {
		fmt.Println("readFeed(): sending push notification to channel: ", channelID)
		sendLongMessage(discordSession, channelID, releaseNotes)
	}
}
