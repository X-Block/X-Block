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


func (e *Event) Subscribe(eventtype EventType,eventfunc EventFunc) Subscriber {
	e.m.Lock()
	defer e.m.Unlock()

	sub := make(chan interface{})
	_,ok := e.subscribers[eventtype]
	if !ok {
		e.subscribers[eventtype] =  make(map[Subscriber]EventFunc)
	}
	e.subscribers[eventtype][sub] = eventfunc

	return sub
}


