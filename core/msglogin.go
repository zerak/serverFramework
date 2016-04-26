package core

import (
	"fmt"
)

type MsgLogin struct {
}

func (m *MsgLogin) ProcessMsg(p *ProtocolV1, client *ClientV1) {
	Info("msg login")

	err := p.Send(client, []byte("string send to client"))
	if err != nil {
		err = fmt.Errorf("failed to send response ->%s", err)
		client.Close()
	}
}
