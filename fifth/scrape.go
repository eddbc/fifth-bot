package fifth

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/antihax/goesi/esi"
	"github.com/antihax/goesi/optional"
	"github.com/bwmarrin/discordgo"
	"github.com/eddbc/fifth-bot/mux"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func (f *Fifth) ScrapeCorpSupers(ds *discordgo.Session, dm *discordgo.Message, ctx *mux.Context) {

	name := ""

	for k, v := range ctx.Fields {
		if k != 0 {
			name += v
			if k < len(ctx.Fields) {
				name += " "
			}
		}
	}

	id, err := getCorpID(name)
	if err == nil {
		chars, _ := getCorpMemberList(id)
		SendMsgToChan(dm.ChannelID, fmt.Sprintf("Starting search: May take some time due to ZKill rate limits...\n"))
		for _, c := range chars {
			stats, err := zkillstats.get(int32(c.CharacterID))
			if err == nil {
				var supersFlown []string
				for _, st := range stats.TopAllTime {
					if st.Type == "ship" {
						for _, stat := range st.Data {
							ship, _, _ := Eve.UniverseApi.GetUniverseTypesTypeId(context.Background(), int32(stat.ShipTypeID), nil)

							for _, v := range supers {
								if ship.Name == v {
									supersFlown = append(supersFlown, v)
								}
							}
						}
					}
				}
				if supersFlown != nil {
					SendMsgToChan(dm.ChannelID, fmt.Sprintf("%v has flown %v <https://zkillboard.com/character/%v/>", c.Name, supersFlown, c.CharacterID))
				}
			} else {
				log.Printf("Error for %v : %v", c.Name, err)
			}
			time.Sleep(time.Second)
		}
		SendMsgToChan(dm.ChannelID, fmt.Sprintf("Search Finished\n"))
	} else {
		SendMsgToChan(dm.ChannelID, fmt.Sprintf("Error :warning: %s\n", err))
	}

}

func getCorpID(name string) (id int32, e error) {

	bgCtx := context.Background()
	strings := []string{"corporation"}
	res, _, err := Eve.SearchApi.GetSearch(bgCtx, strings, name, &esi.GetSearchOpts{
		Strict: optional.NewBool(true),
	})
	if err != nil {
		log.Printf("error getting corporation id, %s\n", err)
		e = fmt.Errorf("corp search failed")
		return
	}

	if !(len(res.Corporation) > 0) {
		log.Printf("corporation not found\n")
		e = fmt.Errorf("corp not found")
		return
	}

	id = res.Corporation[0]

	return
}

func getCorpMemberList(id int32) (chars []EveWhoCharacter, err error) {
	r, err := http.Get(fmt.Sprintf("https://evewho.com/api/corplist/%v", id))
	if err == nil {
		var res EveWhoCorpLookup

		defer r.Body.Close()
		body, _ := ioutil.ReadAll(r.Body)

		json.Unmarshal(body, &res)
		chars = res.Characters
	}
	return
}

type EveWhoCorpLookup struct {
	Info []struct {
		CorporationID int    `json:"corporation_id"`
		Name          string `json:"name"`
		MemberCount   int    `json:"memberCount"`
	} `json:"info"`
	Characters []EveWhoCharacter `json:"characters"`
}

type EveWhoCharacter struct {
	CharacterID int    `json:"character_id"`
	Name        string `json:"name"`
}
