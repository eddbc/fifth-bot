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

	Router.Route("tq", "Get EVE Tranquility server status", f.Status)
	Router.Route("setstatus", "Set 'Playing' status", f.SetStatus)
	Router.Route("time", "Get current EVE time, or time until a given EVE time (eg. ~time 20:00)", f.EveTime)
	Router.Route("who", "Get info about an EVE character", f.Who)
	Router.Route("range", "Get ranges for various ship types from a given system", f.Range)
	if debug {
		Router.Route("test", "", f.Test)
		Router.Route("servers", "", f.Servers)
		Router.Route("caps", "Search public contracts in a region for capitals", f.SearchCapitalContracts)
	}

	Session.AddHandlerOnce(connectedMsg)
}

func connectedMsg(ds *discordgo.Session, _ *discordgo.TypingStart) {
	log.Printf("Connected to %v servers :", len(Session.State.Guilds))
	for _, server := range ds.State.Guilds {
		log.Printf("%v - (%v)", server.Name, server.ID)
	}
}
