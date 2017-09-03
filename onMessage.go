package main

import (
	"fmt"
	"strings"

	"github.com/harvestalys/discordgo"
)

/*
 * This function will be called (due to AddHandler call in main()) every time a new
 * message is created on any channel that the autenticated bot has access to.
 */
func onMessage(session *discordgo.Session, message *discordgo.MessageCreate) {

	// ignore messages created by the bot
	if message.Author.ID == session.State.User.ID {

		return
	}

	if shouldBotHandleCommand(session, message) {

		// bot was mentioned in message, choose command handler

		handleCommand(session, message)
	}
}

func shouldBotHandleCommand(session *discordgo.Session, message *discordgo.MessageCreate) bool {

	var isBotMentioned = checkForMention(message.Mentions, session.State.User)

	var isPrefix = strings.HasPrefix(message.Content, "!") // use "!gw2api "?

	channel, err := session.State.Channel(message.ChannelID)
	if err != nil {

		fmt.Println("shouldBotHandleCommand(): could not get channel by ID, " + err.Error())
		return isBotMentioned || isPrefix
	}

	return channel.IsPrivate || isBotMentioned || isPrefix
}

func checkForMention(mentions []*discordgo.User, botUser *discordgo.User) bool {

	// check if bot is in the list of mentions
	for _, mentionedUser := range mentions {

		if mentionedUser.ID == botUser.ID {

			// bot was mentioned in message

			return true
		}
	}

	return false
}
