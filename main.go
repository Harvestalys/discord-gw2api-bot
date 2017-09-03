package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/harvestalys/discordgo"
)

var discordToken string
var userList Users
var resources Resources

func init() {

	fmt.Println("init(): parsing arguments")

	flag.StringVar(&discordToken, "token", "", "Authentication Token for the Bot")

	flag.Parse()

	fmt.Printf("init(): token = %s\n", discordToken)
}

func main() {

	fmt.Println("main(): starting")

	if discordToken == "" {

		fmt.Println("main(): authentication token is missing")
		return
	}

	dg, err := discordgo.New("Bot " + discordToken)

	if err != nil {

		fmt.Println("main(): creating discord session failed, " + err.Error())
		return
	}

	// register callback for all MessageCreate events
	dg.AddHandler(onMessage)

	// open websocket to Discord and begin listening
	err = dg.Open()

	if err != nil {

		fmt.Println("main(): connecting to discord failed, " + err.Error())
		return
	}

	userList.fromJsonFile()
	resources.fromJsonFile()

	//TODO: funktioniert noch nicht!?
	dummy := resources.Translations["tomorrowsDailies"]
	fmt.Printf("translation: %s\n", dummy)

	// wait for CTRL+C or other terminate signal is received
	fmt.Println("main(): bot is now running (CTRL+C to exit)")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	fmt.Println("main(): terminate signal received, closing connection")

	dg.Close()

	userList.toJsonFile()

	fmt.Println("main(): exit")
}
