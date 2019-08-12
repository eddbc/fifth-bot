package fifth

import (
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/eddbc/fifth-bot/evescout"
	"github.com/eddbc/fifth-bot/mux"
	"log"
)

func (f *Fifth) GetCurrentTheraHoles(ds *discordgo.Session, dm *discordgo.Message, ctx *mux.Context) {

	targetSystem  := int32(0)
	if len(ctx.Fields) > 1 {
		strings := []string{"solar_system"}
		res, _, err := Eve.SearchApi.GetSearch(context.Background(), strings, ctx.Fields[1], nil)
		if err != nil {
			log.Printf("Error searching for system: %v", err)
			SendMsgToChan(dm.ChannelID, "Error searching for system")
			return
		}

		if len(res.SolarSystem) == 0 {
			SendMsgToChan(dm.ChannelID, "Couldn't find that system")
			return
		}

		if len(res.SolarSystem) > 1 {
			systems := ""
			for _, s := range res.SolarSystem {
				system, _, err := Eve.UniverseApi.GetUniverseSystemsSystemId(context.Background(), s, nil)
				if err == nil {
					systems = fmt.Sprintf("%v \n%v", systems, system.Name)

				}
			}
			SendMsgToChan(dm.ChannelID, fmt.Sprintf("Found %v possible matches : %v", len(res.SolarSystem), systems))
			return
		}

		targetSystem = res.SolarSystem[0]
	}

	holes, err := evescout.GetTheraHoles()

	if err != nil {
		return
	}

	msg := "```"
	for _, wh := range holes {
		//if targetSystem == 0 {
			msg = fmt.Sprintf("%v\n%v - %v  (%v)",msg, wh.SignatureID, wh.DestinationSolarSystem.Name, wh.DestinationSolarSystem.Region.Name)
		//}

	}
	msg = fmt.Sprintf("%v\ntarget system: %v```",msg, targetSystem)
	_, _ = SendMsgToChan(dm.ChannelID,msg)
}

