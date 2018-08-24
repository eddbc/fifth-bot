package fifth

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

var zkillstats ZKillCharStats

type ZKillCharStats struct {
	client http.Client
	url    string
}

func init() {
	client := http.Client{
		Timeout: time.Second * 5, // Maximum of 5 secs
	}
	zkillstats = ZKillCharStats{client: client, url: "https://zkillboard.com/api/stats/characterID/%d/"}
}

func (z ZKillCharStats) get(id int32) (*ZKillCharacterStatsResp, error) {
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf(z.url, id), nil)
	req.Header.Set("User-Agent", useragent)

	res, getErr := z.client.Do(req)
	if getErr != nil {
		log.Printf("Error! %s\n", getErr)
		return nil, getErr
	}

	body, _ := ioutil.ReadAll(res.Body)
	stats := ZKillCharacterStatsResp{}
	jsonErr := json.Unmarshal(body, &stats)
	if jsonErr != nil {
		log.Printf("Error! %s\n", jsonErr)
		return nil, jsonErr
	}

	return &stats, nil
}

type ZKillCharacterStatsResp struct {
	AllTimeSum   int  `json:"allTimeSum"`
	CalcTrophies bool `json:"calcTrophies"`
	DangerRatio  int  `json:"dangerRatio"`
	GangRatio    int  `json:"gangRatio"`
	Groups       struct {
		Num25 struct {
			GroupID         int   `json:"groupID"`
			ShipsLost       int   `json:"shipsLost"`
			PointsLost      int   `json:"pointsLost"`
			IskLost         int   `json:"iskLost"`
			ShipsDestroyed  int   `json:"shipsDestroyed"`
			PointsDestroyed int   `json:"pointsDestroyed"`
			IskDestroyed    int64 `json:"iskDestroyed"`
		} `json:"25"`
		Num26 struct {
			GroupID         int   `json:"groupID"`
			ShipsLost       int   `json:"shipsLost"`
			PointsLost      int   `json:"pointsLost"`
			IskLost         int64 `json:"iskLost"`
			ShipsDestroyed  int   `json:"shipsDestroyed"`
			PointsDestroyed int   `json:"pointsDestroyed"`
			IskDestroyed    int64 `json:"iskDestroyed"`
		} `json:"26"`
		Num27 struct {
			GroupID         int   `json:"groupID"`
			ShipsLost       int   `json:"shipsLost"`
			PointsLost      int   `json:"pointsLost"`
			IskLost         int64 `json:"iskLost"`
			ShipsDestroyed  int   `json:"shipsDestroyed"`
			PointsDestroyed int   `json:"pointsDestroyed"`
			IskDestroyed    int64 `json:"iskDestroyed"`
		} `json:"27"`
		Num28 struct {
			GroupID         int `json:"groupID"`
			ShipsLost       int `json:"shipsLost"`
			PointsLost      int `json:"pointsLost"`
			IskLost         int `json:"iskLost"`
			ShipsDestroyed  int `json:"shipsDestroyed"`
			PointsDestroyed int `json:"pointsDestroyed"`
			IskDestroyed    int `json:"iskDestroyed"`
		} `json:"28"`
		Num29 struct {
			GroupID         int   `json:"groupID"`
			ShipsLost       int   `json:"shipsLost"`
			PointsLost      int   `json:"pointsLost"`
			IskLost         int64 `json:"iskLost"`
			ShipsDestroyed  int   `json:"shipsDestroyed"`
			PointsDestroyed int   `json:"pointsDestroyed"`
			IskDestroyed    int64 `json:"iskDestroyed"`
		} `json:"29"`
		Num30 struct {
			GroupID         int   `json:"groupID"`
			ShipsDestroyed  int   `json:"shipsDestroyed"`
			PointsDestroyed int   `json:"pointsDestroyed"`
			IskDestroyed    int64 `json:"iskDestroyed"`
		} `json:"30"`
		Num31 struct {
			GroupID         int `json:"groupID"`
			ShipsDestroyed  int `json:"shipsDestroyed"`
			PointsDestroyed int `json:"pointsDestroyed"`
			IskDestroyed    int `json:"iskDestroyed"`
		} `json:"31"`
		Num237 struct {
			GroupID         int `json:"groupID"`
			ShipsDestroyed  int `json:"shipsDestroyed"`
			PointsDestroyed int `json:"pointsDestroyed"`
			IskDestroyed    int `json:"iskDestroyed"`
		} `json:"237"`
		Num324 struct {
			GroupID         int `json:"groupID"`
			ShipsLost       int `json:"shipsLost"`
			PointsLost      int `json:"pointsLost"`
			IskLost         int `json:"iskLost"`
			ShipsDestroyed  int `json:"shipsDestroyed"`
			PointsDestroyed int `json:"pointsDestroyed"`
			IskDestroyed    int `json:"iskDestroyed"`
		} `json:"324"`
		Num358 struct {
			GroupID         int   `json:"groupID"`
			ShipsLost       int   `json:"shipsLost"`
			PointsLost      int   `json:"pointsLost"`
			IskLost         int   `json:"iskLost"`
			ShipsDestroyed  int   `json:"shipsDestroyed"`
			PointsDestroyed int   `json:"pointsDestroyed"`
			IskDestroyed    int64 `json:"iskDestroyed"`
		} `json:"358"`
		Num361 struct {
			GroupID         int `json:"groupID"`
			ShipsLost       int `json:"shipsLost"`
			PointsLost      int `json:"pointsLost"`
			IskLost         int `json:"iskLost"`
			ShipsDestroyed  int `json:"shipsDestroyed"`
			PointsDestroyed int `json:"pointsDestroyed"`
			IskDestroyed    int `json:"iskDestroyed"`
		} `json:"361"`
		Num365 struct {
			GroupID         int `json:"groupID"`
			ShipsDestroyed  int `json:"shipsDestroyed"`
			PointsDestroyed int `json:"pointsDestroyed"`
			IskDestroyed    int `json:"iskDestroyed"`
		} `json:"365"`
		Num380 struct {
			GroupID         int   `json:"groupID"`
			ShipsDestroyed  int   `json:"shipsDestroyed"`
			PointsDestroyed int   `json:"pointsDestroyed"`
			IskDestroyed    int64 `json:"iskDestroyed"`
		} `json:"380"`
		Num404 struct {
			GroupID         int `json:"groupID"`
			ShipsDestroyed  int `json:"shipsDestroyed"`
			PointsDestroyed int `json:"pointsDestroyed"`
			IskDestroyed    int `json:"iskDestroyed"`
		} `json:"404"`
		Num416 struct {
			GroupID         int `json:"groupID"`
			ShipsDestroyed  int `json:"shipsDestroyed"`
			PointsDestroyed int `json:"pointsDestroyed"`
			IskDestroyed    int `json:"iskDestroyed"`
		} `json:"416"`
		Num417 struct {
			GroupID         int `json:"groupID"`
			ShipsDestroyed  int `json:"shipsDestroyed"`
			PointsDestroyed int `json:"pointsDestroyed"`
			IskDestroyed    int `json:"iskDestroyed"`
		} `json:"417"`
		Num419 struct {
			GroupID         int   `json:"groupID"`
			ShipsLost       int   `json:"shipsLost"`
			PointsLost      int   `json:"pointsLost"`
			IskLost         int   `json:"iskLost"`
			ShipsDestroyed  int   `json:"shipsDestroyed"`
			PointsDestroyed int   `json:"pointsDestroyed"`
			IskDestroyed    int64 `json:"iskDestroyed"`
		} `json:"419"`
		Num420 struct {
			GroupID         int `json:"groupID"`
			ShipsLost       int `json:"shipsLost"`
			PointsLost      int `json:"pointsLost"`
			IskLost         int `json:"iskLost"`
			ShipsDestroyed  int `json:"shipsDestroyed"`
			PointsDestroyed int `json:"pointsDestroyed"`
			IskDestroyed    int `json:"iskDestroyed"`
		} `json:"420"`
		Num426 struct {
			GroupID         int `json:"groupID"`
			ShipsDestroyed  int `json:"shipsDestroyed"`
			PointsDestroyed int `json:"pointsDestroyed"`
			IskDestroyed    int `json:"iskDestroyed"`
		} `json:"426"`
		Num430 struct {
			GroupID         int `json:"groupID"`
			ShipsDestroyed  int `json:"shipsDestroyed"`
			PointsDestroyed int `json:"pointsDestroyed"`
			IskDestroyed    int `json:"iskDestroyed"`
		} `json:"430"`
		Num441 struct {
			GroupID         int `json:"groupID"`
			ShipsDestroyed  int `json:"shipsDestroyed"`
			PointsDestroyed int `json:"pointsDestroyed"`
			IskDestroyed    int `json:"iskDestroyed"`
		} `json:"441"`
		Num443 struct {
			GroupID         int `json:"groupID"`
			ShipsDestroyed  int `json:"shipsDestroyed"`
			PointsDestroyed int `json:"pointsDestroyed"`
			IskDestroyed    int `json:"iskDestroyed"`
		} `json:"443"`
		Num444 struct {
			GroupID         int `json:"groupID"`
			ShipsDestroyed  int `json:"shipsDestroyed"`
			PointsDestroyed int `json:"pointsDestroyed"`
			IskDestroyed    int `json:"iskDestroyed"`
		} `json:"444"`
		Num449 struct {
			GroupID         int `json:"groupID"`
			ShipsDestroyed  int `json:"shipsDestroyed"`
			PointsDestroyed int `json:"pointsDestroyed"`
			IskDestroyed    int `json:"iskDestroyed"`
		} `json:"449"`
		Num463 struct {
			GroupID         int `json:"groupID"`
			ShipsLost       int `json:"shipsLost"`
			PointsLost      int `json:"pointsLost"`
			IskLost         int `json:"iskLost"`
			ShipsDestroyed  int `json:"shipsDestroyed"`
			PointsDestroyed int `json:"pointsDestroyed"`
			IskDestroyed    int `json:"iskDestroyed"`
		} `json:"463"`
		Num485 struct {
			GroupID         int   `json:"groupID"`
			ShipsDestroyed  int   `json:"shipsDestroyed"`
			PointsDestroyed int   `json:"pointsDestroyed"`
			IskDestroyed    int64 `json:"iskDestroyed"`
		} `json:"485"`
		Num513 struct {
			GroupID         int `json:"groupID"`
			ShipsDestroyed  int `json:"shipsDestroyed"`
			PointsDestroyed int `json:"pointsDestroyed"`
			IskDestroyed    int `json:"iskDestroyed"`
		} `json:"513"`
		Num540 struct {
			GroupID         int   `json:"groupID"`
			ShipsLost       int   `json:"shipsLost"`
			PointsLost      int   `json:"pointsLost"`
			IskLost         int   `json:"iskLost"`
			ShipsDestroyed  int   `json:"shipsDestroyed"`
			PointsDestroyed int   `json:"pointsDestroyed"`
			IskDestroyed    int64 `json:"iskDestroyed"`
		} `json:"540"`
		Num541 struct {
			GroupID         int   `json:"groupID"`
			ShipsLost       int   `json:"shipsLost"`
			PointsLost      int   `json:"pointsLost"`
			IskLost         int   `json:"iskLost"`
			ShipsDestroyed  int   `json:"shipsDestroyed"`
			PointsDestroyed int   `json:"pointsDestroyed"`
			IskDestroyed    int64 `json:"iskDestroyed"`
		} `json:"541"`
		Num543 struct {
			GroupID         int   `json:"groupID"`
			ShipsDestroyed  int   `json:"shipsDestroyed"`
			PointsDestroyed int   `json:"pointsDestroyed"`
			IskDestroyed    int64 `json:"iskDestroyed"`
		} `json:"543"`
		Num547 struct {
			GroupID         int   `json:"groupID"`
			ShipsDestroyed  int   `json:"shipsDestroyed"`
			PointsDestroyed int   `json:"pointsDestroyed"`
			IskDestroyed    int64 `json:"iskDestroyed"`
		} `json:"547"`
		Num659 struct {
			GroupID         int   `json:"groupID"`
			ShipsDestroyed  int   `json:"shipsDestroyed"`
			PointsDestroyed int   `json:"pointsDestroyed"`
			IskDestroyed    int64 `json:"iskDestroyed"`
		} `json:"659"`
		Num707 struct {
			GroupID         int `json:"groupID"`
			ShipsDestroyed  int `json:"shipsDestroyed"`
			PointsDestroyed int `json:"pointsDestroyed"`
			IskDestroyed    int `json:"iskDestroyed"`
		} `json:"707"`
		Num830 struct {
			GroupID         int `json:"groupID"`
			ShipsDestroyed  int `json:"shipsDestroyed"`
			PointsDestroyed int `json:"pointsDestroyed"`
			IskDestroyed    int `json:"iskDestroyed"`
		} `json:"830"`
		Num831 struct {
			GroupID         int   `json:"groupID"`
			ShipsLost       int   `json:"shipsLost"`
			PointsLost      int   `json:"pointsLost"`
			IskLost         int   `json:"iskLost"`
			ShipsDestroyed  int   `json:"shipsDestroyed"`
			PointsDestroyed int   `json:"pointsDestroyed"`
			IskDestroyed    int64 `json:"iskDestroyed"`
		} `json:"831"`
		Num832 struct {
			GroupID         int   `json:"groupID"`
			ShipsLost       int   `json:"shipsLost"`
			PointsLost      int   `json:"pointsLost"`
			IskLost         int64 `json:"iskLost"`
			ShipsDestroyed  int   `json:"shipsDestroyed"`
			PointsDestroyed int   `json:"pointsDestroyed"`
			IskDestroyed    int64 `json:"iskDestroyed"`
		} `json:"832"`
		Num833 struct {
			GroupID         int   `json:"groupID"`
			ShipsDestroyed  int   `json:"shipsDestroyed"`
			PointsDestroyed int   `json:"pointsDestroyed"`
			IskDestroyed    int64 `json:"iskDestroyed"`
		} `json:"833"`
		Num834 struct {
			GroupID         int `json:"groupID"`
			ShipsLost       int `json:"shipsLost"`
			PointsLost      int `json:"pointsLost"`
			IskLost         int `json:"iskLost"`
			ShipsDestroyed  int `json:"shipsDestroyed"`
			PointsDestroyed int `json:"pointsDestroyed"`
			IskDestroyed    int `json:"iskDestroyed"`
		} `json:"834"`
		Num837 struct {
			GroupID         int `json:"groupID"`
			ShipsDestroyed  int `json:"shipsDestroyed"`
			PointsDestroyed int `json:"pointsDestroyed"`
			IskDestroyed    int `json:"iskDestroyed"`
		} `json:"837"`
		Num883 struct {
			GroupID         int   `json:"groupID"`
			ShipsDestroyed  int   `json:"shipsDestroyed"`
			PointsDestroyed int   `json:"pointsDestroyed"`
			IskDestroyed    int64 `json:"iskDestroyed"`
		} `json:"883"`
		Num893 struct {
			GroupID         int `json:"groupID"`
			ShipsLost       int `json:"shipsLost"`
			PointsLost      int `json:"pointsLost"`
			IskLost         int `json:"iskLost"`
			ShipsDestroyed  int `json:"shipsDestroyed"`
			PointsDestroyed int `json:"pointsDestroyed"`
			IskDestroyed    int `json:"iskDestroyed"`
		} `json:"893"`
		Num894 struct {
			GroupID         int `json:"groupID"`
			ShipsLost       int `json:"shipsLost"`
			PointsLost      int `json:"pointsLost"`
			IskLost         int `json:"iskLost"`
			ShipsDestroyed  int `json:"shipsDestroyed"`
			PointsDestroyed int `json:"pointsDestroyed"`
			IskDestroyed    int `json:"iskDestroyed"`
		} `json:"894"`
		Num898 struct {
			GroupID         int   `json:"groupID"`
			ShipsDestroyed  int   `json:"shipsDestroyed"`
			PointsDestroyed int   `json:"pointsDestroyed"`
			IskDestroyed    int64 `json:"iskDestroyed"`
		} `json:"898"`
		Num900 struct {
			GroupID         int   `json:"groupID"`
			ShipsDestroyed  int   `json:"shipsDestroyed"`
			PointsDestroyed int   `json:"pointsDestroyed"`
			IskDestroyed    int64 `json:"iskDestroyed"`
		} `json:"900"`
		Num902 struct {
			GroupID         int   `json:"groupID"`
			ShipsDestroyed  int   `json:"shipsDestroyed"`
			PointsDestroyed int   `json:"pointsDestroyed"`
			IskDestroyed    int64 `json:"iskDestroyed"`
		} `json:"902"`
		Num906 struct {
			GroupID         int   `json:"groupID"`
			ShipsDestroyed  int   `json:"shipsDestroyed"`
			PointsDestroyed int   `json:"pointsDestroyed"`
			IskDestroyed    int64 `json:"iskDestroyed"`
		} `json:"906"`
		Num941 struct {
			GroupID         int `json:"groupID"`
			ShipsDestroyed  int `json:"shipsDestroyed"`
			PointsDestroyed int `json:"pointsDestroyed"`
			IskDestroyed    int `json:"iskDestroyed"`
		} `json:"941"`
		Num963 struct {
			GroupID         int   `json:"groupID"`
			ShipsLost       int   `json:"shipsLost"`
			PointsLost      int   `json:"pointsLost"`
			IskLost         int64 `json:"iskLost"`
			ShipsDestroyed  int   `json:"shipsDestroyed"`
			PointsDestroyed int   `json:"pointsDestroyed"`
			IskDestroyed    int64 `json:"iskDestroyed"`
		} `json:"963"`
		Num1003 struct {
			GroupID         int `json:"groupID"`
			ShipsDestroyed  int `json:"shipsDestroyed"`
			PointsDestroyed int `json:"pointsDestroyed"`
			IskDestroyed    int `json:"iskDestroyed"`
		} `json:"1003"`
		Num1005 struct {
			GroupID         int `json:"groupID"`
			ShipsDestroyed  int `json:"shipsDestroyed"`
			PointsDestroyed int `json:"pointsDestroyed"`
			IskDestroyed    int `json:"iskDestroyed"`
		} `json:"1005"`
		Num1025 struct {
			GroupID         int `json:"groupID"`
			ShipsDestroyed  int `json:"shipsDestroyed"`
			PointsDestroyed int `json:"pointsDestroyed"`
			IskDestroyed    int `json:"iskDestroyed"`
		} `json:"1025"`
		Num1201 struct {
			GroupID         int `json:"groupID"`
			ShipsLost       int `json:"shipsLost"`
			PointsLost      int `json:"pointsLost"`
			IskLost         int `json:"iskLost"`
			ShipsDestroyed  int `json:"shipsDestroyed"`
			PointsDestroyed int `json:"pointsDestroyed"`
			IskDestroyed    int `json:"iskDestroyed"`
		} `json:"1201"`
		Num1202 struct {
			GroupID         int `json:"groupID"`
			ShipsDestroyed  int `json:"shipsDestroyed"`
			PointsDestroyed int `json:"pointsDestroyed"`
			IskDestroyed    int `json:"iskDestroyed"`
		} `json:"1202"`
		Num1246 struct {
			GroupID         int `json:"groupID"`
			ShipsDestroyed  int `json:"shipsDestroyed"`
			PointsDestroyed int `json:"pointsDestroyed"`
			IskDestroyed    int `json:"iskDestroyed"`
		} `json:"1246"`
		Num1250 struct {
			GroupID         int `json:"groupID"`
			ShipsLost       int `json:"shipsLost"`
			PointsLost      int `json:"pointsLost"`
			IskLost         int `json:"iskLost"`
			ShipsDestroyed  int `json:"shipsDestroyed"`
			PointsDestroyed int `json:"pointsDestroyed"`
			IskDestroyed    int `json:"iskDestroyed"`
		} `json:"1250"`
		Num1273 struct {
			GroupID         int `json:"groupID"`
			ShipsDestroyed  int `json:"shipsDestroyed"`
			PointsDestroyed int `json:"pointsDestroyed"`
			IskDestroyed    int `json:"iskDestroyed"`
		} `json:"1273"`
		Num1276 struct {
			GroupID         int `json:"groupID"`
			ShipsDestroyed  int `json:"shipsDestroyed"`
			PointsDestroyed int `json:"pointsDestroyed"`
			IskDestroyed    int `json:"iskDestroyed"`
		} `json:"1276"`
		Num1283 struct {
			GroupID         int `json:"groupID"`
			ShipsDestroyed  int `json:"shipsDestroyed"`
			PointsDestroyed int `json:"pointsDestroyed"`
			IskDestroyed    int `json:"iskDestroyed"`
		} `json:"1283"`
		Num1305 struct {
			GroupID         int   `json:"groupID"`
			ShipsLost       int   `json:"shipsLost"`
			PointsLost      int   `json:"pointsLost"`
			IskLost         int   `json:"iskLost"`
			ShipsDestroyed  int   `json:"shipsDestroyed"`
			PointsDestroyed int   `json:"pointsDestroyed"`
			IskDestroyed    int64 `json:"iskDestroyed"`
		} `json:"1305"`
		Num1404 struct {
			GroupID         int   `json:"groupID"`
			ShipsDestroyed  int   `json:"shipsDestroyed"`
			PointsDestroyed int   `json:"pointsDestroyed"`
			IskDestroyed    int64 `json:"iskDestroyed"`
		} `json:"1404"`
		Num1406 struct {
			GroupID         int   `json:"groupID"`
			ShipsDestroyed  int   `json:"shipsDestroyed"`
			PointsDestroyed int   `json:"pointsDestroyed"`
			IskDestroyed    int64 `json:"iskDestroyed"`
		} `json:"1406"`
		Num1527 struct {
			GroupID         int `json:"groupID"`
			ShipsDestroyed  int `json:"shipsDestroyed"`
			PointsDestroyed int `json:"pointsDestroyed"`
			IskDestroyed    int `json:"iskDestroyed"`
		} `json:"1527"`
		Num1534 struct {
			GroupID         int `json:"groupID"`
			ShipsLost       int `json:"shipsLost"`
			PointsLost      int `json:"pointsLost"`
			IskLost         int `json:"iskLost"`
			ShipsDestroyed  int `json:"shipsDestroyed"`
			PointsDestroyed int `json:"pointsDestroyed"`
			IskDestroyed    int `json:"iskDestroyed"`
		} `json:"1534"`
		Num1537 struct {
			GroupID         int `json:"groupID"`
			ShipsDestroyed  int `json:"shipsDestroyed"`
			PointsDestroyed int `json:"pointsDestroyed"`
			IskDestroyed    int `json:"iskDestroyed"`
		} `json:"1537"`
		Num1538 struct {
			GroupID         int   `json:"groupID"`
			ShipsDestroyed  int   `json:"shipsDestroyed"`
			PointsDestroyed int   `json:"pointsDestroyed"`
			IskDestroyed    int64 `json:"iskDestroyed"`
		} `json:"1538"`
		Num1652 struct {
			GroupID         int `json:"groupID"`
			ShipsDestroyed  int `json:"shipsDestroyed"`
			PointsDestroyed int `json:"pointsDestroyed"`
			IskDestroyed    int `json:"iskDestroyed"`
		} `json:"1652"`
		Num1653 struct {
			GroupID         int `json:"groupID"`
			ShipsDestroyed  int `json:"shipsDestroyed"`
			PointsDestroyed int `json:"pointsDestroyed"`
			IskDestroyed    int `json:"iskDestroyed"`
		} `json:"1653"`
		Num1657 struct {
			GroupID         int   `json:"groupID"`
			ShipsDestroyed  int   `json:"shipsDestroyed"`
			PointsDestroyed int   `json:"pointsDestroyed"`
			IskDestroyed    int64 `json:"iskDestroyed"`
		} `json:"1657"`
	} `json:"groups"`
	ID           int   `json:"id"`
	IskDestroyed int64 `json:"iskDestroyed"`
	IskLost      int64 `json:"iskLost"`
	Months       struct {
		Num201204 struct {
			Year       int `json:"year"`
			Month      int `json:"month"`
			ShipsLost  int `json:"shipsLost"`
			PointsLost int `json:"pointsLost"`
			IskLost    int `json:"iskLost"`
		} `json:"201204"`
		Num201207 struct {
			Year       int `json:"year"`
			Month      int `json:"month"`
			ShipsLost  int `json:"shipsLost"`
			PointsLost int `json:"pointsLost"`
			IskLost    int `json:"iskLost"`
		} `json:"201207"`
		Num201208 struct {
			Year       int `json:"year"`
			Month      int `json:"month"`
			ShipsLost  int `json:"shipsLost"`
			PointsLost int `json:"pointsLost"`
			IskLost    int `json:"iskLost"`
		} `json:"201208"`
		Num201209 struct {
			Year       int `json:"year"`
			Month      int `json:"month"`
			ShipsLost  int `json:"shipsLost"`
			PointsLost int `json:"pointsLost"`
			IskLost    int `json:"iskLost"`
		} `json:"201209"`
		Num201210 struct {
			Year       int `json:"year"`
			Month      int `json:"month"`
			ShipsLost  int `json:"shipsLost"`
			PointsLost int `json:"pointsLost"`
			IskLost    int `json:"iskLost"`
		} `json:"201210"`
		Num201212 struct {
			Year       int `json:"year"`
			Month      int `json:"month"`
			ShipsLost  int `json:"shipsLost"`
			PointsLost int `json:"pointsLost"`
			IskLost    int `json:"iskLost"`
		} `json:"201212"`
		Num201301 struct {
			Year            int `json:"year"`
			Month           int `json:"month"`
			ShipsLost       int `json:"shipsLost"`
			PointsLost      int `json:"pointsLost"`
			IskLost         int `json:"iskLost"`
			ShipsDestroyed  int `json:"shipsDestroyed"`
			PointsDestroyed int `json:"pointsDestroyed"`
			IskDestroyed    int `json:"iskDestroyed"`
		} `json:"201301"`
		Num201302 struct {
			Year            int `json:"year"`
			Month           int `json:"month"`
			ShipsLost       int `json:"shipsLost"`
			PointsLost      int `json:"pointsLost"`
			IskLost         int `json:"iskLost"`
			ShipsDestroyed  int `json:"shipsDestroyed"`
			PointsDestroyed int `json:"pointsDestroyed"`
			IskDestroyed    int `json:"iskDestroyed"`
		} `json:"201302"`
		Num201303 struct {
			Year       int `json:"year"`
			Month      int `json:"month"`
			ShipsLost  int `json:"shipsLost"`
			PointsLost int `json:"pointsLost"`
			IskLost    int `json:"iskLost"`
		} `json:"201303"`
		Num201304 struct {
			Year       int `json:"year"`
			Month      int `json:"month"`
			ShipsLost  int `json:"shipsLost"`
			PointsLost int `json:"pointsLost"`
			IskLost    int `json:"iskLost"`
		} `json:"201304"`
		Num201305 struct {
			Year       int `json:"year"`
			Month      int `json:"month"`
			ShipsLost  int `json:"shipsLost"`
			PointsLost int `json:"pointsLost"`
			IskLost    int `json:"iskLost"`
		} `json:"201305"`
		Num201312 struct {
			Year       int `json:"year"`
			Month      int `json:"month"`
			ShipsLost  int `json:"shipsLost"`
			PointsLost int `json:"pointsLost"`
			IskLost    int `json:"iskLost"`
		} `json:"201312"`
		Num201407 struct {
			Year       int `json:"year"`
			Month      int `json:"month"`
			ShipsLost  int `json:"shipsLost"`
			PointsLost int `json:"pointsLost"`
			IskLost    int `json:"iskLost"`
		} `json:"201407"`
		Num201411 struct {
			Year       int `json:"year"`
			Month      int `json:"month"`
			ShipsLost  int `json:"shipsLost"`
			PointsLost int `json:"pointsLost"`
			IskLost    int `json:"iskLost"`
		} `json:"201411"`
		Num201503 struct {
			Year            int `json:"year"`
			Month           int `json:"month"`
			ShipsLost       int `json:"shipsLost"`
			PointsLost      int `json:"pointsLost"`
			IskLost         int `json:"iskLost"`
			ShipsDestroyed  int `json:"shipsDestroyed"`
			PointsDestroyed int `json:"pointsDestroyed"`
			IskDestroyed    int `json:"iskDestroyed"`
		} `json:"201503"`
		Num201504 struct {
			Year            int   `json:"year"`
			Month           int   `json:"month"`
			ShipsLost       int   `json:"shipsLost"`
			PointsLost      int   `json:"pointsLost"`
			IskLost         int   `json:"iskLost"`
			ShipsDestroyed  int   `json:"shipsDestroyed"`
			PointsDestroyed int   `json:"pointsDestroyed"`
			IskDestroyed    int64 `json:"iskDestroyed"`
		} `json:"201504"`
		Num201505 struct {
			Year            int   `json:"year"`
			Month           int   `json:"month"`
			ShipsLost       int   `json:"shipsLost"`
			PointsLost      int   `json:"pointsLost"`
			IskLost         int   `json:"iskLost"`
			ShipsDestroyed  int   `json:"shipsDestroyed"`
			PointsDestroyed int   `json:"pointsDestroyed"`
			IskDestroyed    int64 `json:"iskDestroyed"`
		} `json:"201505"`
		Num201506 struct {
			Year            int   `json:"year"`
			Month           int   `json:"month"`
			ShipsLost       int   `json:"shipsLost"`
			PointsLost      int   `json:"pointsLost"`
			IskLost         int   `json:"iskLost"`
			ShipsDestroyed  int   `json:"shipsDestroyed"`
			PointsDestroyed int   `json:"pointsDestroyed"`
			IskDestroyed    int64 `json:"iskDestroyed"`
		} `json:"201506"`
		Num201507 struct {
			Year            int   `json:"year"`
			Month           int   `json:"month"`
			ShipsLost       int   `json:"shipsLost"`
			PointsLost      int   `json:"pointsLost"`
			IskLost         int   `json:"iskLost"`
			ShipsDestroyed  int   `json:"shipsDestroyed"`
			PointsDestroyed int   `json:"pointsDestroyed"`
			IskDestroyed    int64 `json:"iskDestroyed"`
		} `json:"201507"`
		Num201508 struct {
			Year            int   `json:"year"`
			Month           int   `json:"month"`
			ShipsLost       int   `json:"shipsLost"`
			PointsLost      int   `json:"pointsLost"`
			IskLost         int   `json:"iskLost"`
			ShipsDestroyed  int   `json:"shipsDestroyed"`
			PointsDestroyed int   `json:"pointsDestroyed"`
			IskDestroyed    int64 `json:"iskDestroyed"`
		} `json:"201508"`
		Num201509 struct {
			Year            int   `json:"year"`
			Month           int   `json:"month"`
			ShipsLost       int   `json:"shipsLost"`
			PointsLost      int   `json:"pointsLost"`
			IskLost         int   `json:"iskLost"`
			ShipsDestroyed  int   `json:"shipsDestroyed"`
			PointsDestroyed int   `json:"pointsDestroyed"`
			IskDestroyed    int64 `json:"iskDestroyed"`
		} `json:"201509"`
		Num201510 struct {
			Year            int   `json:"year"`
			Month           int   `json:"month"`
			ShipsLost       int   `json:"shipsLost"`
			PointsLost      int   `json:"pointsLost"`
			IskLost         int   `json:"iskLost"`
			ShipsDestroyed  int   `json:"shipsDestroyed"`
			PointsDestroyed int   `json:"pointsDestroyed"`
			IskDestroyed    int64 `json:"iskDestroyed"`
		} `json:"201510"`
		Num201511 struct {
			Year            int   `json:"year"`
			Month           int   `json:"month"`
			ShipsLost       int   `json:"shipsLost"`
			PointsLost      int   `json:"pointsLost"`
			IskLost         int   `json:"iskLost"`
			ShipsDestroyed  int   `json:"shipsDestroyed"`
			PointsDestroyed int   `json:"pointsDestroyed"`
			IskDestroyed    int64 `json:"iskDestroyed"`
		} `json:"201511"`
		Num201512 struct {
			Year            int   `json:"year"`
			Month           int   `json:"month"`
			ShipsLost       int   `json:"shipsLost"`
			PointsLost      int   `json:"pointsLost"`
			IskLost         int   `json:"iskLost"`
			ShipsDestroyed  int   `json:"shipsDestroyed"`
			PointsDestroyed int   `json:"pointsDestroyed"`
			IskDestroyed    int64 `json:"iskDestroyed"`
		} `json:"201512"`
		Num201601 struct {
			Year            int   `json:"year"`
			Month           int   `json:"month"`
			ShipsLost       int   `json:"shipsLost"`
			PointsLost      int   `json:"pointsLost"`
			IskLost         int   `json:"iskLost"`
			ShipsDestroyed  int   `json:"shipsDestroyed"`
			PointsDestroyed int   `json:"pointsDestroyed"`
			IskDestroyed    int64 `json:"iskDestroyed"`
		} `json:"201601"`
		Num201602 struct {
			Year            int   `json:"year"`
			Month           int   `json:"month"`
			ShipsLost       int   `json:"shipsLost"`
			PointsLost      int   `json:"pointsLost"`
			IskLost         int   `json:"iskLost"`
			ShipsDestroyed  int   `json:"shipsDestroyed"`
			PointsDestroyed int   `json:"pointsDestroyed"`
			IskDestroyed    int64 `json:"iskDestroyed"`
		} `json:"201602"`
		Num201603 struct {
			Year            int   `json:"year"`
			Month           int   `json:"month"`
			ShipsLost       int   `json:"shipsLost"`
			PointsLost      int   `json:"pointsLost"`
			IskLost         int   `json:"iskLost"`
			ShipsDestroyed  int   `json:"shipsDestroyed"`
			PointsDestroyed int   `json:"pointsDestroyed"`
			IskDestroyed    int64 `json:"iskDestroyed"`
		} `json:"201603"`
		Num201604 struct {
			Year            int   `json:"year"`
			Month           int   `json:"month"`
			ShipsLost       int   `json:"shipsLost"`
			PointsLost      int   `json:"pointsLost"`
			IskLost         int   `json:"iskLost"`
			ShipsDestroyed  int   `json:"shipsDestroyed"`
			PointsDestroyed int   `json:"pointsDestroyed"`
			IskDestroyed    int64 `json:"iskDestroyed"`
		} `json:"201604"`
		Num201605 struct {
			Year            int   `json:"year"`
			Month           int   `json:"month"`
			ShipsLost       int   `json:"shipsLost"`
			PointsLost      int   `json:"pointsLost"`
			IskLost         int   `json:"iskLost"`
			ShipsDestroyed  int   `json:"shipsDestroyed"`
			PointsDestroyed int   `json:"pointsDestroyed"`
			IskDestroyed    int64 `json:"iskDestroyed"`
		} `json:"201605"`
		Num201606 struct {
			Year            int   `json:"year"`
			Month           int   `json:"month"`
			ShipsLost       int   `json:"shipsLost"`
			PointsLost      int   `json:"pointsLost"`
			IskLost         int   `json:"iskLost"`
			ShipsDestroyed  int   `json:"shipsDestroyed"`
			PointsDestroyed int   `json:"pointsDestroyed"`
			IskDestroyed    int64 `json:"iskDestroyed"`
		} `json:"201606"`
		Num201607 struct {
			Year            int   `json:"year"`
			Month           int   `json:"month"`
			ShipsLost       int   `json:"shipsLost"`
			PointsLost      int   `json:"pointsLost"`
			IskLost         int   `json:"iskLost"`
			ShipsDestroyed  int   `json:"shipsDestroyed"`
			PointsDestroyed int   `json:"pointsDestroyed"`
			IskDestroyed    int64 `json:"iskDestroyed"`
		} `json:"201607"`
		Num201608 struct {
			Year            int   `json:"year"`
			Month           int   `json:"month"`
			ShipsLost       int   `json:"shipsLost"`
			PointsLost      int   `json:"pointsLost"`
			IskLost         int   `json:"iskLost"`
			ShipsDestroyed  int   `json:"shipsDestroyed"`
			PointsDestroyed int   `json:"pointsDestroyed"`
			IskDestroyed    int64 `json:"iskDestroyed"`
		} `json:"201608"`
		Num201609 struct {
			Year            int   `json:"year"`
			Month           int   `json:"month"`
			ShipsLost       int   `json:"shipsLost"`
			PointsLost      int   `json:"pointsLost"`
			IskLost         int   `json:"iskLost"`
			ShipsDestroyed  int   `json:"shipsDestroyed"`
			PointsDestroyed int   `json:"pointsDestroyed"`
			IskDestroyed    int64 `json:"iskDestroyed"`
		} `json:"201609"`
		Num201610 struct {
			Year            int   `json:"year"`
			Month           int   `json:"month"`
			ShipsLost       int   `json:"shipsLost"`
			PointsLost      int   `json:"pointsLost"`
			IskLost         int   `json:"iskLost"`
			ShipsDestroyed  int   `json:"shipsDestroyed"`
			PointsDestroyed int   `json:"pointsDestroyed"`
			IskDestroyed    int64 `json:"iskDestroyed"`
		} `json:"201610"`
		Num201611 struct {
			Year            int   `json:"year"`
			Month           int   `json:"month"`
			ShipsLost       int   `json:"shipsLost"`
			PointsLost      int   `json:"pointsLost"`
			IskLost         int   `json:"iskLost"`
			ShipsDestroyed  int   `json:"shipsDestroyed"`
			PointsDestroyed int   `json:"pointsDestroyed"`
			IskDestroyed    int64 `json:"iskDestroyed"`
		} `json:"201611"`
		Num201612 struct {
			Year            int   `json:"year"`
			Month           int   `json:"month"`
			ShipsLost       int   `json:"shipsLost"`
			PointsLost      int   `json:"pointsLost"`
			IskLost         int   `json:"iskLost"`
			ShipsDestroyed  int   `json:"shipsDestroyed"`
			PointsDestroyed int   `json:"pointsDestroyed"`
			IskDestroyed    int64 `json:"iskDestroyed"`
		} `json:"201612"`
		Num201701 struct {
			Year            int   `json:"year"`
			Month           int   `json:"month"`
			ShipsLost       int   `json:"shipsLost"`
			PointsLost      int   `json:"pointsLost"`
			IskLost         int64 `json:"iskLost"`
			ShipsDestroyed  int   `json:"shipsDestroyed"`
			PointsDestroyed int   `json:"pointsDestroyed"`
			IskDestroyed    int64 `json:"iskDestroyed"`
		} `json:"201701"`
		Num201702 struct {
			Year            int   `json:"year"`
			Month           int   `json:"month"`
			ShipsDestroyed  int   `json:"shipsDestroyed"`
			PointsDestroyed int   `json:"pointsDestroyed"`
			IskDestroyed    int64 `json:"iskDestroyed"`
		} `json:"201702"`
		Num201703 struct {
			Year            int   `json:"year"`
			Month           int   `json:"month"`
			ShipsLost       int   `json:"shipsLost"`
			PointsLost      int   `json:"pointsLost"`
			IskLost         int   `json:"iskLost"`
			ShipsDestroyed  int   `json:"shipsDestroyed"`
			PointsDestroyed int   `json:"pointsDestroyed"`
			IskDestroyed    int64 `json:"iskDestroyed"`
		} `json:"201703"`
		Num201704 struct {
			Year            int   `json:"year"`
			Month           int   `json:"month"`
			ShipsLost       int   `json:"shipsLost"`
			PointsLost      int   `json:"pointsLost"`
			IskLost         int64 `json:"iskLost"`
			ShipsDestroyed  int   `json:"shipsDestroyed"`
			PointsDestroyed int   `json:"pointsDestroyed"`
			IskDestroyed    int64 `json:"iskDestroyed"`
		} `json:"201704"`
		Num201705 struct {
			Year            int   `json:"year"`
			Month           int   `json:"month"`
			ShipsLost       int   `json:"shipsLost"`
			PointsLost      int   `json:"pointsLost"`
			IskLost         int   `json:"iskLost"`
			ShipsDestroyed  int   `json:"shipsDestroyed"`
			PointsDestroyed int   `json:"pointsDestroyed"`
			IskDestroyed    int64 `json:"iskDestroyed"`
		} `json:"201705"`
		Num201706 struct {
			Year            int   `json:"year"`
			Month           int   `json:"month"`
			ShipsLost       int   `json:"shipsLost"`
			PointsLost      int   `json:"pointsLost"`
			IskLost         int   `json:"iskLost"`
			ShipsDestroyed  int   `json:"shipsDestroyed"`
			PointsDestroyed int   `json:"pointsDestroyed"`
			IskDestroyed    int64 `json:"iskDestroyed"`
		} `json:"201706"`
		Num201707 struct {
			Year            int `json:"year"`
			Month           int `json:"month"`
			ShipsLost       int `json:"shipsLost"`
			PointsLost      int `json:"pointsLost"`
			IskLost         int `json:"iskLost"`
			ShipsDestroyed  int `json:"shipsDestroyed"`
			PointsDestroyed int `json:"pointsDestroyed"`
			IskDestroyed    int `json:"iskDestroyed"`
		} `json:"201707"`
		Num201708 struct {
			Year            int `json:"year"`
			Month           int `json:"month"`
			ShipsLost       int `json:"shipsLost"`
			PointsLost      int `json:"pointsLost"`
			IskLost         int `json:"iskLost"`
			ShipsDestroyed  int `json:"shipsDestroyed"`
			PointsDestroyed int `json:"pointsDestroyed"`
			IskDestroyed    int `json:"iskDestroyed"`
		} `json:"201708"`
		Num201710 struct {
			Year            int   `json:"year"`
			Month           int   `json:"month"`
			ShipsLost       int   `json:"shipsLost"`
			PointsLost      int   `json:"pointsLost"`
			IskLost         int   `json:"iskLost"`
			ShipsDestroyed  int   `json:"shipsDestroyed"`
			PointsDestroyed int   `json:"pointsDestroyed"`
			IskDestroyed    int64 `json:"iskDestroyed"`
		} `json:"201710"`
		Num201711 struct {
			Year            int   `json:"year"`
			Month           int   `json:"month"`
			ShipsLost       int   `json:"shipsLost"`
			PointsLost      int   `json:"pointsLost"`
			IskLost         int   `json:"iskLost"`
			ShipsDestroyed  int   `json:"shipsDestroyed"`
			PointsDestroyed int   `json:"pointsDestroyed"`
			IskDestroyed    int64 `json:"iskDestroyed"`
		} `json:"201711"`
		Num201712 struct {
			Year            int   `json:"year"`
			Month           int   `json:"month"`
			ShipsLost       int   `json:"shipsLost"`
			PointsLost      int   `json:"pointsLost"`
			IskLost         int   `json:"iskLost"`
			ShipsDestroyed  int   `json:"shipsDestroyed"`
			PointsDestroyed int   `json:"pointsDestroyed"`
			IskDestroyed    int64 `json:"iskDestroyed"`
		} `json:"201712"`
		Num201801 struct {
			Year            int `json:"year"`
			Month           int `json:"month"`
			ShipsDestroyed  int `json:"shipsDestroyed"`
			PointsDestroyed int `json:"pointsDestroyed"`
			IskDestroyed    int `json:"iskDestroyed"`
		} `json:"201801"`
		Num201802 struct {
			Year            int   `json:"year"`
			Month           int   `json:"month"`
			ShipsDestroyed  int   `json:"shipsDestroyed"`
			PointsDestroyed int   `json:"pointsDestroyed"`
			IskDestroyed    int64 `json:"iskDestroyed"`
		} `json:"201802"`
		Num201803 struct {
			Year            int `json:"year"`
			Month           int `json:"month"`
			ShipsLost       int `json:"shipsLost"`
			PointsLost      int `json:"pointsLost"`
			IskLost         int `json:"iskLost"`
			ShipsDestroyed  int `json:"shipsDestroyed"`
			PointsDestroyed int `json:"pointsDestroyed"`
			IskDestroyed    int `json:"iskDestroyed"`
		} `json:"201803"`
		Num201804 struct {
			Year            int   `json:"year"`
			Month           int   `json:"month"`
			ShipsLost       int   `json:"shipsLost"`
			PointsLost      int   `json:"pointsLost"`
			IskLost         int   `json:"iskLost"`
			ShipsDestroyed  int   `json:"shipsDestroyed"`
			PointsDestroyed int   `json:"pointsDestroyed"`
			IskDestroyed    int64 `json:"iskDestroyed"`
		} `json:"201804"`
		Num201805 struct {
			Year            int   `json:"year"`
			Month           int   `json:"month"`
			ShipsDestroyed  int   `json:"shipsDestroyed"`
			PointsDestroyed int   `json:"pointsDestroyed"`
			IskDestroyed    int64 `json:"iskDestroyed"`
		} `json:"201805"`
		Num201806 struct {
			Year            int `json:"year"`
			Month           int `json:"month"`
			ShipsDestroyed  int `json:"shipsDestroyed"`
			PointsDestroyed int `json:"pointsDestroyed"`
			IskDestroyed    int `json:"iskDestroyed"`
		} `json:"201806"`
	} `json:"months"`
	NextTopRecalc   int `json:"nextTopRecalc"`
	PointsDestroyed int `json:"pointsDestroyed"`
	PointsLost      int `json:"pointsLost"`
	Sequence        int `json:"sequence"`
	ShipsDestroyed  int `json:"shipsDestroyed"`
	ShipsLost       int `json:"shipsLost"`
	SoloKills       int `json:"soloKills"`
	SoloLosses      int `json:"soloLosses"`
	TopAllTime      []struct {
		Type string `json:"type"`
		Data []struct {
			Kills int `json:"kills"`
			//CharacterID   int    `json:"characterID"`
			//CorporationID int    `json:"corporationID"`
			ShipTypeId int `json:"shipTypeID"`
		} `json:"data"`
	} `json:"topAllTime"`
	Trophies struct {
		Levels int `json:"levels"`
		Max    int `json:"max"`
	} `json:"trophies"`
	Type      string `json:"type"`
	Activepvp struct {
		Kills struct {
			Type  string `json:"type"`
			Count int    `json:"count"`
		} `json:"kills"`
	} `json:"activepvp"`
	Info struct {
		AllianceID    int `json:"allianceID"`
		CorporationID int `json:"corporationID"`
		FactionID     int `json:"factionID"`
		ID            int `json:"id"`
		KillID        int `json:"killID"`
		LastAPIUpdate struct {
			Sec  int `json:"sec"`
			Usec int `json:"usec"`
		} `json:"lastApiUpdate"`
		Name      string  `json:"name"`
		SecStatus float64 `json:"secStatus"`
		Skip      int     `json:"skip"`
		Type      string  `json:"type"`
	} `json:"info"`
	TopLists []struct {
		Type   string `json:"type"`
		Title  string `json:"title"`
		Values []struct {
			Kills         int         `json:"kills"`
			CharacterID   int         `json:"characterID"`
			CharacterName string      `json:"characterName"`
			ID            int         `json:"id"`
			TypeID        interface{} `json:"typeID"`
			Name          string      `json:"name"`
		} `json:"values"`
	} `json:"topLists"`
	TopIskKillIDs []int `json:"topIskKillIDs"`
}
