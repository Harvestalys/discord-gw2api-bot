package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/harvestalys/discordgo"
	"github.com/harvestalys/gw2api"
)

var discordToken string
var language string

var userList Users
var resources Resources

var unauthorizedGW2API *gw2api.GW2Api
var achievements map[int]gw2api.Achievement

const maxRunesPerMessage = 2000

func init() {

	fmt.Println("init(): parsing arguments")

	flag.StringVar(&discordToken, "token", "", "Authentication token")
	flag.StringVar(&language, "lang", "en", "Language to use")

	flag.Parse()

	fmt.Println("init(): token:", discordToken, " language:", language)
}

func main() {

	fmt.Println("main(): starting")

	if discordToken == "" {

		fmt.Println("main(): authentication token is missing")
		return
	}

	dg, err := discordgo.New("Bot " + discordToken)

	if err != nil {

		fmt.Println("main(): creating discord session failed:", err.Error())
		return
	}

	// register callback for all MessageCreate events
	dg.AddHandler(onMessage)

	// open websocket to Discord and begin listening
	err = dg.Open()

	if err != nil {

		fmt.Println("main(): connecting to discord failed:", err.Error())
		return
	}

	userList.fromJsonFile()
	resources.fromJsonFile()

	unauthorizedGW2API = gw2api.NewGW2Api()
	requestAchievements()

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

func requestAchievements() {

	if achievements == nil {

		achievements = make(map[int]gw2api.Achievement)
	}

	fmt.Println("requestAchievements(): requesting achievements...")

	finished := false
	for i := 0; !finished; i++ {

		pageSize := 200

		achievs, err := unauthorizedGW2API.AchievementPages(language, i, pageSize)

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
