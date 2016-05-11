package core

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"sync/atomic"
	"time"

	. "serverFramework/client"
	"serverFramework/protocol"
	"strconv"
)

const (
	ProtocolHeaderLen = 9 // header1 + cmd4 + length4
)

/*
//the protocol v1 be made of [header cmd len data]
// header: one byte
// cmd: four byte,
// len: four byte,length of the data
// data: the contents

//   	 [x]    [x][x][x][x]  [x][x][x][x]  [x][x][x][x]...
//   |  byte  |  (uint32)  |   (uint32)  	 |   (binary)
//   | 1-byte |   4-byte   |    4-byte   	 |	 N-byte
//   --------------------------------------------...
//    header    cmd	      length         data
//	 [  0   1234  5678  ...]
//   [0x05  cmd   len  data]
*/
type ProtocolV1 struct {
}

func init() {
	protocol.Register("  v1", &ProtocolV1{})
}

func (p *ProtocolV1) IOLoop(conn net.Conn) error {
	var err error
	var header byte
	var cmd uint32
	var length uint32
	clientId := atomic.AddInt64(&ServerApp.clientIDSequence, 1)
	client := NewClient(clientId, conn)

	// synchronize the startup of messagePump in order
	// to guarantee that it gets a chance to initialize
	// goroutine local state derived from client attributes
	// and avoid a potential race with IDENTIFY (where a client
	// could have changed or disabled said attributes)
	msgPumpStartedChan := make(chan bool)
	go p.messagePump(client, msgPumpStartedChan)
	<-msgPumpStartedChan

	buf := make([]byte, ProtocolHeaderLen)
	for {
		_, err = io.ReadFull(client.Reader, buf)
		if err != nil {
			ServerLogger.Error("ProtocolV1 read head from client[%v] err->%v buffed->%v", client.Conn.RemoteAddr(), err, client.Reader.Buffered())
			break
		}

		// header
		header = buf[0]
		if header != 0x05 {
			err = fmt.Errorf("ProtocolV1 header[%s] err", header)
			break
		}

		// cmd
		cmd = binary.BigEndian.Uint32(buf[1:5])

		// length
		length = binary.BigEndian.Uint32(buf[5:9])

		// data
		data := make([]byte, length)
		_, err = io.ReadFull(client.Reader, data)
		if err != nil {
			ServerLogger.Error("ProtocolV1 read data from client[%v] err->%v buffed->%v", client.Conn.RemoteAddr(), err, client.Reader.Buffered())
			break
		}

		ServerLogger.Info("ProtocolV1 header[%v] cmd[%v] len[%d] data[%v]", header, cmd, length, string(data))

		// new msg
		msg := client.Pool.Get().(*Message)
		//msg := NewMessage(int(cmd), data, client)
		msg.ID = int(cmd)
		msg.Body = data
		msg.Conn = client.Conn

		client.MsgChan <- msg
	}

	defer func() {
		client.ExitChan <- 1
		client.Exit()
		ServerLogger.Warn("ProtocolV1 client[%v] exit loop err->%v", client.RemoteAddr(), err)
	}()
	return err
}

func (p *ProtocolV1) Send(c Client, data []byte) (int, error) {
	c.WLock()

	//// todo
	//// set write deadline or not
	//var zeroTime time.Time
	//if c.GetHBInterval() > 0 {
	//	c.SetWriteDeadline(time.Now().Add(c.GetHBInterval()))
	//} else {
	//	c.SetWriteDeadline(zeroTime)
	//}

	// todo
	// check write len(data) size buf
	n, err := c.Write(data)
	if err != nil {
		c.WUnlock()
		return n, err
	}
	c.Flush()
	c.WUnlock()

	return n, nil
}

func (p *ProtocolV1) messagePump(client *ClientV1, startedChan chan bool) {
	var err error

	hbTicker := time.NewTicker(client.HeartbeatInterval)
	hbChan := hbTicker.C
	msgChan := client.MsgChan

	//msgTimeOut := client.MsgTimeout

	// signal to the goroutine that started the messagePump
	// that we've started up
	close(startedChan)

	for {
		select {
		case <-hbChan:
			_, err = p.Send(client, []byte("s2c heartbeat"))
			if err != nil {
				goto exit
			}
		case msg, ok := <-msgChan:
			if ok {
				ServerLogger.Info("cid[%v] recv msg id[%v] b[%v] t[%v]", client.GetID(), msg.ID, string(msg.Body), msg.Timestamp)

				if _, ok := msgHandle[strconv.Itoa(msg.ID)]; ok {
					msgHandle[strconv.Itoa(msg.ID)].ProcessMsg(p, client)
				} else {
					ServerLogger.Warn("cid[%v] unhandle msg[%v]", client.GetID(), msg.ID)
					goto exit
				}
			}
		case <-client.ExitChan:
			goto exit
		}
	}

exit:
	hbTicker.Stop()
	if err != nil {
		ServerLogger.Warn("message pump error[%v]", err)
	}
}

func (p *ProtocolV1) decodePack(client *ClientV1) (dd []byte, err error) {
	return nil, err
	var header byte
	var cmd uint32
	var length uint32
	buf := make([]byte, ProtocolHeaderLen)
	_, err = io.ReadFull(client.Reader, buf)
	if err != nil {
		if err == io.EOF {
			ServerLogger.Error("ProtocolV1 read from client %v may be closed", client.Conn.RemoteAddr())
		}
		ServerLogger.Error("ProtocolV1 read from client %v err-%v buffed-%v", client.Conn.RemoteAddr(), err, client.Reader.Buffered())
		return nil, err
	}

	ServerLogger.Info("ProtocolV1 recv buf [%v]", buf)

	// header
	header = buf[0]
	if header != 0x05 {
		err = fmt.Errorf("ProtocolV1 header[%s] err", header)
		return nil, err
	}

	// cmd
	cmd = binary.BigEndian.Uint32(buf[1:5])
	ServerLogger.Info("ProtocolV1 cmd [%v]", cmd)

	// length
	length = binary.BigEndian.Uint32(buf[5:9])
	ServerLogger.Info("ProtocolV1 length [%v]", length)

	// data
	data := make([]byte, length)
	_, err = io.ReadFull(client.Reader, data)
	if err != nil {
		if err == io.EOF {
			ServerLogger.Warn("ProtocolV1 read err-%v client-%v may be closed", err, client.Conn.RemoteAddr())
		}
		ServerLogger.Error("ProtocolV1 read err-%v buffed-%v", err, client.Reader.Buffered())
		return nil, err
	}

	ServerLogger.Info("ProtocolV1 header[%v] cmd[%v] len[%d] data[%v]", header, cmd, length, data)
	return nil, err
}

func (p *ProtocolV1) encodePack(header byte, cmd, length int, data []byte) {

}
