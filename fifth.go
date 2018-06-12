package main

import (
"github.com/bwmarrin/discordgo"
"github.com/eddbc/fifth-bot/mux"
"log"
)

type Fifth struct {}

func (f *Fifth) test(ds *discordgo.Session, dm *discordgo.Message, ctx *mux.Context) {
	_, err := ds.ChannelMessageSend(dm.ChannelID, "Test String")

	if err != nil {
		log.Printf("error sending message, %s\n", err)
	}
}
