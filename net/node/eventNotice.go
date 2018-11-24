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

