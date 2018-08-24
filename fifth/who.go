package fifth

import (
	"context"
	"errors"
	"fmt"
	"github.com/antihax/goesi/esi"
	"github.com/antihax/goesi/optional"
	"github.com/bwmarrin/discordgo"
	"log"
	"net/url"
)

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
		for _, st := range stats.TopAllTime {
			if st.Type == "ship" {
				for i := 0; i < 5; i++ {
					stat := st.Data[i]
					ship, _, err := Eve.UniverseApi.GetUniverseTypesTypeId(ctx, int32(stat.ShipTypeId), nil)
					if err == nil {
						topShips = fmt.Sprintf("%v%v: %v kills\n", topShips, ship.Name, stat.Kills)
					}
				}
			}
		}
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:  "Top Ships",
			Value: fmt.Sprintf("%s", topShips),
		})
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
