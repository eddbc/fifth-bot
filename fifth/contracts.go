package fifth

import (
	"fmt"
	"github.com/antihax/goesi/esi"
	"github.com/antihax/goesi/optional"
	"github.com/bwmarrin/discordgo"
	"github.com/eddbc/fifth-bot/isk"
	"github.com/eddbc/fifth-bot/mux"
	"log"
	"strconv"
	"sync"
)

//SearchCapitalContracts Bot command to fetch super-capital contracts
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

	rgnID := sres.Region[0]

	region, _, err := Eve.UniverseApi.GetUniverseRegionsRegionId(ctx, rgnID, nil)
	_, err = ds.ChannelMessageSend(dm.ChannelID,
		fmt.Sprintf("Searching for contracts in %v... (May take a while)", region.Name),
	)

	contracts, err := getContractsForRegion(rgnID)
	if err != nil {
		return
	}

	log.Printf("%v contracts found", len(contracts))
	_, err = ds.ChannelMessageSend(dm.ChannelID, fmt.Sprintf("%v contracts found", len(contracts)))

	msg := ""
	for _, contract := range contracts {
		msg = msg + makeContractMessage(contract) + "\n"
	}
	_, err = ds.ChannelMessageSend(dm.ChannelID, msg)
}

func getAllContracts() {
	regions, _, err := Eve.UniverseApi.GetUniverseRegions(ctx, nil)

	if err != nil {
		log.Print("error getting regions")
	}

	var wg sync.WaitGroup
	contractRegions := make(chan []contractReport, len(regions))
	for _, regionID := range regions {
		wg.Add(1)
		go func(regionId int32) {
			defer wg.Done()
			rc, err := getContractsForRegion(regionId)
			if err != nil {
				log.Print(err)
			} else {
				contractRegions <- rc
			}

		}(regionID)
	}
	wg.Wait()
	close(contractRegions)
	for region := range contractRegions {
		for _, contract := range region {
			log.Print(makeContractMessage(contract))
		}
	}
}

func getContractsForRegion(regionID int32) (superContracts []contractReport, err error) {
	superContracts = make([]contractReport, 0)

	pages := 1
	currentPage := 1

	for currentPage <= pages {
		log.Printf("processing page %v", currentPage)
		contracts, resp, err := Eve.ContractsApi.GetContractsPublicRegionId(ctx, regionID, &esi.GetContractsPublicRegionIdOpts{
			Page: optional.NewInt32(int32(currentPage)),
		})
		pages, _ = strconv.Atoi(resp.Header.Get("X-Pages"))
		if err == nil {
			for _, contract := range contracts {
				if contract.Type_ == "item_exchange" && contract.Volume > 1300000 {
					items, _, err := Eve.ContractsApi.GetContractsPublicItemsContractId(ctx, contract.ContractId, nil)
					if err != nil {
						return superContracts, err
					}

					isSuper := false
					superType := ""
					itemList := ""

					for _, item := range items {
						t, _, err := Eve.UniverseApi.GetUniverseTypesTypeId(ctx, item.TypeId, nil)
						if err != nil {
							return superContracts, err
						}
						itemList = fmt.Sprintf("%v %v %v,", itemList, t.Name, item.Quantity)
						for _, v := range superTypes {
							if t.GroupId == v {
								isSuper = true
								superType = t.Name
							}
						}
					}
					if isSuper {
						superContracts = append(superContracts, contractReport{
							super:    superType,
							region:   regionID,
							contract: contract,
							itemList: itemList,
						})
					}

				}
			}
		}
		currentPage++
	}

	return
}

func makeContractMessage(cr contractReport) (msg string) {

	return fmt.Sprintf("*%v* - %v <url=contract:%v//%v>", cr.super, isk.NearestThousandFormat(cr.contract.Price), cr.region, cr.contract.ContractId)
}

type contractReport struct {
	super    string
	region   int32
	contract esi.GetContractsPublicRegionId200Ok
	itemList string
}
