package pub

import (
	"reflect"
	"runtime"
)

func getFunctionName(f interface{}) string {
	name := runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
	return name
}

// SubscribersCount returns the count of subscribers to event.
func(d *Publisher) SubscribersCount(event string) int {
	if !d.EventExist(event) {
		return 0
	}
	return len(d.subscribers[event].subscribers)
}

func(d *Publisher) subscriberToEventAlreadyExists(subscriber subFunc, event string) bool{
	if !d.EventExist(event) || d.SubscribersCount(event) == 0 {
		return false
	}
	subscriberName := getFunctionName(subscriber)
	eventBox := d.subscribers[event]
	for _, f := range eventBox.subscribers {
		name := getFunctionName(f)
		if subscriberName == name {
			return true
		}
	}
	return false
}

//EventExist Checks if an event with eventName has been created.
func(d *Publisher) EventExist(event string) bool {
	_, ok := d.subscribers[event]
	return ok
}