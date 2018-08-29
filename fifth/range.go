package fifth

import (
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/eddbc/fifth-bot/mux"
	"log"
)

func (f *Fifth) Range(ds *discordgo.Session, dm *discordgo.Message, ctx *mux.Context) {
	if len(ctx.Fields) <= 1 {
		ds.ChannelMessageSend(dm.ChannelID, fmt.Sprintf("<@%v> You need to specify a system you fucking noob", dm.Author.ID))
		return
	}

	strings := []string{"solar_system"}
	res, _, err := Eve.SearchApi.GetSearch(context.Background(), strings, ctx.Fields[1], nil)
	if err != nil {
		log.Printf("Error searching for system: %v", err)
		ds.ChannelMessageSend(dm.ChannelID, fmt.Sprintf("<@%v> Error searching for system", dm.Author.ID))
		return
	}

	if len(res.SolarSystem) == 0 {
		ds.ChannelMessageSend(dm.ChannelID, fmt.Sprintf("<@%v> Couldn't find that system", dm.Author.ID))
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
		ds.ChannelMessageSend(dm.ChannelID, fmt.Sprintf("<@%v> Found %v possible matches : %v", dm.Author.ID, len(res.SolarSystem), systems))
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
		ds.ChannelMessageSend(dm.ChannelID, fmt.Sprintf("<@%v> Ranges from %v: %v", dm.Author.ID, system.Name, ranges))
	}

}

type rangeClass struct {
	size  string
	class string
}
