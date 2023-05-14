---
id: compress
title: Compress
---

Compression middleware for [Woody](https://github.com/ximispot/woody) that will compress the response using `gzip`, `deflate` and `brotli` compression depending on the [Accept-Encoding](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Accept-Encoding) header.

## Signatures

```go
func New(config ...Config) woody.Handler
```

## Examples

Import the middleware package that is part of the Woody web framework

```go
import (
  "github.com/ximispot/woody"
  "github.com/ximispot/woody/middleware/compress"
)
```

After you initiate your Woody app, you can use the following possibilities:

```go
// Initialize default config
app.Use(compress.New())

// Or extend your config for customization
app.Use(compress.New(compress.Config{
    Level: compress.LevelBestSpeed, // 1
}))

// Skip middleware for specific routes
app.Use(compress.New(compress.Config{
  Next:  func(c *woody.Ctx) bool {
    return c.Path() == "/dont_compress"
  },
  Level: compress.LevelBestSpeed, // 1
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

    // Level determines the compression algoritm
    //
    // Optional. Default: LevelDefault
    // LevelDisabled:         -1
    // LevelDefault:          0
    // LevelBestSpeed:        1
    // LevelBestCompression:  2
    Level int
}
```

## Default Config

```go
var ConfigDefault = Config{
    Next:  nil,
    Level: LevelDefault,
}
```

## Constants

```go
// Compression levels
const (
    LevelDisabled        = -1
    LevelDefault         = 0
    LevelBestSpeed       = 1
    LevelBestCompression = 2
)
```
