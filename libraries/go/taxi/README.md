# Taxi 🚕

![](https://media.giphy.com/media/SLr8qaoRH6Hmw/source.gif)

Taxi is a lightweight, remote procedure call (RPC) framework for services written in Go. Requests are sent as JSON over HTTP, making it easy to interact with Taxi-based services. Taxi does not impose any particular URL scheme or request structure, and responses are fully-customisable, giving you full control over the design of your API.

## Installation

```go
go get github.com/jakewright/taxi
```

## Usage

Taxi is made up of two components: a client for dispatching RPCs to a remote service, and a router for handling incoming RPCs.

### Client

#### Creating a new client

```go
// Create a new client
client := taxi.NewClient()

// Taxi uses an http.Client{} to send requests. 
// You can provide your own if you wish using
// NewClientUsing(doer Doer).

type Doer interface {
	Do(r *http.Request) (*http.Response, error)
}

custom := &http.Client{Timeout: 30 * time.Second}
client := taxi.NewClientUsing(custom)
```

#### Mock client

It's useful in tests to use a mock client. The mock client has an http.Handler that it uses to serve requests. Set this to an instance of `TestFixture` to create a handler that expects and responds to particular requests.


### Router

Use the router to build a server that responds to RPCs.

```go
// Create a new router (optionally set a logger)
router := taxi.NewRouter().WithLogger(log.Printf)

// Set global middleware
router.UseMiddleware(...)

// Register a handler
router.RegisterHandlerFund("GET", "/foo", fooHandler)

func fooHandler(ctx context.Context, decode taxi.Decoder) (interface{}, error) {
	body := &sruct{
    	Bar string `json:"bar"`
    }{}
    if err := decode(bar); err != nil {
    	return err
    }
    
	...
}
```
