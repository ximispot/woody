---
id: requestid
title: RequestID
---

RequestID middleware for [Woody](https://github.com/gowoody/woody) that adds an indentifier to the response.

## Signatures

```go
func New(config ...Config) woody.Handler
```

## Examples

Import the middleware package that is part of the Woody web framework

```go
import (
  "github.com/gowoody/woody/v2"
  "github.com/ximispot/woody/middleware/requestid"
)
```

After you initiate your Woody app, you can use the following possibilities:

```go
// Initialize default config
app.Use(requestid.New())

// Or extend your config for customization
app.Use(requestid.New(requestid.Config{
    Header:    "X-Custom-Header",
    Generator: func() string {
        return "static-id"
    },
}))
```

## Config

```go
// Config defines the config for middleware.
type Config struct {
    // Next defines a function to skip this middleware when returned true.
    //
    // Optional. Default: nil
    Next func(c *woody.Ctx) bool

    // Header is the header key where to get/set the unique request ID
    //
    // Optional. Default: "X-Request-ID"
    Header string

    // Generator defines a function to generate the unique identifier.
    //
    // Optional. Default: utils.UUID
    Generator func() string

    // ContextKey defines the key used when storing the request ID in
    // the locals for a specific request.
    //
    // Optional. Default: requestid
    ContextKey interface{}
}
```

## Default Config
The default config uses a fast UUID generator which will expose the number of
requests made to the server. To conceal this value for better privacy, use the
`utils.UUIDv4` generator.

```go
var ConfigDefault = Config{
    Next:       nil,
    Header:     woody.HeaderXRequestID,
	Generator:  utils.UUID,
	ContextKey: "requestid",
}
```
