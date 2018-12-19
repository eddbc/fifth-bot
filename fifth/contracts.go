package fifth

import (
	"context"
	"fmt"
	"github.com/antihax/goesi/esi"
	"github.com/bwmarrin/discordgo"
	"github.com/eddbc/fifth-bot/mux"
	"log"
	"strconv"
	"sync"
)

func (f *Fifth) SearchCapitalContracts(ds *discordgo.Session, dm *discordgo.Message, muxCtx *mux.Context) {
	rgnStr := ""
	for k, v := range muxCtx.Fields {
		if k != 0 {
			rgnStr += v
			if k < len(muxCtx.Fields) {
				rgnStr += " "
			}
		}
	}

	log.Printf("searching for region %v", rgnStr)

	sres, _, err := Eve.SearchApi.GetSearch(ctx, []string{"region"}, rgnStr, nil)
	if err != nil {
		return
	}

	rgnId := sres.Region[0]

	region, _, err := Eve.UniverseApi.GetUniverseRegionsRegionId(ctx, rgnId, nil)
	_, err = ds.ChannelMessageSend(dm.ChannelID,
		fmt.Sprintf("Searching for contracts in %v... (May take a while)", region.Name),
	)

	contracts, err := getContractsForRegion(rgnId)
	if err != nil {
		return
	}

	log.Printf("%v contracts found", len(contracts))
	_, err = ds.ChannelMessageSend(dm.ChannelID, fmt.Sprintf("%v contracts found", len(contracts)))
}

func getAllContracts() {
	log.Println("cron job firing")

	regions, _, err := Eve.UniverseApi.GetUniverseRegions(ctx, nil)

	if err != nil {
		log.Print("error getting regions")
	}

	var wg sync.WaitGroup
	contractChannel := make (chan []esi.GetContractsPublicRegionId200Ok)
	for _, regionId := range regions {
		wg.Add(1)
		go func(regionId int32) {
			defer wg.Done()
			regionContracts, err :=getContractsForRegion(regionId)
			if err != nil {
				log.Fatal(err)
			} else {
				contractChannel <- regionContracts
			}

		}(regionId)
	}
	wg.Wait()
}

func getContractsForRegion(regionId int32) (expCon []esi.GetContractsPublicRegionId200Ok, err error) {
	expCon = make([]esi.GetContractsPublicRegionId200Ok, 0)

	pages := 1
	currentPage := 1

	for currentPage <= pages {
		log.Printf("processing page %v", currentPage)
		c := context.WithValue(ctx, "Page", currentPage)
		contracts, resp, err := Eve.ContractsApi.GetContractsPublicRegionId(c, regionId, nil)
		pages, _ = strconv.Atoi(resp.Header.Get("X-Pages"))
		if err == nil {
			for _, c := range contracts {
				if c.Price > 1000000000 && c.Type_ == "item_exchange" {
					items, _, err := Eve.ContractsApi.GetContractsPublicItemsContractId(ctx, c.ContractId, nil)
					if err != nil {
						return expCon, err
					}
					for _, item := range items {
						t, _, err := Eve.UniverseApi.GetUniverseTypesTypeId(ctx, item.TypeId, nil)
						if err != nil {
							return expCon, err
						}
						for _, v := range superTypes {
							if t.GroupId == v {
								expCon = append(expCon, c)
							}
						}
					}

				}
			}
		}
		currentPage++
	}

	return
}
