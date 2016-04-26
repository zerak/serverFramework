package protocol

import (
	"bufio"
	"errors"
	"net"
)

const (
	ProtocolHeaderLen = 9 // header1 + cmd4 + length4
)

func init() {
	Register("  v1", &ProtocolV1{})
}

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
	//ctx *context
}

func (p *ProtocolV1) IOLoop(conn net.Conn) error {
	var err error
	//var header byte
	////var cmd uint32
	//var length uint32
	//clientId := atomic.AddInt64(&p.ctx.core.clientIDSequence, 1)
	//client := newClient(clientId, conn, p.ctx)
	//
	//// synchronize the startup of messagePump in order
	//// to guarantee that it gets a chance to initialize
	//// goroutine local state derived from client attributes
	//// and avoid a potential race with IDENTIFY (where a client
	//// could have changed or disabled said attributes)
	//msgPumpStartedChan := make(chan bool)
	//go p.messagePump(client, msgPumpStartedChan)
	//<-msgPumpStartedChan
	//
	//buf := make([]byte, ProtocolHeaderLen)
	//for {
	//	_, err = io.ReadFull(client.Reader, buf)
	//	if err != nil {
	//		ServerLogger.Error("ProtocolV1 read head from client[%v] err->%v buffed->%v", client.Conn.RemoteAddr(), err, client.Reader.Buffered())
	//		break
	//	}
	//
	//	//fmt.Printf("ProtocolV1 recv buf [%v]\n", buf)
	//
	//	// header
	//	header = buf[0]
	//	if header != 0x05 {
	//		err = fmt.Errorf("ProtocolV1 header[%s] err", header)
	//		break
	//	}
	//
	//	// cmd
	//	//cmd = binary.BigEndian.Uint32(buf[1:5])
	//	//fmt.Printf("ProtocolV1 cmd [%v]\n", cmd)
	//
	//	// length
	//	length = binary.BigEndian.Uint32(buf[5:9])
	//	//fmt.Printf("ProtocolV1 length [%v]\n", length)
	//
	//	// data
	//	data := make([]byte, length)
	//	_, err = io.ReadFull(client.Reader, data)
	//	if err != nil {
	//		ServerLogger.Error("ProtocolV1 read data from client[%v] err->%v buffed->%v", client.Conn.RemoteAddr(), err, client.Reader.Buffered())
	//		break
	//	}
	//
	//	//fmt.Printf("ProtocolV1 header[%v] cmd[%v] len[%d] data[%v]\n", header, cmd, length, data)
	//
	//	err = p.Send(client, []byte("string send to client"))
	//	if err != nil {
	//		err = fmt.Errorf("failed to send response ->%s", err)
	//		break
	//	}
	//
	//} // END CLIENT LOOP
	//
	//defer func() {
	//	defer conn.Close()
	//	client.ExitChan <- 1
	//	ServerLogger.Warn("ProtocolV1 client[%v] exit loop err->%v", client.RemoteAddr(), err)
	//}()
	return err
}

func (p *ProtocolV1) Send(w *bufio.Writer, data []byte) (int, error) {
	return 0, errors.New("")
}

//func (p *ProtocolV1) messagePump(client *ClientV1, startedChan chan bool) {
//	var err error
//
//	hbTicker := time.NewTicker(client.HeartbeatInterval)
//	hbChan := hbTicker.C
//
//	//msgTimeOut := client.MsgTimeout
//
//	// signal to the goroutine that started the messagePump
//	// that we've started up
//	close(startedChan)
//
//	for {
//		select {
//		case <-hbChan:
//			err = p.Send(client, []byte("heartBeat"))
//			if err != nil {
//				goto exit
//			}
//		case <-client.ExitChan:
//			goto exit
//		}
//	}
//
//exit:
//	hbTicker.Stop()
//}
//
//func (p *ProtocolV1) Send(client *ClientV1, data []byte) error {
//	client.writeLock.Lock()
//
//	// var zeroTime time.Time
//	// if client.HeartbeatInterval > 0 {
//	// 	client.SetWriteDeadline(time.Now().Add(client.HeartbeatInterval))
//	// } else {
//	// 	client.SetWriteDeadline(zeroTime)
//	// }
//
//	// _, err := SendFramedResponse(client.Writer, frameType, data)
//	_, err := protocol.SendResponse(client.Writer, data)
//	if err != nil {
//		client.writeLock.Unlock()
//		return err
//	}
//
//	// if frameType != frameTypeMessage {
//	// 	err = client.Flush()
//	// }
//
//	client.Writer.Flush()
//
//	client.writeLock.Unlock()
//
//	return err
//}
//
//func parsePack(client *ClientV1) (dd []byte, err error) {
//	return nil, err
//	var header byte
//	var cmd uint32
//	var length uint32
//	buf := make([]byte, ProtocolHeaderLen)
//	_, err = io.ReadFull(client.Reader, buf)
//	if err != nil {
//		if err == io.EOF {
//			ServerLogger.Error("ProtocolV1 read from client %v may be closed", client.Conn.RemoteAddr())
//		}
//		ServerLogger.Error("ProtocolV1 read from client %v err-%v buffed-%v", client.Conn.RemoteAddr(), err, client.Reader.Buffered())
//		return nil, err
//	}
//
//	ServerLogger.Info("ProtocolV1 recv buf [%v]", buf)
//
//	// header
//	header = buf[0]
//	if header != 0x05 {
//		err = fmt.Errorf("ProtocolV1 header[%s] err", header)
//		return nil, err
//	}
//
//	// cmd
//	cmd = binary.BigEndian.Uint32(buf[1:5])
//	ServerLogger.Info("ProtocolV1 cmd [%v]", cmd)
//
//	// length
//	length = binary.BigEndian.Uint32(buf[5:9])
//	ServerLogger.Info("ProtocolV1 length [%v]", length)
//
//	// data
//	data := make([]byte, length)
//	_, err = io.ReadFull(client.Reader, data)
//	if err != nil {
//		if err == io.EOF {
//			ServerLogger.Warn("ProtocolV1 read err-%v client-%v may be closed", err, client.Conn.RemoteAddr())
//		}
//		ServerLogger.Error("ProtocolV1 read err-%v buffed-%v", err, client.Reader.Buffered())
//		return nil, err
//	}
//
//	ServerLogger.Info("ProtocolV1 header[%v] cmd[%v] len[%d] data[%v]", header, cmd, length, data)
//	return nil, err
//}
