package core

type MsgHandler interface {
	ProcessMsg(p *ProtocolV1, client *ClientV1)
}

var msgHandle = make(map[string]MsgHandler)

func ResigerMsgHandler(name string, adapter MsgHandler) {
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
