package main

import (
"github.com/bwmarrin/discordgo"
"github.com/eddbc/fifth-bot/mux"
"log"
	"context"
	"fmt"
)

type Fifth struct {}

func (f *Fifth) status(ds *discordgo.Session, dm *discordgo.Message, ctx *mux.Context) {

	status, _, err := eve.StatusApi.GetStatus(context.Background(), nil)

	if err != nil {
		log.Printf("error getting TQ status, %s\n", err)
		_, err = ds.ChannelMessageSend(dm.ChannelID, "Error getting TQ status! The logs show nothing...\n")
		return
	}

	if status.Vip {
		_, err = ds.ChannelMessageSend(dm.ChannelID, fmt.Sprintf("Server is in VIP Mode! Current Players: %d\n", status.Players))
	} else {
		_, err = ds.ChannelMessageSend(dm.ChannelID, fmt.Sprintf("Current Players: %d\n", status.Players))
	}

	if err != nil {
		log.Printf("error sending message, %s\n", err)
	}
}
