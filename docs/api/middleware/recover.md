---
id: recover
title: Recover
---

Recover middleware for [Woody](https://github.com/gowoody/woody) that recovers from panics anywhere in the stack chain and handles the control to the centralized [ErrorHandler](https://docs.gowoody.io/error-handling).

## Signatures

```go
func New(config ...Config) woody.Handler
```

## Examples

Import the middleware package that is part of the Woody web framework

```go
import (
  "github.com/gowoody/woody/v2"
  "github.com/ximispot/woody/middleware/recover"
)
```

After you initiate your Woody app, you can use the following possibilities:

```go
// Initialize default config
app.Use(recover.New())

// This panic will be caught by the middleware
app.Get("/", func(c *woody.Ctx) error {
    panic("I'm an error")
})
```

## Config

```go
// Config defines the config for middleware.
type Config struct {
    // Next defines a function to skip this middleware when returned true.
    //
    // Optional. Default: nil
    Next func(c *woody.Ctx) bool

    // EnableStackTrace enables handling stack trace
    //
    // Optional. Default: false
    EnableStackTrace bool

    // StackTraceHandler defines a function to handle stack trace
    //
    // Optional. Default: defaultStackTraceHandler
    StackTraceHandler func(c *woody.Ctx, e interface{})
}
```

## Default Config

```go
var ConfigDefault = Config{
    Next:              nil,
    EnableStackTrace:  false,
    StackTraceHandler: defaultStackTraceHandler,
}
```
