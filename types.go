package pub

import "time"

type Header struct {
	name string
	eventTime time.Time
}

func(h *Header) EventTime() time.Time {
	return h.eventTime
}

func(h *Header) EventName() string {
	return h.name
}

type EventPayload struct {
	Data interface{}
	Header
}

type subFunc func(data EventPayload)
