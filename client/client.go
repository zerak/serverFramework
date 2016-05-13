package client

import (
	"time"
)

type Client interface {
	String() string
	Exit()

	WLock()
	WUnlock()

	Write(data []byte) (int, error)
	Flush()

	GetID() int64
	GetHBInterval() time.Duration
}
