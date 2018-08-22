// Declare this file to be part of the main package so it can be compiled into
// an executable.
package main

// Import all Go packages required for this file.
import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/antihax/goesi"
	"github.com/antihax/goesi/esi"
	"github.com/bwmarrin/discordgo"
	"github.com/gregjones/httpcache"
)

import _ "github.com/joho/godotenv/autoload"

// Version is a constant that stores the Disgord version information.
const Version = "v0.0.0-alpha"

const useragent = "fifth-bot, edd_reynolds on slack"

var debug = false

var debugChannel = "459341365562572803"     // testing-lab
var spamChannel = "459997163787649025"      // pit of dispair
var importantChannel = "385195528360820739" // general

// Session is declared in the global space so it can be easily used
// throughout this program.
// In this use case, there is no error that would be returned.
var Session, _ = discordgo.New()

var eve *esi.APIClient

var botChannel string

// Read in all configuration options from both environment variables and
// command line arguments.
func init() {

	// Set up caching for esi http client (in-memory for now)
	transport := httpcache.NewTransport(httpcache.NewMemoryCache())
	transport.Transport = &http.Transport{Proxy: http.ProxyFromEnvironment}
	client := &http.Client{Transport: transport}

	// Get our API Client.
	esiClient := goesi.NewAPIClient(client, useragent)
	eve = esiClient.ESI

	// Get Default Bot Channel
	botChannel = os.Getenv("FTH_BT_CHANNEL")

	// Discord Authentication Token
	Session.Token = os.Getenv("FTH_BT_TOKEN")
	if Session.Token == "" {
		flag.StringVar(&Session.Token, "t", "", "Discord Authentication Token")
	}

	flag.BoolVar(&debug, "d", false, "Debug Mode")
}

func main() {

	// Declare any variables needed later.
	var err error

	// Print out a fancy logo!
	fmt.Printf(` 
___________.__  _____  __  .__   __________        __   
\_   _____/|__|/ ____\/  |_|  |__\______   \ _____/  |_ 
 |    __)  |  \   __\\   __\  |  \|    |  _//  _ \   __\
 |     \   |  ||  |   |  | |   Y  \    |   (  <_> )  |  
 \___  /   |__||__|   |__| |___|  /______  /\____/|__|  
     \/                         \/       \/  %-16s`+"\n\n", Version)

	// Parse command line arguments
	flag.Parse()

	if debug {
		log.Println("debug enabled")
	}

	// Verify a Token was provided
	if Session.Token == "" {
		log.Println("You must provide a Discord authentication token.")
		return
	}

	// Verify the Token is valid and grab user information
	Session.State.User, err = Session.User("@me")
	if err != nil {
		log.Printf("error fetching user information, %s\n", err)
	}

	// Open a websocket connection to Discord
	err = Session.Open()
	if err != nil {
		log.Printf("error opening connection to Discord, %s\n", err)
		os.Exit(1)
	}

	// Open ZKill websocket for new killmails
	go listenZKill()

	// Wait for a CTRL-C
	log.Printf(`Now running. Press CTRL-C to exit.`)
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Clean up
	Session.Close()

	// Exit Normally.
}

func sendMsg(msg string) {
	sendMsgToChan(spamChannel, msg)
}

func sendImportantMsg(msg string) {
	sendMsgToChan(importantChannel, msg)
}

func sendDebugMsg(msg string) {
	Session.ChannelMessageSend(debugChannel, msg)
}

func sendMsgToChan(chann string, msg string) {
	log.Println(msg)
	if !debug {
		Session.ChannelMessageSend(chann, msg)
	} else {
		Session.ChannelMessageSend(debugChannel, msg)
	}
}
