package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/harvestalys/discordgo"
	"github.com/harvestalys/gw2api"
)

var discordToken string

var userList Users
var resources Resources
var configuration Configuration

var unauthorizedGW2API *gw2api.GW2Api
var achievements map[int]gw2api.Achievement

var ticker *time.Ticker
var discordSession *discordgo.Session

const serviceNameShort = "discordgw2apibot"
const serviceNameLong = "Discord GW2 API Bot"
const displayName = "T.O.N." // Tyria Notifier and Observer

const maxRunesPerMessage = 2000

func init() {

	fmt.Println("init(): parsing arguments")

	flag.StringVar(&discordToken, "token", "", "Authentication token")

	flag.Parse()

	fmt.Println("init(): token:", discordToken)
}

func main() {

	fmt.Println("main(): starting")

	if discordToken == "" {

		fmt.Println("main(): authentication token is missing")
		return
	}

	var err error
	discordSession, err = discordgo.New("Bot " + discordToken)

	if err != nil {

		fmt.Println("main(): creating discord session failed:", err.Error())
		return
	}

	// register callback for all MessageCreate events
	discordSession.AddHandler(onMessage)

	// open websocket to Discord and begin listening
	err = discordSession.Open()

	if err != nil {

		fmt.Println("main(): connecting to discord failed:", err.Error())
		return
	}

	// read persistent bot data from files
	configuration.fromJsonFile()
	userList.fromJsonFile()
	resources.fromJsonFile(configuration.Language)

	// initialize GW2 API data
	unauthorizedGW2API = gw2api.NewGW2Api()
	requestAchievements()

	checkForNewGW2Version()
	startGW2UpdateWatcher()

	fmt.Println("main(): bot is now running (CTRL+C to exit)")

	// wait for CTRL+C or other terminate signal is received
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	fmt.Println("main(): terminate signal received, closing connection")

	// cleanup
	ticker.Stop()
	discordSession.Close()

	// save persistent bot data to files
	configuration.toJsonFile()
	userList.toJsonFile()

	fmt.Println("main(): exit")
}

func readVersionNumber() int {

	build, err := unauthorizedGW2API.Build()

	if err != nil {
		return configuration.LatestGW2Version
	}

	return build
}

func handleGW2Update(newVersion int) {

	configuration.LatestGW2Version = newVersion

	readFeed(newVersion, configuration.Language)
}

func checkForNewGW2Version() {

	currentVersion := readVersionNumber()
	if currentVersion > configuration.LatestGW2Version {
		handleGW2Update(currentVersion)
	}
}

func startGW2UpdateWatcher() {

	ticker = time.NewTicker(time.Minute * time.Duration(configuration.UpdateCheckMinutes))

	go func() {
		for range ticker.C {
			checkForNewGW2Version()
		}
	}()
}

func requestAchievements() {

	if achievements == nil {

		achievements = make(map[int]gw2api.Achievement)
	}

	fmt.Println("requestAchievements(): requesting achievements...")

	finished := false
	for i := 0; !finished; i++ {

		pageSize := 200

		achievs, err := unauthorizedGW2API.AchievementPages(configuration.Language, i, pageSize)

		if err != nil {

			finished = true

			fmt.Println("GW2: failed to get achievement details for page "+strconv.Itoa(i), err.Error())
		} else {

			for _, a := range achievs {

				achievements[a.ID] = a
			}

			if len(achievs) < pageSize {

				finished = true
			}
		}
	}

	fmt.Println("requestAchievements(): ...finished")
}
