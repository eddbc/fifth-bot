package main

import (
	"github.com/gorilla/websocket"
	"log"
	"fmt"
	"time"
	"github.com/eddbc/fifth-bot/isk"
	"context"
)

var entitiesOfInterest = []int{
	1354830081, // goons
	99005338,	// horde
	//99008312, 	// escl8
	//99006411, 	// nsh
	98481691, 	// nogrl
}

func listenZKill() {
	url := "wss://zkillboard.com:2096"
	c, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	sub := &subscribe{"sub", "killstream"}

	err = c.WriteJSON(sub)
	if err != nil {
		log.Fatal("subscribe:", err)
	}

	log.Printf("listening to zkill")

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			kill := Kill{}
			err := c.ReadJSON(&kill)
			if err != nil {
				log.Println("read:", err)
				return
			}
			processKill(kill)
		}
	}()

	<-done
}

func processKill(kill Kill) {

	important := false
	msg := ""
	value := kill.Zkb.TotalValue

	// ignore all kills under 1M, to reduce spam
	if value < 1000000 {
		return
	}

	// ignore kills more that 24 hours old
	if kill.KillmailTime.Before( time.Now().Add(-1*(24*time.Hour)) ) {
		return
	}

	// ignore structures, maybe?
	if kill.Victim.CharacterID == 0 {
		return
	}

	// get filtering criteria
	isExpsv := isExpensive(kill)
	isKill, isLoss := isEntityRelated(kill)

	// ignore unrelated kills in highsec
	if !isKill && !isLoss {
		sys, _, err := eve.UniverseApi.GetUniverseSystemsSystemId(context.Background(), int32(kill.SolarSystemID), nil)
		if err == nil {
			log.Printf("system security for %v: %v",sys.Name, sys.SecurityStatus)
			if sys.SecurityStatus >= 0.5 {
				return
			}
		} else {
			log.Printf("error getting system: %v", err)
		}
	}

	// get kill details if matches any criteria, exit if not
	if isExpsv || isKill || isLoss {
		kill.inflate()
	} else {
		return
	}

	if isLoss { // ship lost by entity of interest
		msg = fmt.Sprintf("%v is a disgusting feeder", kill.Victim.CharacterName)
	} else if isKill { // ship killed by entity of interest
		fnlBlw, err := kill.getFinalBlow()
		if err != nil {
			fnlBlw.CharacterName = "Someone"
		}

		switch  {
		case value <=3000000: // Sub 3M kills
			msg = fmt.Sprintf("%v is kb padding. Disgusting.", fnlBlw.CharacterName)
			break
		case value <= 1000000000: // 3M - 1B kills
			msg = fmt.Sprintf("%v isn't completely useless.", fnlBlw.CharacterName)
			break
		default: // 1B+ kills
			msg = fmt.Sprintf("%v killed something big. Good job team.", fnlBlw.CharacterName)
			important = true
		}
	} else if isExpsv { // kill is expensive
		msg =  fmt.Sprintf("%v worth %v ISK died!", kill.Victim.ShipTypeName, isk.NearestThousandFormat(kill.Zkb.TotalValue))
	}

	// put zKill link in message
	if msg == "" {
		msg = kill.Zkb.url
	} else {
		msg = fmt.Sprintf("%v %v", msg, kill.Zkb.url)
	}

	// send message to appropriate channels
	if important {
		sendImportantMsg(msg)
	} else {
		sendMsg(msg)
	}

}

func isExpensive(km Kill) (bool){
	expsvLimit := float64(15000000000)
	return km.Zkb.TotalValue > expsvLimit
}

func isEntityRelated(km Kill) (bool, bool) {
	loss := false
	kill := false

	for _, id := range entitiesOfInterest {
		kill = km.isAttacker(id)
		loss = km.isVictim(id)
	}

	return kill, loss
}

type subscribe struct {
	Action  string `json:"action"`
	Channel string `json:"channel"`
}
