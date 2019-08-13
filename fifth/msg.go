package fifth

import (
	"github.com/bwmarrin/discordgo"
	"log"
)

var debugChannel = "459341365562572803"     // testing-lab
var importantChannel = "610240336958062593" // general
var spamChannel = "610240336958062593"      // kill-feed
var timerChannel = "610240336958062593"

//SendMsg Send a message to the default channel
func SendMsg(msg string) (*discordgo.Message, error) {
	return SendMsgToChan(spamChannel, msg)
}

//SendImportantMsg Send message to the default priority channel
func SendImportantMsg(msg string) (*discordgo.Message, error) {
	return SendMsgToChan(importantChannel, msg)
}

//SendDebugMsg Send message to the debug channel
func SendDebugMsg(msg string) (*discordgo.Message, error) {
	return SendMsgToChan(debugChannel, msg)
}

//SendMsgToChan Send message to given channel
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
