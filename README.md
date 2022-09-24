# vandathron/pub

Package `vandathron/pub` is designed based on the publish/subscribe model to 
execute independently executing functions,tasks or subscribers by sending events to its subscribers. 
Main Features:

* Event creation. Events can be created by calling the CreateEvent method
* Registering/Unregistering subscribers to events
* Configurable logging. Interface `Logger` can be implemented to configure logging
* Simple Concurrent execution of subscribers

---

* [Install](#install)
* [Examples](#examples)
* [Full Example](#full-example)

## Install

With [go](https://go.dev/doc/install) installed, you can install this package with this command

```sh
go get -u github.com/vandathron/pub
```

---
## Examples

Let's start by creating a publisher with a single event

```go
func main() {
	p := pub.NewPublisher()
	p.CreateEvent("USER.CREATED")
}
```

We can then declare subscribers, listeners or functions. Parameter must conform with the example below
```go
func SendWelcomeEmail(data pub.EventPayload){
    u := data.Data.(User)
    fmt.Println(u.name)
    fmt.Println("Sending welcome mail..")
    ...
}

func SendSMS(data pub.EventPayload){
    u := data.Data.(User)
    fmt.Println(u.name)
    fmt.Println("Sending phone message..")
    ...
}
```
Subscribe listeners to event:
```go
p.Subscribe("USER.CREATED", SendSMS, SendWelcomeEmail)
 ```
6
Event can be emitted by calling `Publish(event string)` method. This publishes event and execute subscribers concurrently.

```go
p.Publish("USER.CREATED", pub.EventPayload{
    Data:   User{name: "van"},
    Header: pub.Header{},
})
```

Listeners can be unsubscribed from an event:
```go
p.Unsubscribe(SendSMS, "USER.CREATED")
```
Logging is disabled by default. We can create a logger that implements the `Logger` interface.

```go
type EvtLogger struct {

}

func (e *EvtLogger) LogInfo(msg string) {
    ...
}

func (e *EvtLogger) LogErr(msg string) {
    ...
}
```

Logging can now be enabled 
```go
p.DisableLogs = false
p.SetLogger(&EvtLogger)
```
---
## Full Example

```go
type EvtLogger struct {

}

var _ pub.Logger = new(EvtLogger)

func (e EvtLogger) LogInfo(msg string) {
   
}

func (e EvtLogger) LogErr(msg string) {
    
}

type User struct {
    name string
}

func main(){
	p := pub.NewPublisher()
	p.DisableLogs = false;
	p.SetLogger(&EvtLogger{})
	
	p.CreateEvent("USER.CREATED")
	p.CreateEvent("USER.DELETED")

	p.Subscribe(SendSMS, "USER.CREATED")
	p.Subscribe(SendWelcomeEmail, "USER.CREATED")

	p.Publish("USER.CREATED", pub.EventPayload{
		Data:   User{name: "van"},
		Header: pub.Header{},
	})

}
```