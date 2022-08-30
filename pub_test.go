package pub

import (
	"fmt"
	"testing"
	"time"
)

type publisherTest struct {
	title string
	event string
	subscribers map[string] *subQueue
	disableLogs bool
	pub *Publisher


}

func TestEvents(t *testing.T) {
	tests := []publisherTest {
		{
			title: "Event should exist after created",
			event: "event-1",
			pub: NewPublisher(),
		},
	}

	for _, test := range tests {
		t.Run(test.title, func(t *testing.T) {
			testCreateEvent(t,test)
		})
	}
}


func TestSubscription(t *testing.T){
	tests := []publisherTest {
		{
			title: "Subscribers should be called",
			event: "event-1",
			pub: NewPublisher(),
		},
	}

	for _, test := range tests {
		t.Run(test.title, func(t *testing.T) {
			testSubscribeToEvent(t,test)
			testUnsubscribeFromEvent(t,test)
		})
	}
}

func testCreateEvent(t *testing.T, publisher publisherTest){
	publisher.pub.CreateEvent(publisher.event)

	if !publisher.pub.EventExist(publisher.event){
		t.Errorf("%s does not exist. Expected: %v", publisher.event, true)
	}
}

func testSubscribeToEvent(t *testing.T, publisher publisherTest) {
	publisher.pub.CreateEvent(publisher.event)
	called := 0
	publisher.pub.Subscribe(func(data EventPayload) {
		fmt.Println("Subscriber!!!")
		called++
	}, publisher.event)

	publisher.pub.Publish(publisher.event, EventPayload{})
	time.Sleep(2 * time.Second)
	if called != 1 {
		t.Errorf("TestSubscribeToEvent failed to call subscribers. Expected: %v. Actual: %v", 1, called)
	}
}

func testUnsubscribeFromEvent(t *testing.T, publisher publisherTest){
	publisher.pub.CreateEvent(publisher.event)
	sub := func(data EventPayload) {
		fmt.Println("Unsub")
	}
	publisher.pub.Subscribe(sub, publisher.event)

	publisher.pub.Unsubscribe(sub, publisher.event)
	got := len(publisher.pub.subscribers)
	if got != 0 {
		t.Errorf("testUnsubscribeFromEvent failed to remove subscriber. Expected: %v. Actual: %v",0, got )
	}
}