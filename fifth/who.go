package fifth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/antihax/goesi/esi"
	"github.com/antihax/goesi/optional"
	"github.com/bwmarrin/discordgo"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

var supers = []string{"Hel", "Aeon", "Wyvern", "Nyx", "Vendetta", "Revenant", "Avatar", "Erebus", "Leviathan", "Ragnarok", "Molok", "Vanquisher", "Komodo"}
var caps = []string{"Apostle", "Lif", "Ninazu", "Minokawa", "Chimera", "Archon", "Thanatos", "Nidhoggur", "Moros", "Phoenix", "Naglfar", "Revelation", "Vehement"}
var fittedSlots = []int32{27, 34, 19, 26, 11, 118, 87, 92, 98, 125, 132}

func getCharacterInfoEmbed(name string) (*discordgo.MessageEmbed, error) {
	defer timeTrack(time.Now(), "getCharacterInfoEmbed")

	ctx := context.Background()
	strings := []string{"character"}
	res, _, err := Eve.SearchApi.GetSearch(ctx, strings, name, &esi.GetSearchOpts{
		Strict: optional.NewBool(true),
	})
	if err != nil {
		log.Printf("error getting character id, %s\n", err)
		return nil, err
	}

	if len(res.Character) != 1 {
		return nil, errors.New("character not found")
	}

	cid := res.Character[0]
	char, _, err := Eve.CharacterApi.GetCharactersCharacterId(context.Background(), cid, nil)
	if err != nil {
		log.Printf("error getting character details, %s\n", err)
		return nil, err
	}

	corp, _, err := Eve.CorporationApi.GetCorporationsCorporationId(context.Background(), char.CorporationId, nil)
	if err != nil {
		log.Printf("error getting corp details, %s\n", err)
		return nil, err
	}

	fields := []*discordgo.MessageEmbedField{
		{
			Name:   "Corporation",
			Value:  fmt.Sprintf("%s", corp.Name),
			Inline: true,
		},
	}

	if char.AllianceId != 0 {
		alli, _, err := Eve.AllianceApi.GetAlliancesAllianceId(context.Background(), char.AllianceId, nil)
		if err != nil {
			log.Printf("error getting alliance details, %s\n", err)
			return nil, err
		}
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   "Alliance",
			Value:  fmt.Sprintf("%s", alli.Name),
			Inline: true,
		})
	}

	start := time.Now()

	stats, e := zkillstats.get(cid)
	if e == nil {
		topShips := ""
		recentShips := ""
		topSystems := ""
		recentSystems := ""
		var capsFlown []string
		var supersFlown []string
		for _, st := range stats.TopAllTime {
			if st.Type == "ship" {
				for i, stat := range st.Data {
					ship, _, err := Eve.UniverseApi.GetUniverseTypesTypeId(ctx, int32(stat.ShipTypeID), nil)
					if i < 5 {
						if err == nil {
							topShips = fmt.Sprintf("%v%v: %v kills\n", topShips, ship.Name, stat.Kills)
						}
					}

					for _, v := range supers {
						if ship.Name == v {
							supersFlown = append(supersFlown, v)
						}
					}
					for _, v := range caps {
						if ship.Name == v {
							capsFlown = append(capsFlown, v)
						}
					}
				}
			}
		}

		for _, recent := range stats.TopLists[3].Values {
			recentShips = fmt.Sprintf("%v%v: %v kills\n", recentShips, recent.Name, recent.Kills)
		}

		for _, recent := range stats.TopLists[4].Values {
			recentSystems = fmt.Sprintf("%v%v: %v kills\n", recentSystems, recent.Name, recent.Kills)
		}

		for _, recent := range stats.TopAllTime[5].Data[0:5] {
			sys, _, err := Eve.UniverseApi.GetUniverseSystemsSystemId(ctx, int32(recent.SolarSystemID), nil)
			if err == nil {
				topSystems = fmt.Sprintf("%v%v: %v kills\n", topSystems, sys.Name, recent.Kills)
			}
		}

		elapsed := time.Since(start)
		if Debug {
			log.Printf("Stats took %s", elapsed)
		}

		if len(capsFlown) > 0 {
			c := ""
			for _, capital := range capsFlown {
				c = fmt.Sprintf("%v%v\n", c, capital)
			}
			fields = append(fields, &discordgo.MessageEmbedField{
				Name:   "Caps Flown",
				Value:  c,
				Inline: true,
			})
		}

		if len(supersFlown) > 0 {
			s := ""
			for _, soup := range supersFlown {
				s = fmt.Sprintf("%v%v\n", s, soup)
			}
			fields = append(fields, &discordgo.MessageEmbedField{
				Name:   "Supers Flown",
				Value:  s,
				Inline: true,
			})
		}

		if recentShips != "" {
			fields = append(fields, &discordgo.MessageEmbedField{
				Name:   "Recent Ships",
				Value:  fmt.Sprintf("%s", recentShips),
				Inline: true,
			})
		}

		if topShips != "" {
			fields = append(fields, &discordgo.MessageEmbedField{
				Name:   "Top Ships",
				Value:  fmt.Sprintf("%s", topShips),
				Inline: true,
			})
		}

		if recentShips != "" {
			fields = append(fields, &discordgo.MessageEmbedField{
				Name:   "Recent Systems",
				Value:  fmt.Sprintf("%s", recentSystems),
				Inline: true,
			})
		}

		if topShips != "" {
			fields = append(fields, &discordgo.MessageEmbedField{
				Name:   "Top Systems",
				Value:  fmt.Sprintf("%s", topSystems),
				Inline: true,
			})
		}

		//cyno, _ := isCynoChar(cid)
		//if cyno != "" {
		//	fields = append(fields, &discordgo.MessageEmbedField{
		//		Name:  "**Potential Cyno**",
		//		Value: cyno,
		//	})
		//}
	}

	fields = append(fields, &discordgo.MessageEmbedField{
		Name:   "EveWho",
		Value:  fmt.Sprintf("https://evewho.com/pilot/%v/", url.QueryEscape(char.Name)),
		Inline: true,
	})
	fields = append(fields, &discordgo.MessageEmbedField{
		Name:   "zKill",
		Value:  fmt.Sprintf("https://zkillboard.com/character/%v/", cid),
		Inline: true,
	})

	embed := &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{},
		Color:  0x00ff00, // Green
		Fields: fields,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: fmt.Sprintf("https://imageserver.eveonline.com/Character/%d_128.jpg", cid),
		},
		//Timestamp: time.Now().Format(time.RFC3339), // Discord wants ISO8601; RFC3339 is an extension of ISO8601 and should be completely compatible.
		Title: char.Name,
	}

	return embed, nil
}

func isCynoChar(characterID int32) (cyno string, err error) {
	defer timeTrack(time.Now(), "isCynoChar")

	cyno = ""

	r, err := http.Get(fmt.Sprintf("https://zkillboard.com/api/characterID/%v/losses/", characterID))
	if err != nil {
		return cyno, err
	}

	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)

	var kills []Kill
	var cynoLoss int32
	json.Unmarshal(body, &kills)

	//start = time.Now()
	var killFetch time.Duration
	var cynoSearch time.Duration

	for _, k := range kills {

		start := time.Now()
		kill, _, err := Eve.KillmailsApi.GetKillmailsKillmailIdKillmailHash(ctx, k.Zkb.Hash, k.KillmailID, nil)
		elapsed := time.Since(start)
		killFetch += elapsed

		start = time.Now()
		if err == nil {
			for _, item := range kill.Victim.Items {
				if item.ItemTypeId == 21096 || item.ItemTypeId == 28646 {
					for _, slot := range fittedSlots {
						if item.Flag == slot {
							if cynoLoss == 0 {
								cynoLoss = kill.KillmailId
							}
						}
					}
				}
			}
		}
		elapsed = time.Since(start)
		cynoSearch += elapsed
	}

	if Debug {
		log.Printf("time fetching kills %v", killFetch)
		log.Printf("Time searching for cynos %v", cynoSearch)
	}

	if cynoLoss != 0 {
		cyno = fmt.Sprintf("https://zkillboard.com/kill/%v/", cynoLoss)
	}

	return cyno, err
}
