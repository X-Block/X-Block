package events

import (
	"testing"
	"fmt"
)

func TestNewEvent(t *testing.T) {
	event := NewEvent()

	var subscriber1 EventFunc = func(v interface{}){
		fmt.Println("subscriber1 event func.")
	}

	var subscriber2 EventFunc = func(v interface{}){
		fmt.Println("subscriber2 event func.")
	}


}
