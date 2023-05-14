---
id: redirect
title: Redirect
---

Redirection middleware for Woody.

## Signatures

```go
func New(config ...Config) woody.Handler
```

## Examples

```go
package main

import (
  "github.com/gowoody/woody/v2"
  "github.com/ximispot/woody/middleware/redirect"
)

func main() {
  app := woody.New()
  
  app.Use(redirect.New(redirect.Config{
    Rules: map[string]string{
      "/old":   "/new",
      "/old/*": "/new/$1",
    },
    StatusCode: 301,
  }))
  
  app.Get("/new", func(c *woody.Ctx) error {
    return c.SendString("Hello, World!")
  })
  app.Get("/new/*", func(c *woody.Ctx) error {
    return c.SendString("Wildcard: " + c.Params("*"))
  })
  
  app.Listen(":3000")
}
```

**Test:**

```curl
curl http://localhost:3000/old
curl http://localhost:3000/old/hello
```

## Config

```go
// Config defines the config for middleware.
type Config struct {
	// Filter defines a function to skip middleware.
	// Optional. Default: nil
	Next func(*woody.Ctx) bool

	// Rules defines the URL path rewrite rules. The values captured in asterisk can be
	// retrieved by index e.g. $1, $2 and so on.
	// Required. Example:
	// "/old":              "/new",
	// "/api/*":            "/$1",
	// "/js/*":             "/public/javascripts/$1",
	// "/users/*/orders/*": "/user/$1/order/$2",
	Rules map[string]string

	// The status code when redirecting
	// This is ignored if Redirect is disabled
	// Optional. Default: 302 (woody.StatusFound)
	StatusCode int

	rulesRegex map[*regexp.Regexp]string
}
```

## Default Config

```go
var ConfigDefault = Config{
	StatusCode: woody.StatusFound,
}
```
