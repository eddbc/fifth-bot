package main

// This file adds the Disgord message route multiplexer, aka "command router".
// to the Disgord bot. This is an optional addition however it is included
// by default to demonstrate how to extend the Disgord bot.

import (
	"github.com/eddbc/fifth-bot/mux"
)

// Router is registered as a global variable to allow easy access to the
// multiplexer throughout the bot.
var Router = mux.New()

func init() {
	// Register the mux OnMessageCreate handler that listens for and processes
	// all messages received.
	Session.AddHandler(Router.OnMessageCreate)

	f := Fifth{}

	// Register the build-in help command.
	Router.Route("help", "Display this message", Router.Help)

	Router.Route("status", "Get EVE Tranquility server status", f.status)
	//Router.Route("who", "Get info about an EVE character", f.who)
	Router.Route("time", "Get current EVE time, or time until a given EVE time (eg. ~time 20:00)", f.eveTime)
}
