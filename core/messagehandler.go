package core

import (
	. "serverFramework/client"
	. "serverFramework/protocol"
)

type MsgHandler interface {
	ProcessMsg(p Protocol, client Client, msg *Message)
}

var msgHandle = make(map[string]MsgHandler)

func ResigerMsg(name string, adapter MsgHandler) {
	Info("msg handler ", name)

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
