package events

import (
	"sync"
	"errors"
)



type Event struct {
	m           sync.RWMutex
	subscribers map[EventType]map[Subscriber]EventFunc
}

func NewEvent() *Event {
	return &Event{
		subscribers: make(map[EventType]map[Subscriber]EventFunc),
	}
}

