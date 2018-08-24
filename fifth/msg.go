package fifth

import "log"

var debugChannel = "459341365562572803"     // testing-lab
var spamChannel = "459997163787649025"      // pit of dispair
var importantChannel = "385195528360820739" // general

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