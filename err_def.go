package pub

import "errors"

var (
	// ErrDuplicateSubscriber is returned when trying to subscribe a function that has already been subscribed
	// to a particular event.
	ErrDuplicateSubscriber = errors.New("duplicate subscriber")

	// ErrSubscriberDoesNotExist is returned when a subscriber that does not exist is tried
	// to unsubscribe
	ErrSubscriberDoesNotExist = errors.New("subscriber does not exist")
	//ErrNoSubscribers is returned if an event with no subscribers is emitted.
	// This can be ignored when the Emit() is called.
	ErrNoSubscribers = errors.New("no subscriber for this event")

	// ErrEventDoesNotExist is returned when an event not created is emitted.
	ErrEventDoesNotExist = errors.New("event does not exist")

)