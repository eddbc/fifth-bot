package main

import (
	"github.com/bwmarrin/discordgo"
	"github.com/eddbc/fifth-bot/mux"
	"log"
	"context"
	"fmt"
	"time"
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
		_, err = ds.ChannelMessageSend(dm.ChannelID, fmt.Sprintf("Disgusting subhumans currently on TQ: %d\n", status.Players))
	}

	if err != nil {
		log.Printf("error sending message, %s\n", err)
	}
}

func (f *Fifth) eveTime (ds *discordgo.Session, dm *discordgo.Message, ctx *mux.Context) {

	//init the loc
	loc, _ := time.LoadLocation("Atlantic/Reykjavik")
	et := time.Now().In(loc)
	etStr := et.Format("15:04")

	if len(ctx.Fields) == 1 {
		_, _ = ds.ChannelMessageSend(dm.ChannelID, fmt.Sprintf("Current EVE Time: **%v**\n", etStr))
	} else {
		dd := et.Format("2006/01/02")
		tt := ctx.Fields[1]
		target, _ := time.ParseInLocation("2006/01/02 15:04", fmt.Sprintf("%v %v", dd, tt), loc)
		if target.Before(et) {
			target = target.AddDate(0,0,1)
		}

		log.Printf("target time %v for input %v", target, ctx.Fields[1])
		timeTil := target.Sub(et)
		_, _ = ds.ChannelMessageSend(dm.ChannelID, fmt.Sprintf("Time until %v EVE: **%v**. You should probably learn simple maths and figure it out yourself though.\n(Current EVE Time: %v)", target.Format("15:04"), fmtDuration(timeTil), etStr))
	}

}

func fmtDuration(d time.Duration) string {
	d = d.Round(time.Minute)
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	return fmt.Sprintf("%dh %dm", h, m)
}

func (f *Fifth) who(ds *discordgo.Session, dm *discordgo.Message, ctx *mux.Context) {

	name := ""

	for k, v := range ctx.Fields {
		if k != 0 {
			name+=v
			if k<len(ctx.Fields) {
				name+=" "
			}
		}
	}

	embed, err := getCharacterInfoEmbed(name)

	if err == nil{
		ds.ChannelMessageSendEmbed(dm.ChannelID, embed)
	} else {
		ds.ChannelMessageSend(dm.ChannelID, fmt.Sprintf("Error :warning: %s\n", err))
	}
}