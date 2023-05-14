---
id: rewrite
title: Rewrite
---

Rewrite middleware rewrites the URL path based on provided rules. It can be helpful for backward compatibility or just creating cleaner and more descriptive links.


## Signatures

```go
func New(config ...Config) woody.Handler
```

### Examples
```go
package main

import (
  "github.com/ximispot/woody"
  "github.com/ximispot/woody/middleware/rewrite"
)

func main() {
  app := woody.New()
  
  app.Use(rewrite.New(rewrite.Config{
    Rules: map[string]string{
      "/old":   "/new",
      "/old/*": "/new/$1",
    },
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
