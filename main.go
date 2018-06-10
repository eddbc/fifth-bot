package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"flag"
)

// Version is a constant that stores the Fifth-Bot version information.
const Version = "v0.0.0-alpha"

var (
	commandPrefix string
	botID         string
	token 		  string
)

func init() {
	flag.StringVar(&token, "t", "", "Discord Authentication Token")
}

func main() {

	// Print out a fancy logo!
	fmt.Printf(` 
___________.__  _____  __  .__   __________        __   
\_   _____/|__|/ ____\/  |_|  |__\______   \ _____/  |_ 
 |    __)  |  \   __\\   __\  |  \|    |  _//  _ \   __\
 |     \   |  ||  |   |  | |   Y  \    |   (  <_> )  |  
 \___  /   |__||__|   |__| |___|  /______  /\____/|__|  
     \/                         \/       \/             `+"\n\n")

	flag.Parse()

	discord, err := discordgo.New("Bot "+token)
	errCheck("error creating discord session", err)
	user, err := discord.User("@me")
	errCheck("error retrieving account", err)

	botID = user.ID
	discord.AddHandler(commandHandler)
	discord.AddHandler(func(discord *discordgo.Session, ready *discordgo.Ready) {
		err = discord.UpdateStatus(0, "absolutely nothing.")
		if err != nil {
			fmt.Println("Error attempting to set my status")
		}
		servers := discord.State.Guilds
		fmt.Printf("Fifth-Bot has started on %d servers", len(servers))
	})

	err = discord.Open()
	errCheck("Error opening connection to Discord", err)
	defer discord.Close()

	commandPrefix = "!"

	<-make(chan struct{})

}

func errCheck(msg string, err error) {
	if err != nil {
		fmt.Printf("%s: %+v", msg, err)
		panic(err)
	}
}

func commandHandler(discord *discordgo.Session, message *discordgo.MessageCreate) {
	user := message.Author
	if user.ID == botID || user.Bot {
		//Do nothing because the bot is talking
		return
	}

	//content := message.Content

	fmt.Printf("Message: %+v || From: %s\n", message.Message, message.Author)
}