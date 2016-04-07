package core

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
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
				// no enough data in the buff continue wait
				fmt.Printf("ProtocolV1 read err-%v\n", err)
				continue
			}
			fmt.Printf("ProtocolV1 recv head err [%v]\n", buf)
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

		// length
		len_buf := bytes.NewBuffer(buf[5:9])
		binary.Read(len_buf, binary.BigEndian, &length)

		// 1 check the total buff size in the io reader
		// 2 read buf and move read buf pointer
	checkData:
		_, err = client.Reader.Peek(ProtocolHeaderLen + length)
		if err != nil {
			if err == bufio.ErrBufferFull {
				// data not enough
				fmt.Printf("ProtocolV1 read err-%v\n", err)
				goto checkData
			}
			fmt.Printf("ProtocolV1 read data err - %v\n", err)
			break
		}

		// data
		data := make([]byte, ProtocolHeaderLen+length)
		pos := 0
		for {
			n, err := client.Reader.Read(data[pos:])
			if err != nil {
				if err == io.EOF {
					fmt.Printf("ProtocolV1 read contents ok\n")
					break
				}
			}

			if n > 0 {
				pos += n
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

	defer func() {
		fmt.Printf("[ProtocolV1::Loop] exit client[%v] loop ...\n", client.RemoteAddr())
	}()
	defer conn.Close()
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
