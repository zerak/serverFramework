package core

type EventHandler interface {
	ProcessEvent(p *ProtocolV1, client *ClientV1)
}

var eventsHandler = make(map[string]EventHandler)

func ResigerEvent(event string, handle EventHandler) {
	if handle == nil {
		Error("events register handle is nil")
		return
	}

	if _, ok := eventsHandler[event]; ok {
		Warn("the event", event, "already registed")
		return
	}

	eventsHandler[event] = handle
}
