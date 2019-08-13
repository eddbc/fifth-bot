package fifth

import (
	"context"
	"errors"
	"fmt"
	"github.com/eddbc/fifth-bot/isk"
	"github.com/gorilla/websocket"
	"log"
	"math"
	"time"
)

var superTypes = []int32{
	30,  // Titan
	659, // Supercarrier
}

//var capTypes = []int32{
//	485,  // Dread
//	547,  // Carrier
//	1538, // FAX
//}

var entitiesOfInterest = []int32{

	// Corps
	98326834,   // NOG8S

	// Characters
	91872672,  // Edd Reynolds

	// Testing Groups
	//1354830081, // goons
	//99005338,	// horde
	//498125261, // test
}

//var stagingSystems = []int32{
//	30000974, // H-8F5Q
//}

//ListenZKill Start live zKill feed
func ListenZKill() {
	url := "wss://zkillboard.com:2096"
	c, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	sub := &subscribe{"sub", "killstream"}

	err = c.WriteJSON(sub)
	if err != nil {
		log.Fatal("Error starting zKillboard websocket: ", err)
	}

	log.Printf("zKillboard feed running")

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			kill := Kill{}
			err := c.ReadJSON(&kill)
			if err != nil {
				log.Fatal("Error, Restarting, zKillboard websocket closed: ", err)
			}
			processKill(&kill)
		}
	}()

	<-done

	defer log.Fatal("zKill websocket closed")
}

func logIgnore(reason string) {
	if Debug {
		//log.Printf("ignoring kill. reason: %v\n", reason)
	}
}

func processKill(kill *Kill) {

	important := false
	msg := ""
	value := kill.Zkb.TotalValue
	react := false

	// ignore all kills under 1M, to reduce spam
	if value < 1000000 {
		logIgnore("<1M ISK")
		return
	}

	// ignore kills more that 24 hours old
	if kill.KillmailTime.Before(time.Now().Add(-1 * (24 * time.Hour))) {
		logIgnore("too old")
		return
	}

	// ignore structures, maybe?
	if kill.Victim.CharacterID == 0 {
		logIgnore("no victim character")
		return
	}

	// get filtering criteria
	//isExpsv := isExpensive(kill)
	//isNearbyCap := isNearbyCap(kill)
	isKill, isLoss, err := isEntityRelated(kill)
	if err != nil {
		return
	}

	// ignore unrelated kills in highsec
	if !isKill && !isLoss {
		sys, _, err := Eve.UniverseApi.GetUniverseSystemsSystemId(context.Background(), int32(kill.SolarSystemID), nil)
		if err == nil {
			if sys.SecurityStatus >= 0.5 {
				logIgnore("high-sec")
				return
			}
		} else {
			log.Printf("error getting system: %v", err)
		}
	}

	// get kill details if matches any criteria, exit if not
	if isKill || isLoss {
		kill.inflate()
	} else {
		logIgnore("doesn't match criteria")
		return
	}

	if isLoss { // ship lost by entity of interest
		msg = fmt.Sprintf("**%v is a disgusting feeder**\n", kill.Victim.CharacterName)
		react = true
		if value > 100000000 {
			important = true
		}
	} else if isKill { // ship killed by entity of interest
		name, _ := kill.interestingName()

		if value >= 1000000000 {
			// 1B+ kills are important
			msg = fmt.Sprintf("%v killed something big. Good job team.\n", name)
			important = true
		} else {
			// <1B kills are not important
			msg = fmt.Sprintf("%v isn't completely useless.\n", name)
		}
	//} else if isExpsv { // kill is expensive
	//	important = true
	//} else if isNearbyCap {
	//	important = true
	//	msg = fmt.Sprintf("Capital activity in range!\n")
	}

	// put zKill link in message
	msg = fmt.Sprintf("%v%v [%v] %v ISK - %v (%v) <%v>",
		msg, kill.Victim.ShipTypeName, kill.Victim.ticker(),
		isk.NearestThousandFormat(kill.Zkb.TotalValue), kill.SolarSystemName,
		kill.RegionName, kill.Zkb.url)

	// send message to appropriate channels
	f := SendMsg
	if important {
		f = SendImportantMsg
	}

	m, err := f(msg)
	if err == nil && react {
		Session.MessageReactionAdd(m.ChannelID, m.ID, ":rip:486665154356969496")
	}
}

//func isNearbyCap(km *Kill) bool {
//	isCap := false
//	isNearby := false
//
//	for _, stagingId := range stagingSystems {
//		distance, err := distanceBetweenSystems(stagingId, km.SolarSystemID)
//		if err == nil {
//			if distance < 8 {
//				isNearby = true
//			}
//		}
//	}
//
//	if isNearby {
//		for _, attacker := range km.Attackers {
//			t, _, err := Eve.UniverseApi.GetUniverseTypesTypeId(ctx, attacker.ShipTypeID, nil)
//			if err != nil {
//				return false
//			}
//			for _, capId := range capTypes {
//				if t.GroupId == capId {
//					isCap = true
//				}
//			}
//			for _, superId := range superTypes {
//				if t.GroupId == superId {
//					isCap = true
//				}
//			}
//		}
//	}
//
//	return isCap && isNearby
//}

type xyz struct {
	x float64
	y float64
	z float64
}

func distanceBetweenSystems(idA int32, idB int32) (distance float64, err error) {

	a, _, err := Eve.UniverseApi.GetUniverseSystemsSystemId(ctx, idA, nil)
	if err != nil {
		return
	}

	b, _, err := Eve.UniverseApi.GetUniverseSystemsSystemId(ctx, idB, nil)
	if err != nil {
		return
	}

	distance = distanceBetweenPoints(
		xyz{a.Position.X, a.Position.Y, a.Position.Z},
		xyz{b.Position.X, b.Position.Y, b.Position.Z},
	)
	distance = distance / 9460730472580800      // convert from Meters to LY
	distance = math.Round(distance*1000) / 1000 // round to 3dp
	return
}

func distanceBetweenPoints(a xyz, b xyz) float64 {
	x := math.Pow(b.x-a.x, 2)
	y := math.Pow(b.y-a.y, 2)
	z := math.Pow(b.z-a.z, 2)
	return math.Pow(x+y+z, 0.5)
}

func isExpensive(km *Kill) bool {
	expsvLimit := float64(15000000000)
	return km.Zkb.TotalValue > expsvLimit
}

func isEntityRelated(km *Kill) (kill bool, loss bool, err error) {
	kill = false
	loss = false

	defer func() {
		if r := recover(); r != nil {
			err = errors.New("error getting related information")
			log.Printf("error getting related information for kill %v: %v+", km.KillmailID, r)
			//SendDebugMsg(fmt.Sprintf("error getting related information for kill: <%v>", km.getURL()))
		}
	}()

	for _, id := range entitiesOfInterest {
		kill = kill || km.isAttacker(id)
		loss = loss || km.isVictim(id)
	}

	return kill, loss, err
}

type subscribe struct {
	Action  string `json:"action"`
	Channel string `json:"channel"`
}
