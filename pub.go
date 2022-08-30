package pub

import (
	"fmt"
	"reflect"
	"time"
)

// Publisher manages events, subscribers and emit events.
type Publisher struct {
	subscribers    map[string] *subQueue
	l  Logger
	// DisableLogs is false by default. Set DisableLogs to true to turn off logging.
	DisableLogs    bool
}

// NewPublisher creates a publisher for creating events, subscribing functions and emitting events.
func NewPublisher() *Publisher {
	return &Publisher{
		subscribers: make(map[string]*subQueue),
		l:      nil,
	}
}

// SetLogger configures logger. Implement Logger interface.
func(d *Publisher) SetLogger(logger Logger){
	d.l = logger
}

// CreateEvent creates a new event.
// Calling this method will remove all existing subscribers for this event if it already exist.
func(d *Publisher) CreateEvent(event string){
	if d.EventExist(event){
		d.subscribers[event].popAll()
	}
	d.subscribers[event] = newQueue(nil)
}

// Subscribe subscribes subscriber to event. Creates the event if event does not exist.
//
// If subscriber already exists, error ErrDuplicateSubscriber is returned.
func(d *Publisher) Subscribe(subscriber subFunc, event string) (bool, error){
	if _, ok := d.subscribers[event]; !ok {
		d.subscribers[event] = newQueue(subscriber)
		return true, nil
	}

	if d.subscriberToEventAlreadyExists(subscriber, event) {
			return false, ErrDuplicateSubscriber
		}

	d.subscribers[event].pushFunc(subscriber)

	subscriberName := getFunctionName(subscriber)

	d.logInf(fmt.Sprintf("Subscriber: %s subscribed to Event: %s", subscriberName, event))

	return true, nil
}

// Unsubscribe unsubscribes subscriber from the event.
// Delete event if no subscriber is registered to the event.
func (d *Publisher) Unsubscribe(subscriber subFunc, event string) (bool, error) {
	if _, ok := d.subscribers[event]; !ok {
		return false, ErrEventDoesNotExist
	}

	eventSubs := d.subscribers[event]

	if !d.subscriberToEventAlreadyExists(subscriber, event) {
		return false, ErrSubscriberDoesNotExist
	}

	subscriberName := getFunctionName(subscriber)

	for idx, f := range eventSubs.subscribers {
		name := getFunctionName(f)

		if subscriberName == name {
			eventSubs.popFuncAt(idx)
		}
	}

	if d.SubscribersCount(event) == 0 {
		delete(d.subscribers, event)
	}

	d.logInf(fmt.Sprintf("Subscriber: %s unsubscribed from Event: %s", subscriberName, event))
	return true, nil
}

// Publish execute each subscriber registered to this event on separate goroutines.
// Returns ErrEventDoesNotExist if event does not exist.
func(d *Publisher) Publish(event string, data EventPayload) (bool, error){
	if _, ok := d.subscribers[event]; !ok {
		return false, ErrEventDoesNotExist
	}

	subscribers := d.subscribers[event].getAllSubs()

	if len(subscribers) == 0 {
		return false, ErrNoSubscribers
	}

	data.Header = Header{
		name:      event,
		eventTime: time.Now(),
	}

	for _, subscriber := range subscribers {
		 execSubscriber := d.subWrapper(subscriber, event).(subFunc)
		 go execSubscriber(data)
	}
	return true, nil
}

func(d *Publisher) subWrapper(f interface{}, event string) interface{} {
	fv := reflect.ValueOf(f)
	subscriberName := getFunctionName(f)
	wrapperFunc := reflect.MakeFunc(reflect.TypeOf(f), func(args []reflect.Value) (results []reflect.Value) {
		defer func(){
			if x := recover(); x != nil {
				d.logErr(fmt.Sprintf("Subscriber: %s processing failed for Event: %s", subscriberName, event))
			}
		}()
		d.logInf(fmt.Sprintf("Executing subscriber: %s. Event: %s", subscriberName, event))
		out := fv.Call(args)
		d.logInf(fmt.Sprintf("Done executing subscriber %s. Event: %s", subscriberName, event))
		return out
	})

	return wrapperFunc.Interface()
}

func(d *Publisher) logErr(msg string){
	if !d.DisableLogs && d.l != nil{
		d.l.LogErr(msg)
	}
}

func(d *Publisher) logInf(msg string){
	if !d.DisableLogs && d.l != nil{
		d.l.LogInfo(msg)
	}
}