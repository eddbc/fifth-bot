package fifth

import (
	"context"
	"fmt"
	"github.com/antihax/goesi/esi"
	"github.com/bwmarrin/discordgo"
	"github.com/eddbc/fifth-bot/esistatus"
	"github.com/eddbc/fifth-bot/mux"
	"log"
	"strings"
	"time"
)

//Eve ESI Client
var Eve *esi.APIClient

//Session Discord client session
var Session *discordgo.Session

//Debug Debug mode enabled
var Debug bool

const useragent = "fifth-bot, edd_reynolds on slack"

//Fifth Discord Bot Main Struct
type Fifth struct{}

//Status Bot command to get EVE Online server status
func (f *Fifth) Status(ds *discordgo.Session, dm *discordgo.Message, ctx *mux.Context) {

	msg := ""

	status, _, err := Eve.StatusApi.GetStatus(context.Background(), nil)

	if err != nil {
		log.Printf("error getting TQ status, %s\n", err)
		_, err = SendMsgToChan(dm.ChannelID, "Error getting TQ status! The logs show nothing...\n")
		return
	}

	if status.Vip && status.Players == 0 {
		msg = fmt.Sprintf("TQ is offline :(\n")
	} else if status.Vip {
		msg = fmt.Sprintf("TQ is in VIP Mode. Current Players: %d\n", status.Players)
	} else {
		msg = fmt.Sprintf("Disgusting subhumans currently on TQ: %d\n", status.Players)
	}

	g, y, r, err := esistatus.GetEsiStatus()
	if err != nil {
		msg = fmt.Sprintf("%vError getting ESI status\n", msg)
	} else {
		if r > 0 || g < 150 {
			msg = fmt.Sprintf("%v:warning: **ESI is reporting errors** :warning:", msg)
		} else if y > 0 {
			msg = fmt.Sprintf("%v*ESI is reporting minor issues*", msg)
		} else {
			msg = fmt.Sprintf("%vESI is :ok_hand:", msg)
		}

		if r > 0 || y > 0 {
			msg = fmt.Sprintf("%v\n(ESI Endpoints - %v :white_check_mark:,  %v :warning:,  %v :broken_heart:)\n", msg, g, y, r)
		}
	}

	_, err = SendMsgToChan(dm.ChannelID, msg)
}

//EveTime bot command to give current EVE time, or time until given EVE time
func (f *Fifth) EveTime(ds *discordgo.Session, dm *discordgo.Message, ctx *mux.Context) {

	//init the loc
	loc, _ := time.LoadLocation("Atlantic/Reykjavik")
	et := time.Now().In(loc)
	etStr := et.Format("15:04")

	if len(ctx.Fields) == 1 {
		SendMsgToChan(dm.ChannelID, fmt.Sprintf("Current EVE Time: **%v**\n", etStr))
	} else {
		dd := et.Format("2006/01/02")
		tt := ctx.Fields[1]
		target, _ := time.ParseInLocation("2006/01/02 15:04", fmt.Sprintf("%v %v", dd, tt), loc)
		if target.Before(et) {
			target = target.AddDate(0, 0, 1)
		}

		log.Printf("target time %v for input %v", target, ctx.Fields[1])
		timeTil := target.Sub(et)
		SendMsgToChan(dm.ChannelID, fmt.Sprintf("Time until %v EVE: **%v**. You should probably learn simple maths and figure it out yourself though.\n(Current EVE Time: %v)", target.Format("15:04"), fmtDuration(timeTil), etStr))
	}

}

//Who Bot command to give information about a given character
func (f *Fifth) Who(ds *discordgo.Session, dm *discordgo.Message, ctx *mux.Context) {

	name := ""

	for k, v := range ctx.Fields {
		if k != 0 {
			name += v
			if k < len(ctx.Fields) {
				name += " "
			}
		}
	}

	embed, err := getCharacterInfoEmbed(name)

	if err == nil {
		_, err := ds.ChannelMessageSendEmbed(dm.ChannelID, embed)
		if err != nil {
			log.Printf("error: %+v", err)
		}
	} else {
		SendMsgToChan(dm.ChannelID, fmt.Sprintf("Error :warning: %s\n", err))
	}
}

//SetStatus Bot command to set the 'currently playing' status of the bot
func (f *Fifth) SetStatus(ds *discordgo.Session, dm *discordgo.Message, ctx *mux.Context) {
	status := ""

	for k, v := range ctx.Fields {
		if k != 0 {
			status += v
			if k < len(ctx.Fields) {
				status += " "
			}
		}
	}

	_ = ds.UpdateGameStatus(0, status)
}

func (f *Fifth) ListZKillTracked(ds *discordgo.Session, dm *discordgo.Message, ctx *mux.Context) {

	entitiesOfInterest = getTrackedEntities()
	entStr := ""
	for _, e := range entitiesOfInterest {
		entStr = fmt.Sprintf("%v\n%v", entStr, e.name)
	}
	_, _ = SendMsgToChan(dm.ChannelID, fmt.Sprintf("```%v```", entStr))
}

func (f *Fifth) AddZKillTracked(ds *discordgo.Session, dm *discordgo.Message, ctx *mux.Context) {

	name := ""

	for k, v := range ctx.Fields {
		if k != 0 {
			name += v
			if k < len(ctx.Fields) {
				name += " "
			}
		}
	}

	err := addTrackedEntityByName(strings.TrimSpace(name))
	if err == nil {
		_, _ = SendMsgToChan(dm.ChannelID, fmt.Sprintf("Added to list"))
	} else {
		_, _ = SendMsgToChan(dm.ChannelID, fmt.Sprintf("Error: %v", err))
	}
}

//Servers Bot command to list currently connected servers
func (f *Fifth) Servers(ds *discordgo.Session, dm *discordgo.Message, ctx *mux.Context) {
	msg := fmt.Sprintf("**Connected to %v servers :**", len(Session.State.Guilds))
	for _, guild := range ds.State.Guilds {
		msg = fmt.Sprintf("%v\n * %v - (%v)", msg, guild.Name, guild.ID)
	}
	m, _ := SendDebugMsg(msg)
	Session.MessageReactionAdd(m.ChannelID, m.ID, "509447602291605519")
}

//EmoteTest Bot command to test emote reacting
func (f *Fifth) EmoteTest(ds *discordgo.Session, dm *discordgo.Message, ctx *mux.Context) {
	m, _ := SendMsgToChan(dm.ChannelID, "reaction test")
	err := Session.MessageReactionAdd(m.ChannelID, m.ID, ":rip:486665154356969496")

	if err != nil {
		log.Printf("error: %+v", err)
	}
}

func fmtDuration(d time.Duration) string {
	d = d.Round(time.Minute)
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	return fmt.Sprintf("%dh %dm", h, m)
}

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	if Debug {
		log.Printf("%s took %s", name, elapsed)
	}
}
