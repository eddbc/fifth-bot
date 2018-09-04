package fifth

import "log"

var debugChannel = "459341365562572803"     // testing-lab
var importantChannel = "485560181732802578" // general
var spamChannel = "485560181732802578"      // also general

func SendMsg(msg string) {
	SendMsgToChan(spamChannel, msg)
}

func SendImportantMsg(msg string) {
	SendMsgToChan(importantChannel, msg)
}

func SendDebugMsg(msg string) {
	Session.ChannelMessageSend(debugChannel, msg)
}

func SendMsgToChan(chann string, msg string) {
	log.Println(msg)
	if !Debug {
		Session.ChannelMessageSend(chann, msg)
	} else {
		Session.ChannelMessageSend(debugChannel, msg)
	}
}
