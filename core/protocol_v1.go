package core

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"sync/atomic"
)

const (
	ProtocolHeaderLen = 9 // header1 + cmd4 + length4
)

/*
//the protocol v1 be made of [header cmd len data]
// header: one byte
// cmd: two byte,
// len: two byte,length of the data
// data: the contents

//   	[x]	   [x][x]      [x][x]  [x][x][x][x]...
//   | byte | (uint32) |  (uint32)  |   (binary)
//   |1-byte|  4-byte  |   4-byte   |	 N-byte
//   --------------------------------------------...
//    header    cmd	      length         data
//	 [  0   1234  5678  ...]
//   [0x05  cmd   len  data]
*/
type ProtocolV1 struct {
	ctx *context
}

func (p *ProtocolV1) IOLoop(conn net.Conn) error {
	fmt.Printf("[ProtocolV1::Loop] loop ...\n")
	var err error
	var header byte
	var cmd uint32
	var length int32

	clientId := atomic.AddInt64(&p.ctx.core.clientIDSequence, 1)
	client := newClient(clientId, conn, p.ctx)

	buf := make([]byte, ProtocolHeaderLen)

	for {
		buf[0] = 0
		buf, err = client.Reader.Peek(ProtocolHeaderLen)
		if err != nil {
			if err == bufio.ErrBufferFull {
				// no data in the buff continue wait
				fmt.Printf("ProtocolV1 read err-%v\n", err)
				continue
			}
			fmt.Printf("ProtocolV1 recv buf[%v]\n", buf)
			break
		}


		// header
		header = buf[0]
		if header != 0x05 {
			err = fmt.Errorf("ProtocolV1 header[%s] err\n", header)
			continue
		}

		// cmd
		cmd_buf := bytes.NewBuffer(buf[1:5])
		binary.Read(cmd_buf, binary.BigEndian, &cmd)

		// data
		len_buf := bytes.NewBuffer(buf[5:9])
		binary.Read(len_buf, binary.BigEndian, &length)
		data := make([]byte, ProtocolHeaderLen+length)
		n, err := client.Reader.Read(data)
		if err == nil {
			if (int32)(n) < (length + ProtocolHeaderLen) {
				fmt.Printf("ProtocolV1 recv[%d] err\n", n, length, data)
				panic(errors.New("recv bytes not enough"))
			}
		}

		data = data[ProtocolHeaderLen:]

		fmt.Printf("ProtocolV1 header[%v] cmd[%v] len[%d] data[%v]\n", header, cmd, length, data)

		err = p.Send(client, []byte("string send to client"))
		if err != nil {
			err = fmt.Errorf("failed to send response - %s", err)
			break
		}

	} // END CLIENT LOOP

	fmt.Printf("[ProtocolV1::Loop] loop exit err - %v\n", err)

	defer conn.Close()
	defer func() {
		fmt.Printf("[ProtocolV1::Loop] defer func ...\n")
	}()
	return err
}

func (p *ProtocolV1) Send(client *ClientV1, data []byte) error {
	client.writeLock.Lock()

	// var zeroTime time.Time
	// if client.HeartbeatInterval > 0 {
	// 	client.SetWriteDeadline(time.Now().Add(client.HeartbeatInterval))
	// } else {
	// 	client.SetWriteDeadline(zeroTime)
	// }

	// _, err := SendFramedResponse(client.Writer, frameType, data)
	_, err := SendResponse(client.Writer, data)
	if err != nil {
		client.writeLock.Unlock()
		return err
	}

	// if frameType != frameTypeMessage {
	// 	err = client.Flush()
	// }

	client.Writer.Flush()

	client.writeLock.Unlock()

	return err
}

func (p *ProtocolV1) readInt(r *bufio.Reader) (int, error) {
	//b, err := r.ReadBytes(delimEnd)
	//if err != nil {
	//	return Resp{}, err
	//}
	//i, err := strconv.ParseInt(string(b[1:len(b)-2]), 10, 64)
	//if err != nil {
	//	return Resp{}, errParse
	//}
	return 1, nil
}
