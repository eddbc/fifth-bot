package main

import (
	"github.com/gorilla/websocket"
	"log"
	"fmt"
	"context"
	"time"
	"github.com/eddbc/fifth-bot/isk"
)

var entitiesOfInterest = []int{
	//1354830081, // goons
	//99005338,	// horde
	//99008312, 	// escl8
	//99006411, 	// nsh
	98481691, 	// nogrl
}

var killPostChannels = []string{
	"459341365562572803", // test lab
	"459997163787649025", // nogrl bot spam
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

	interesting := false
	msg := ""

	// ignore all kills under 1M, to reduce spam
	if kill.Zkb.TotalValue < 1000000 {
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

	// look for expensive kills (+15B)
	if kill.Zkb.TotalValue > 15000000000 {
		interesting = true
		kill.inflate()
		msg =  fmt.Sprintf("%v worth %v ISK died!", kill.Victim.ShipTypeName, isk.NearestThousandFormat(kill.Zkb.TotalValue))
	}

	// look for entities of interest on the kill
	for _, id := range entitiesOfInterest {
		if kill.isAttacker(id){
			interesting = true
			kill.inflate()
			list := []int32{int32(id)}
			res, _, err := eve.UniverseApi.PostUniverseNames(context.Background(), list, nil)
			if err != nil {
				log.Fatal("entity lookup:", err)
			}
			e := res[0]

			killer := e.Name

			fb, err := kill.getFinalBlow()
			if err != nil {
				killer = fb.CharacterName
			}

			msg = fmt.Sprintf("%v managed to kill something. Good job, I guess.", killer)
		}

		if kill.isVictim(id){
			interesting = true
			kill.inflate()
			msg = fmt.Sprintf("%v is a disgusting feeder.", kill.Victim.CharacterName)
		}
	}

	if interesting {
		kill.inflate()
		msg = fmt.Sprintf("%v %v", msg, kill.Zkb.url)
		log.Println(msg)
		for _, ch := range killPostChannels {
			Session.ChannelMessageSend(ch, msg)
		}
	}
}

type subscribe struct {
	Action  string `json:"action"`
	Channel string `json:"channel"`
}
