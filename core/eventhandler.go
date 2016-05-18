package core

import (
	"strings"

	. "github.com/TaXingTianJi/serverFramework/client"
	. "github.com/TaXingTianJi/serverFramework/protocol"
)

type EventHandler interface {
	ProcessEvent(p *Protocol, client *Client)
}

var eventsHandler = make(map[string]EventHandler)

func ResigerEvent(event string, handle EventHandler) {
	if handle == nil {
		Error("events register handle is nil")
		return
	}

	if _, ok := eventsHandler[strings.ToUpper(event)]; ok {
		Warn("the event", event, "already registed")
		return
	}

	eventsHandler[strings.ToUpper(event)] = handle
}
