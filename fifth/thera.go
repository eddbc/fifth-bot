package fifth

import (
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/eddbc/fifth-bot/evescout"
	"github.com/eddbc/fifth-bot/mux"
	"log"
	"sort"
	"strconv"
)

// GetCurrentTheraHoles Bot command to list currently active thera holes, with optional ranges to a target system
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


	if targetSystem == 0 {
		for _, wh := range holes {
			msg = fmt.Sprintf("%v\n%v - %v (%v)",msg, wh.SignatureID, wh.DestinationSolarSystem.Name, wh.DestinationSolarSystem.Region.Name)
		}
	} else {
		var routedWhs []routedWh
		for _, wh := range holes {
			route, _, err := Eve.RoutesApi.GetRouteOriginDestination(context.Background(), wh.DestinationSolarSystem.ID, targetSystem, nil)
			if err == nil {
				jumps := fmt.Sprintf("%v jumps",strconv.Itoa(len(route)))
				m := fmt.Sprintf("%v - %v %v (%v)", wh.SignatureID, jumps, wh.DestinationSolarSystem.Name, wh.DestinationSolarSystem.Region.Name)
				routedWhs = append(routedWhs, routedWh{jumps:len(route), str:m})
			}

		}
		sort.Slice(routedWhs, func(i, j int) bool {
			return routedWhs[i].jumps < routedWhs[j].jumps
		})
		for _, wh := range routedWhs {
			msg = fmt.Sprintf("%v\n%v", msg, wh.str)
		}
	}
	msg = fmt.Sprintf("%v\n```",msg)
	_, _ = SendMsgToChan(dm.ChannelID,msg)
}

type routedWh struct {
	jumps int
	str string
}