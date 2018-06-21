package main

import (
	"github.com/antihax/goesi/esi"
	"github.com/antihax/goesi/optional"
	"log"
	"github.com/bwmarrin/discordgo"
	"fmt"
	"context"
	"errors"
)

func getCharacterInfoEmbed(name string) (*discordgo.MessageEmbed, error ) {

	res, _, err := eve.SearchApi.GetSearch(context.Background(), []string{"character"}, name, &esi.GetSearchOpts{
		Strict:optional.NewBool(true),
	})
	if err != nil {
		log.Printf("error getting character id, %s\n", err)
		return nil, err
	}

	if len(res.Character) != 1 {
		return nil, errors.New("character not found")
	}

	cid := res.Character[0]
	char, _, err := eve.CharacterApi.GetCharactersCharacterId(context.Background(), cid, nil)
	if err != nil {
		log.Printf("error getting character details, %s\n", err)
		return nil, err
	}

	corp, _, err := eve.CorporationApi.GetCorporationsCorporationId(context.Background(), char.CorporationId, nil)
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
		alli, _, err := eve.AllianceApi.GetAlliancesAllianceId(context.Background(), char.AllianceId, nil)
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
		for _, st := range stats.TopAllTime {
			if st.Type == "ship" {
				log.Printf("%v", st.Data)
			}
		}
	}

	embed := &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{},
		Color:  0x00ff00, // Green
		Fields: fields,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL:  fmt.Sprintf("https://imageserver.eveonline.com/Character/%d_128.jpg",cid),
		},
		//Timestamp: time.Now().Format(time.RFC3339), // Discord wants ISO8601; RFC3339 is an extension of ISO8601 and should be completely compatible.
		Title:     char.Name,
	}

	return embed, nil
}
