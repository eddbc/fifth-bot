package fifth

import (
	"github.com/bwmarrin/discordgo"
	"log"
)

var debugChannel = "459341365562572803"     // testing-lab
var importantChannel = "486808823110303754" // general
var spamChannel = "488473067505319938"      // kill-feed
var timerChannel = "517892609155399683"

func SendMsg(msg string) (*discordgo.Message, error) {
	return SendMsgToChan(spamChannel, msg)
}

func SendImportantMsg(msg string) (*discordgo.Message, error) {
	return SendMsgToChan(importantChannel, msg)
}

func SendDebugMsg(msg string) (*discordgo.Message, error) {
	return Session.ChannelMessageSend(debugChannel, msg)
}

func SendMsgToChan(chann string, msg string) (*discordgo.Message, error) {
	log.Println(msg)
	if Debug {
		chann = debugChannel
	}
	m, err := Session.ChannelMessageSend(chann, msg)
	if err != nil {
		log.Println(err)
	}
	return m, err
}
