package fifth

import (
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/eddbc/fifth-bot/mux"
	"github.com/eddbc/fifth-bot/storage"
	"go.etcd.io/bbolt"
	"log"
	"sort"
	"strconv"
	"time"
)

type timer struct {
	ID          int       `json:"id"`
	Time        time.Time `json:"time"`
	Description string    `json:"description"`
	Pinged      bool      `json:"pinged"`
}

func (t *timer) toStr() string {
	loc, _ := time.LoadLocation("Atlantic/Reykjavik")
	_, _, d, h, m, _ := diff(time.Now().In(loc), t.Time)
	return fmt.Sprintf("%v - %v (%vd %vh %vm) [%v]", t.Description, t.Time.Format("Jan 2, 15:04"), d, h, m, t.ID)
}

//AddTimer Bot command to add a timer to the timer board
func (f *Fifth) AddTimer(ds *discordgo.Session, dm *discordgo.Message, ctx *mux.Context) {
	desc := ""
	d := 0
	h := 0
	m := 0
	var err error
	for k, v := range ctx.Fields {
		if k == 1 {
			d, err = removeTrailingChar(v)
		}
		if k == 2 {
			h, err = removeTrailingChar(v)
		}
		if k == 3 {
			m, err = removeTrailingChar(v)
		}
		if k > 3 {
			desc += v
			if k+1 < len(ctx.Fields) {
				desc += " "
			}
		}

		if err != nil {
			log.Println("Invalid Timer")
			ds.ChannelMessageSend(dm.ChannelID, "Invalid Time Given")
			return
		}
	}

	if d > 16 || h > 23 || m > 59 {
		log.Println("Invalid Timer")
		ds.ChannelMessageSend(dm.ChannelID, "Invalid Time Given")
		return
	}

	if desc == "" {
		log.Println("Invalid Timer")
		ds.ChannelMessageSend(dm.ChannelID, "Invalid Description Given")
		return
	}

	loc, _ := time.LoadLocation("Atlantic/Reykjavik")

	then := time.Now().In(loc).Add(
		time.Hour*time.Duration(24*d) +
			time.Hour*time.Duration(h) +
			time.Minute*time.Duration(m))

	s := fmt.Sprintf("Set timer for %v with description \"%s\"", then.Format("Jan 2, 15:04"), desc)
	log.Println(s)
	ds.ChannelMessageSend(dm.ChannelID, s)

	saveTimer(timer{Time: then, Description: desc, Pinged: false})
}

//ListTimers bot command to return a list of current timers
func (f *Fifth) ListTimers(ds *discordgo.Session, dm *discordgo.Message, ctx *mux.Context) {
	var timers []timer
	storage.DB.View(func(tx *bbolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte(storage.TimersKey))

		b.ForEach(func(k, v []byte) error {
			r := timer{}
			json.Unmarshal(v, &r)
			timers = append(timers, r)
			return nil
		})
		return nil
	})

	sort.Slice(timers, func(i, j int) bool {
		return timers[i].Time.Before(timers[j].Time)
	})

	resp := "```\n"
	if len(timers) == 0 {
		resp += "No Timers"
	}
	for _, timer := range timers {
		resp += timer.toStr() + "\n\n"
	}
	resp += "```"
	log.Println(resp)
	SendMsgToChan(dm.ChannelID, resp)
	//ds.ChannelMessageSend(dm.ChannelID, resp)
}

// RemoveTimer bot command to remove a timer from the timer board
func (f *Fifth) RemoveTimer(ds *discordgo.Session, dm *discordgo.Message, ctx *mux.Context) {

	i, err := strconv.Atoi(ctx.Fields[1])
	SendMsgToChan(dm.ChannelID, fmt.Sprintf("Deleting timer [%v]...", i))

	if err != nil {
		return
	}

	storage.DB.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(storage.TimersKey))
		b.Delete(storage.Itob(i))
		return nil
	})
}

func timerCron() {
	storage.DB.Update(func(tx *bbolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte(storage.TimersKey))

		b.ForEach(func(k, v []byte) error {
			timer := timer{}
			json.Unmarshal(v, &timer)

			if !timer.Pinged {
				then := time.Now().Add(30 * time.Minute)
				if timer.Time.Before(then) {
					SendMsgToChan(timerChannel, "@here Timer Warning : "+timer.toStr())
					timer.Pinged = true
					bytes, err := json.Marshal(timer)
					if err == nil {
						b.Put(k, bytes)
					}
				}
			}

			if timer.Time.Before(time.Now()) {
				log.Printf("Timer has expired, deleting : < %v >", timer.toStr())
				b.Delete(k)
			}
			return nil
		})
		return nil
	})
}

func saveTimer(t timer) {
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

func removeTrailingChar(s string) (int, error) {
	if last := len(s) - 1; last >= 0 {
		s = s[:last]
	}
	i, err := strconv.Atoi(s)
	return i, err
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
