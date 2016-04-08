package core

type Message struct {
	ID int
	Body []byte
	Timestamp int64
	Attempts uint16

}
