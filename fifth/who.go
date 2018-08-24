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
)

var supers = []string{"Hel", "Aeon", "Wyvern", "Nyx", "Vendetta", "Revenant", "Avatar", "Erebus", "Leviathan", "Ragnarok", "Molok", "Vanquisher", "Komodo"}
var caps = []string{"Apostle", "Lif", "Ninazu", "Minokawa", "Chimera", "Archon", "Thanatos", "Nidhoggur", "Moros", "Phoenix", "Naglfar", "Revelation", "Vehement"}
var fittedSlots = []int{27, 34, 19, 26, 11, 118, 87, 92, 98, 125, 132}

func getCharacterInfoEmbed(name string) (*discordgo.MessageEmbed, error) {

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

	stats, e := zkillstats.get(cid)
	if e == nil {
		topShips := ""
		var capsFlown []string
		var supersFlown []string
		for _, st := range stats.TopAllTime {
			if st.Type == "ship" {
				for i, stat := range st.Data {
					ship, _, err := Eve.UniverseApi.GetUniverseTypesTypeId(ctx, int32(stat.ShipTypeId), nil)
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

		if len(capsFlown) > 0 {
			c := ""
			for _, cap := range capsFlown {
				c = fmt.Sprintf("%v%v\n", c, cap)
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

		if topShips != "" {
			fields = append(fields, &discordgo.MessageEmbedField{
				Name:   "Top Ships",
				Value:  fmt.Sprintf("%s", topShips),
				Inline: true,
			})
		}

		cyno, _ := isCynoChar(cid)
		if cyno != "" {
			fields = append(fields, &discordgo.MessageEmbedField{
				Name:  "**Potential Cyno**",
				Value: cyno,
			})
		}
	}

	fields = append(fields, &discordgo.MessageEmbedField{
		Name:   "zKill",
		Value:  fmt.Sprintf("https://zkillboard.com/character/%v/", cid),
		Inline: true,
	})
	fields = append(fields, &discordgo.MessageEmbedField{
		Name:   "EveWho",
		Value:  fmt.Sprintf("https://evewho.com/pilot/%v/", url.QueryEscape(char.Name)),
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
	cyno = ""

	r, err := http.Get(fmt.Sprintf("https://zkillboard.com/api/characterID/%v/losses/", characterID))
	if err != nil {
		return cyno, err
	}

	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)

	var kills []Kill
	cynoLoss := 0
	json.Unmarshal(body, &kills)

	for _, kill := range kills {
		for _, item := range kill.Victim.Items {
			if item.ItemTypeID == 21096 {
				for _, slot := range fittedSlots {
					if item.Flag == slot {
						if cynoLoss == 0 {
							cynoLoss = kill.KillmailID
						}
					}
				}
			}
		}
	}

	if cynoLoss != 0 {
		cyno = fmt.Sprintf("https://zkillboard.com/kill/%v/", cynoLoss)
	}

	return cyno, err
}
