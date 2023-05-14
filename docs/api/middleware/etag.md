---
id: etag
title: ETag
---

ETag middleware for [Woody](https://github.com/gowoody/woody) that lets caches be more efficient and save bandwidth, as a web server does not need to resend a full response if the content has not changed.

## Signatures

```go
func New(config ...Config) woody.Handler
```

## Examples

Import the middleware package that is part of the Woody web framework

```go
import (
  "github.com/gowoody/woody/v2"
  "github.com/ximispot/woody/middleware/etag"
)
```

After you initiate your Woody app, you can use the following possibilities:

```go
// Initialize default config
app.Use(etag.New())

// Get / receives Etag: "13-1831710635" in response header
app.Get("/", func(c *woody.Ctx) error {
    return c.SendString("Hello, World!")
})

// Or extend your config for customization
app.Use(etag.New(etag.Config{
    Weak: true,
}))

// Get / receives Etag: "W/"13-1831710635" in response header
app.Get("/", func(c *woody.Ctx) error {
    return c.SendString("Hello, World!")
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

    // Weak indicates that a weak validator is used. Weak etags are easy
    // to generate, but are far less useful for comparisons. Strong
    // validators are ideal for comparisons but can be very difficult
    // to generate efficiently. Weak ETag values of two representations
    // of the same resources might be semantically equivalent, but not
    // byte-for-byte identical. This means weak etags prevent caching
    // when byte range requests are used, but strong etags mean range
    // requests can still be cached.
    Weak bool
}
```

## Default Config

```go
var ConfigDefault = Config{
    Next: nil,
    Weak: false,
}
```
