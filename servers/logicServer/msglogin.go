package main

import (
	"fmt"
	"serverFramework/core"
)

type MsgLogin struct {
}

func (m *MsgLogin) ProcessMsg(p *core.ProtocolV1, client *core.ClientV1) {
	core.ServerLogger.Info("msg login")

	err := p.Send(client, []byte("string send to client"))
	if err != nil {
		err = fmt.Errorf("failed to send response ->%s", err)
		client.Close()
	}
}
