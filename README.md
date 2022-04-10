# Go Server

## Description

Just messing around with Go to figure out basics and build pipeline fundamentals...


## Requirements

- go1.18
- docker 
- REST Client (VSCode Extension; for testing)

## Getting Started

1. Run mongo and mongo-express: `docker-compose up`
2. 


## Testing

For basic HTTP testing, can use the REST Client VSCode Extension to simplify API testing :grinning:

## TODO

- [ ] Get example tests working
- [ ] Generate docs
- [ ] Setup build pipeline
- [ ] Setup mongo within app to write data to
- [x] Watch video on context


## Context 

Main thing thats useful is cancellation, can cancel processing at any point and this propogates to child contexts.
Can send values over context (useful if u want to log etc)

context type
context.Background (dont need any cancellation, this is the root context but you dont need to do anything)


```go
import "context"

ctx := context.background
ctx, cancel := context.WithCancel(ctx)

go func() {
	s := bufio.NewScanner(os.Stdin)
	s.Scan()
	cancel()
}

sleepAndTalk(ctx, 5 * time.Second, "hello")


```

When using values in context, the 2 rules:
- Whatever goes into values should be request specific
- Whatever is in context is extra info that is useful but does not impact the business logic
	- A good example of what is useful is the request id (correlation id)