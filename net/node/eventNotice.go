package node

import (
	"XBlock/events"
	"fmt"
)

type eventQueue struct {
	Consensus  *events.Event
	Block      *events.Event
	Disconnect *events.Event
}

func (eq *eventQueue) init() {
	eq.Consensus = events.NewEvent()
	eq.Block = events.NewEvent()
	eq.Disconnect = events.NewEvent()
}

