package fifth

import (
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/eddbc/fifth-bot/mux"
	"log"
)

//Range Bot command to give dotlan links with the jump ranges fro ma given system
func (f *Fifth) Range(ds *discordgo.Session, dm *discordgo.Message, ctx *mux.Context) {
	if len(ctx.Fields) <= 1 {
		SendMsgToChan(dm.ChannelID, "You need to specify a system you fucking noob")
		return
	}

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

	system, _, err := Eve.UniverseApi.GetUniverseSystemsSystemId(context.Background(), res.SolarSystem[0], nil)
	if err == nil {
		ranges := ""
		classes := []rangeClass{
			{"Capital", "Thanatos"},
			{"Super Capital", "Nyx"},
			{"Rorqual", "Rorqual"},
			{"Black-Ops", "Panther"},
		}
		for _, class := range classes {
			ranges = fmt.Sprintf("%v\n%v : <http://evemaps.dotlan.net/range/%v,5/%v/>", ranges, class.size, class.class, system.Name)
		}
		SendMsgToChan(dm.ChannelID, fmt.Sprintf("Ranges from %v: %v", system.Name, ranges))
	}

}

type rangeClass struct {
	size  string
	class string
}
