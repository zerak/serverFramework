package core

import (
	// "bytes"
	// "encoding/binary"
	"fmt"
	"io"
	"net"
	"sync/atomic"
)

/*
//the protocol v1 be made of [header cmd len data]
// header: one byte
// cmd: two byte,
// len: two byte,length of the data
// data: the contents

//   	[x]	   [x][x]      [x][x]  [x][x][x][x]...
//   | byte | (uint32) |  (uint32)  |   (binary)
//   |1-byte|  2-byte  |   2-byte   |	 N-byte
//   --------------------------------------------...
//    header    cmd	      length         data
//
//   [0x05 cmd len data]
*/
type ProtocolV1 struct {
	ctx *context
}

func (p *ProtocolV1) IOLoop(conn net.Conn) error {
	fmt.Printf("[ProtocolV1::Loop] loop ...\n")
	var err error
	// var header byte
	// var cmd uint32
	// var length int32

	clientId := atomic.AddInt64(&p.ctx.core.clientIDSequence, 1)
	client := newClient(clientId, conn, p.ctx)

	buf := make([]byte, 8)
	len := 0

	for {
		for {
			// var n int
			n, err := client.Reader.Read(buf[len:])
			// n, err := client.Read(buf[len:])
			fmt.Printf("len[%d] n[%d] err[%v]\n",len, n, err)
			if err != nil {
				if err != io.EOF {
					// HANDLE CLIENT READ ERR
					fmt.Printf("ProtocolV1 read err-%v\n", err)
				} else {
					err = nil
					len = 0
				}

				fmt.Printf("ProtocolV1 recv buf[%v]\n", buf)
				break
			}

			if n > 0 {
				len += n
			}

			if n == 0 {
				break
			}
		}

		// // header
		// err = binary.Read(conn, binary.BigEndian, &header)
		// if err != nil {
		// 	break
		// }
		// if header != 0x05 {
		// 	err = fmt.Errorf("ProtocolV1 header[%s] err\n", header)
		// 	break
		// }
		// fmt.Printf("ProtocolV1 header[%v]\n", header)

		// // cmd
		// err = binary.Read(client.Reader, binary.BigEndian, &cmd)
		// if err != nil {
		// 	fmt.Printf("cmd\n")
		// 	break
		// }
		// fmt.Printf("ProtocolV1 cmd[%v]\n", cmd)

		// // data
		// err = binary.Read(client.Reader, binary.BigEndian, &length)
		// if err != nil {
		// 	fmt.Printf("len\n")
		// 	break
		// }
		// data := make([]byte, length)
		// _, err = io.ReadFull(client.Reader, data)
		// if err != nil {
		// 	fmt.Printf("data\n")
		// 	break
		// }
		// fmt.Printf("ProtocolV1 len[%d] data[%v]\n", length, data)

		// err = p.Send(client, []byte("string send to client"))
		// if err != nil {
		// 	err = fmt.Errorf("failed to send response - %s", err)
		// 	break
		// }

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
