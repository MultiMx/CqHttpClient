package state

import "time"

var (
	Start           = time.Now()
	ErrPanicCounter uint
	ErrMsgCounter   uint
	SendMsgCounter  uint
)
