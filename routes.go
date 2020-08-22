package main

// This file adds the Disgord message route multiplexer, aka "command router".
// to the Disgord bot. This is an optional addition however it is included
// by default to demonstrate how to extend the Disgord bot.

import (
	"github.com/bwmarrin/discordgo"
	"github.com/eddbc/fifth-bot/fifth"
	"github.com/eddbc/fifth-bot/mux"
	"log"
)

// Router is registered as a global variable to allow easy access to the
// multiplexer throughout the bot.
var Router = mux.New()

func routes() {
	// Register the mux OnMessageCreate handler that listens for and processes
	// all messages received.
	Session.AddHandler(Router.OnMessageCreate)

	f := fifth.Fifth{}

	// Register the build-in help command.
	Router.Route("help", "Display this message", Router.Help)

	Router.Route("status", "Get EVE Tranquility server status", f.Status)
	Router.Route("setstatus", "", f.SetStatus)
	Router.Route("time", "Get current EVE time, or time until a given EVE time (eg. !time 20:00)", f.EveTime)
	Router.Route("who", "Get info about an EVE character", f.Who)
	Router.Route("thera", "List Current Thera Holes", f.GetCurrentTheraHoles)
	Router.Route("range", "Get ranges for various ship types from a given system", f.Range)
	Router.Route("zkill-list", "List zKill feed tracked entities", f.ListZKillTracked)
	Router.Route("zkill-add", "Add to the list of tracked zKill entities", f.AddZKillTracked)
	Router.Route("scrape", "Scrape a corp for supers/titans", f.ScrapeCorpSupers)

	// Disabled Timer Commands
	//Router.Route("timer", "Add a timer with format: \"!timer 2d 6h 15m description and location goes here\"", f.AddTimer)
	//Router.Route("timers", "List Timers with format \"description - date (time left) [ID]\"", f.ListTimers)
	//Router.Route("timer-remove", "Remove a timer with a given ID", f.RemoveTimer)

	// Unlisted commands (A command with no description will not show up in the !help list of commands
	Router.Route("servers", "", f.Servers)

	// Debug commands (only enabled in debug mode)
	if debug {
		Router.Route("test", "", f.EmoteTest)
		Router.Route("supers", "Search public contracts in a region for super capitals", f.SearchCapitalContracts)
	}

	Session.AddHandlerOnce(connectedMsg)
}

func connectedMsg(ds *discordgo.Session, _ *discordgo.TypingStart) {
	log.Printf("Connected to %v servers :", len(Session.State.Guilds))
	for _, server := range ds.State.Guilds {
		log.Printf("%v - (%v)", server.Name, server.ID)
	}
}
