package main

import (
	"fmt"
	"serverFramework/core"
)

type MsgHeartBeat struct {
}

func (m *MsgHeartBeat) ProcessMsg(p *core.ProtocolV1, client *core.ClientV1) {
	core.ServerLogger.Info("msg heart beat")

	err := p.Send(client, []byte("heartbeat pack send to client"))
	if err != nil {
		err = fmt.Errorf("failed to send response ->%s", err)
		client.Close()
	}
}
