package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/eddbc/fifth-bot/isk"
	"github.com/gorilla/websocket"
	"log"
	"time"
)

var entitiesOfInterest = []int{
	1354830081, // goons
	99005338,	// horde
	98481691, // nogrl
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
				log.Fatal("websocket read: ", err)
			}
			processKill(kill)
		}
	}()

	<-done
}

func logIgnore(reason string) {
	if debug {
		//log.Printf("ignoring kill. reason: %v\n", reason)
	}
}

func processKill(kill Kill) {

	important := false
	msg := ""
	value := kill.Zkb.TotalValue

	// ignore all kills under 5M, to reduce spam
	if value < 5000000 {
		logIgnore("<5M ISK")
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
	isExpsv := isExpensive(kill)
	isKill, isLoss, err := isEntityRelated(kill)
	if err != nil {
		return
	}

	// ignore unrelated kills in highsec
	if !isKill && !isLoss {
		sys, _, err := eve.UniverseApi.GetUniverseSystemsSystemId(context.Background(), int32(kill.SolarSystemID), nil)
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
	if isExpsv || isKill || isLoss {
		kill.inflate()
	} else {
		logIgnore("doesn't match criteria")
		return
	}

	if isLoss { // ship lost by entity of interest
		msg = fmt.Sprintf("%v is a disgusting feeder", kill.Victim.CharacterName)
		if value > 500000000 {
			important = true
		}
	} else if isKill { // ship killed by entity of interest
		name, _ := kill.interestingName()

		if value >= 1000000000 {
			// 1B+ kills are important
			msg = fmt.Sprintf("%v killed something big. Good job team.", name)
			important = true
		} else {
			// <1B kills are not important
			msg = fmt.Sprintf("%v isn't completely useless.", name)
		}
	} else if isExpsv { // kill is expensive
		msg = fmt.Sprintf("%v worth %v ISK died!", kill.Victim.ShipTypeName, isk.NearestThousandFormat(kill.Zkb.TotalValue))
	}

	// put zKill link in message
	if msg == "" {
		msg = kill.Zkb.url
	} else {
		msg = fmt.Sprintf("%v <%v>", msg, kill.Zkb.url)
	}

	// send message to appropriate channels
	if important {
		sendImportantMsg(msg)
	} else {
		sendMsg(msg)
	}

}

func isExpensive(km Kill) bool {
	expsvLimit := float64(15000000000)
	return km.Zkb.TotalValue > expsvLimit
}

func isEntityRelated(km Kill) (kill bool, loss bool, err error) {
	kill = false
	loss = false

	defer func() {
		if r := recover(); r != nil {
			err = errors.New("error getting related information")
			log.Printf("error getting related information for kill %v: %v+", km.KillmailID, r)
			sendDebugMsg(fmt.Sprintf("error: crashing because of this kill: %v", km.getUrl()))
		}
	}()

	for _, id := range entitiesOfInterest {
		kill = km.isAttacker(id)
		loss = km.isVictim(id)
	}

	return kill, loss, nil
}

type subscribe struct {
	Action  string `json:"action"`
	Channel string `json:"channel"`
}
