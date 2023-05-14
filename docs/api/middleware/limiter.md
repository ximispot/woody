---
id: limiter
title: Limiter
---

Limiter middleware for [Woody](https://github.com/ximispot/woody) that is used to limit repeat requests to public APIs and/or endpoints such as password reset. It is also useful for API clients, web crawling, or other tasks that need to be throttled.

:::note
This middleware uses our [Storage](https://github.com/ximispot/storage) package to support various databases through a single interface. The default configuration for this middleware saves data to memory, see the examples below for other databases.
:::

:::note
This module does not share state with other processes/servers by default.
:::

## Signatures

```go
func New(config ...Config) woody.Handler
```

## Examples

Import the middleware package that is part of the Woody web framework

```go
import (
  "github.com/ximispot/woody"
  "github.com/ximispot/woody/middleware/limiter"
)
```

After you initiate your Woody app, you can use the following possibilities:

```go
// Initialize default config
app.Use(limiter.New())

// Or extend your config for customization
app.Use(limiter.New(limiter.Config{
    Next: func(c *woody.Ctx) bool {
        return c.IP() == "127.0.0.1"
    },
    Max:          20,
    Expiration:     30 * time.Second,
    KeyGenerator:          func(c *woody.Ctx) string {
        return c.Get("x-forwarded-for")
    },
    LimitReached: func(c *woody.Ctx) error {
        return c.SendFile("./toofast.html")
    },
    Storage: myCustomStorage{},
}))
```

## Sliding window

Instead of using the standard fixed window algorithm, you can enable the [sliding window](https://en.wikipedia.org/wiki/Sliding_window_protocol) algorithm.

A example of such configuration is:

```go
app.Use(limiter.New(limiter.Config{
    Max:            20,
    Expiration:     30 * time.Second,
    LimiterMiddleware: limiter.SlidingWindow{},
}))
```

This means that every window will take into account the previous window(if there was any). The given formula for the rate is:
```
weightOfPreviousWindpw = previous window's amount request * (whenNewWindow / Expiration)
rate = weightOfPreviousWindpw + current window's amount request.
```

## Config

```go
// Config defines the config for middleware.
type Config struct {
    // Next defines a function to skip this middleware when returned true.
    //
    // Optional. Default: nil
    Next func(c *woody.Ctx) bool

	// Max number of recent connections during `Duration` seconds before sending a 429 response
    //
    // Default: 5
    Max int

    // KeyGenerator allows you to generate custom keys, by default c.IP() is used
    //
    // Default: func(c *woody.Ctx) string {
    //   return c.IP()
    // }
    KeyGenerator func(*woody.Ctx) string

    // Expiration is the time on how long to keep records of requests in memory
    //
    // Default: 1 * time.Minute
    Expiration time.Duration

    // LimitReached is called when a request hits the limit
    //
    // Default: func(c *woody.Ctx) error {
    //   return c.SendStatus(woody.StatusTooManyRequests)
    // }
    LimitReached woody.Handler

    // When set to true, requests with StatusCode >= 400 won't be counted.
    //
    // Default: false
    SkipFailedRequests bool

    // When set to true, requests with StatusCode < 400 won't be counted.
    //
    // Default: false
    SkipSuccessfulRequests bool

    // Store is used to store the state of the middleware
    //
    // Default: an in memory store for this process only
    Storage woody.Storage

    // LimiterMiddleware is the struct that implements limiter middleware.
    //
    // Default: a new Fixed Window Rate Limiter
    LimiterMiddleware LimiterHandler
}
```

:::note
A custom store can be used if it implements the `Storage` interface - more details and an example can be found in `store.go`.
:::

## Default Config

```go
var ConfigDefault = Config{
    Max:        5,
    Expiration: 1 * time.Minute,
    KeyGenerator: func(c *woody.Ctx) string {
        return c.IP()
    },
    LimitReached: func(c *woody.Ctx) error {
        return c.SendStatus(woody.StatusTooManyRequests)
    },
    SkipFailedRequests: false,
    SkipSuccessfulRequests: false,
    LimiterMiddleware: FixedWindow{},
}
```

### Custom Storage/Database

You can use any storage from our [storage](https://github.com/ximispot/storage/) package.

```go
storage := sqlite3.New() // From github.com/ximispot/storage/sqlite3
app.Use(limiter.New(limiter.Config{
	Storage: storage,
}))
```