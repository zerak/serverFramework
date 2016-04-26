package main

import "serverFramework/core"

func init() {
	core.ResigerMsgHandler(string(10010), &MsgHeartBeat{})
	core.ResigerMsgHandler(string(10011), &MsgLogin{})

	//msgHandle["100"] = &MsgLogin{}
	//msgHandle["a"] = &MsgLogin{}
	//
	//for k, v := range msgHandle{
	//	Info("k[", k, "] v[", v, "]")
	//}
	//
	//if _, ok := msgHandle["a"]; ok {
	//	msgHandle["a"].ProcessMsg()
	//	ServerLogger.Warn("handle msg[%v]", "a")
	//} else {
	//	ServerLogger.Warn("unhandle msg[%v]", 10010)
	//}
	//
	//if _, ok := msgHandle[string(10010)]; ok {
	//	msgHandle[string(10010)].ProcessMsg()
	//} else {
	//	ServerLogger.Warn("unhandle msg[%v]", 10010)
	//}
}
