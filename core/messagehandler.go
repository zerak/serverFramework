package core

import (
	. "github.com/TaXingTianJi/serverFramework/client"
	. "github.com/TaXingTianJi/serverFramework/protocol"
)

type MsgHandler interface {
	ProcessMsg(p Protocol, client Client, msg *Message)
}

var msgHandle = make(map[string]MsgHandler)

func RegisterMsg(name string, adapter MsgHandler) {
	Info("register msg handler", name)

	if adapter == nil {
		Error("message handler register adapter is nil")
		return
	}

	if _, ok := msgHandle[name]; ok {
		Warn("the handler", name, "already registed")
		return
	}

	msgHandle[name] = adapter
}
