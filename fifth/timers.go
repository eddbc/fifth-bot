package fifth

import (
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/eddbc/fifth-bot/mux"
	"github.com/eddbc/fifth-bot/storage"
	"go.etcd.io/bbolt"
	"log"
	"strconv"
	"time"
)

type Timer struct {
	ID int `json:"id"`
	Time  time.Time `json:"time"`
	Description string `json:"description"`
}

func (f *Fifth) AddTimer(ds *discordgo.Session, dm *discordgo.Message, ctx *mux.Context) {
	desc := ""
	d:=0
	h:=0
	m:=0
	for k, v := range ctx.Fields {
		if k == 1 {
			d = removeTrailingChar(v)
		}
		if k == 2 {
			h = removeTrailingChar(v)
		}
		if k == 3 {
			m = removeTrailingChar(v)
		}
		if k > 3 {
			desc+=v
			if k+1 < len(ctx.Fields) {
				desc += " "
			}
		}
	}

	if d > 8 || h > 23 || m > 59 {
		log.Println("Invalid Timer")
	}
	loc, _ := time.LoadLocation("Atlantic/Reykjavik")

	then := time.Now().In(loc).Add(
		time.Hour * time.Duration(24*d) +
		time.Hour* time.Duration(h) +
		time.Minute * time.Duration(m))

	s := fmt.Sprintf("Set timer for %v with description \"%s\"", then.Format("Jan 2, 15:04"), desc)
	log.Println(s)
	ds.ChannelMessageSend(dm.ChannelID, s)

	saveTimer(Timer{Time:then, Description:desc})
}

func (f *Fifth) ListTimers(ds *discordgo.Session, dm *discordgo.Message, ctx *mux.Context) {
	var timers []Timer
	storage.DB.View(func(tx *bbolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte(storage.TimersKey))

		b.ForEach(func(k, v []byte) error {
			r := Timer{}
			json.Unmarshal(v, &r)
			timers = append(timers, r)
			return nil
		})
		return nil
	})

	resp := "```\n"
	loc, _ := time.LoadLocation("Atlantic/Reykjavik")
	if len(timers) == 0 {
		resp += "No Timers"
	}
	for _, timer := range timers {
		_, _, d, h, m, _ := diff(time.Now().In(loc), timer.Time)
		resp += fmt.Sprintf("%v (%vd %vh %vm) [%v]\n",timer.Description, d, h, m, timer.ID)
	}
	resp += "```"
	log.Println(resp)
	ds.ChannelMessageSend(dm.ChannelID, resp)
}

func saveTimer(t Timer) {
	storage.DB.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(storage.TimersKey))
		if b == nil {
			return nil
		}

		id, _ := b.NextSequence()
		t.ID = int(id)

		buf, _ := json.Marshal(t)

		b.Put(storage.Itob(t.ID), buf)
		return nil
	})
}

func removeTrailingChar(s string) int {
	if last := len(s) - 1; last >= 0 {
		s = s[:last]
	}
	i, _ := strconv.Atoi(s)
	return i
}

func diff(a, b time.Time) (year, month, day, hour, min, sec int) {
	if a.Location() != b.Location() {
		b = b.In(a.Location())
	}
	if a.After(b) {
		a, b = b, a
	}
	y1, M1, d1 := a.Date()
	y2, M2, d2 := b.Date()

	h1, m1, s1 := a.Clock()
	h2, m2, s2 := b.Clock()

	year = int(y2 - y1)
	month = int(M2 - M1)
	day = int(d2 - d1)
	hour = int(h2 - h1)
	min = int(m2 - m1)
	sec = int(s2 - s1)

	// Normalize negative values
	if sec < 0 {
		sec += 60
		min--
	}
	if min < 0 {
		min += 60
		hour--
	}
	if hour < 0 {
		hour += 24
		day--
	}
	if day < 0 {
		// days in month:
		t := time.Date(y1, M1, 32, 0, 0, 0, 0, time.UTC)
		day += 32 - t.Day()
		month--
	}
	if month < 0 {
		month += 12
		year--
	}

	return
}
