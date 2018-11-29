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

	bolt "go.etcd.io/bbolt"

	"github.com/eddbc/fifth-bot/fifth"
	"github.com/eddbc/fifth-bot/sso"
	_ "github.com/joho/godotenv/autoload"
)

// Version is a constant that stores the Disgord version information.
const Version = "v0.1"

const useragent = "fifth-bot, edd_reynolds on slack"

var debug = false

// Session is declared in the global space so it can be easily used
// throughout this program.
// In this use case, there is no error that would be returned.
var Session, _ = discordgo.New()

var eve *esi.APIClient
var httpClient *http.Client

var eveSSOId string
var eveSSOKey string

// Read in all configuration options from both environment variables and
// command line arguments.
func init() {

	flag.BoolVar(&debug, "d", false, "Debug Mode")

	// Set up caching for esi http client (in-memory for now)
	transport := httpcache.NewTransport(httpcache.NewMemoryCache())
	transport.Transport = &http.Transport{Proxy: http.ProxyFromEnvironment}
	httpClient := &http.Client{Transport: transport}
	sso.Client = httpClient

	// Get our API Client.
	esiClient := goesi.NewAPIClient(httpClient, useragent)
	eve = esiClient.ESI
	fifth.Eve = eve

	// Discord Authentication Token
	Session.Token = os.Getenv("FTH_BT_TOKEN")
	if Session.Token == "" {
		flag.StringVar(&Session.Token, "t", "", "Discord Authentication Token")
	}
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

	log.Printf(`Starting...`)

	// Parse command line arguments
	flag.Parse()

	if debug {
		log.Println("debug enabled")
	}

	// Load defined command routes
	routes()

	// Open the fifth.db data file in current directory.
	// It will be created if it doesn't exist.
	log.Printf(`Initialising storage...`)
	db, err := bolt.Open("fifth.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

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

	Session.UpdateStatus(0, "Bot Development")

	fifth.Session = Session
	fifth.Debug = debug
	fifth.DB = db

	fifth.DBInit()

	// Open ZKill websocket for new killmails
	go fifth.ListenZKill()

	go sso.Load(eveSSOId, eveSSOKey)

	// Wait for a CTRL-C
	log.Printf(`Discord bot running`)
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Clean up
	Session.Close()

	// Exit Normally.
}
