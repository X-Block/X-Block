package events

Type  EventType int16

const (

	EventSaveBlock             EventType = 0
	EventReplyTx               EventType = 1
	EventBlockPersistCompleted EventType = 2
	EventNewInventory          EventType = 3
	EventNodeDisconnect        EventType = 4

)
