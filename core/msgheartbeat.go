package core

import (
	"fmt"
)

type MsgHeartBeat struct {
}

func (m *MsgHeartBeat) ProcessMsg(p *ProtocolV1, client *ClientV1) {
	Info("msg heart beat")

	err := p.Send(client, []byte("heartbeat pack send to client"))
	if err != nil {
		err = fmt.Errorf("failed to send response ->%s", err)
		client.Close()
	}
}
