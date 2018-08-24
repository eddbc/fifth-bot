package fifth

import (
	"github.com/robfig/cron"
	"context"
	"log"
	"sync"
	"github.com/antihax/goesi/esi"
	"github.com/bwmarrin/discordgo"
	"github.com/eddbc/fifth-bot/mux"
	"fmt"
)

var C *cron.Cron

var ctx context.Context

func init() {
	ctx = context.Background()
	C = cron.New()

	//C.AddFunc("@every 10s", func(){getAllContracts()})

	C.Start()

}
func (f *Fifth) SearchCapitalContracts(ds *discordgo.Session, dm *discordgo.Message, muxCtx *mux.Context) {
	region := ""
	for k, v := range muxCtx.Fields {
		if k != 0 {
			region += v
			if k < len(muxCtx.Fields) {
				region += " "
			}
		}
	}

	log.Printf("searching for region %v", region)

	r, _, err := Eve.SearchApi.GetSearch(ctx, []string{"region"}, region, nil)
	if err != nil {
		return
	}

	contracts, err := getContractsForRegion(r.Region[0])
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
	for _, regionId := range regions {
		wg.Add(1)
		go func(regionId int32) {
			defer wg.Done()
			getContractsForRegion(regionId)
		}(regionId)
	}
	wg.Wait()
}

func getContractsForRegion(regionId int32) ([]esi.GetContractsPublicRegionId200Ok, error) {
	expCon := make([]esi.GetContractsPublicRegionId200Ok, 0)
	contracts, _, err := Eve.ContractsApi.GetContractsPublicRegionId(ctx, regionId, nil)
	if err == nil {
		for _, c := range contracts {
			if c.Price > 1000000000 && c.Type_ == "item_exchange" {
				expCon = append(expCon, c)
			}
		}
	}
	return expCon, err
}