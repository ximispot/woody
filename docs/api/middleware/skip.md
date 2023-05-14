---
id: skip
title: Skip
---

Skip middleware for [Woody](https://github.com/gowoody/woody) that skips a wrapped handler if a predicate is true.

## Signatures
```go
func New(handler woody.Handler, exclude func(c *woody.Ctx) bool) woody.Handler
```

## Examples
Import the middleware package that is part of the Woody web framework
```go
import (
  "github.com/gowoody/woody/v2"
  "github.com/ximispot/woody/middleware/skip"
)
```

After you initiate your Woody app, you can use the following possibilities:

```go
func main() {
	app := woody.New()

	app.Use(skip.New(BasicHandler, func(ctx *woody.Ctx) bool {
		return ctx.Method() == woody.MethodGet
	}))

	app.Get("/", func(ctx *woody.Ctx) error {
		return ctx.SendString("It was a GET request!")
	})

	log.Fatal(app.Listen(":3000"))
}

func BasicHandler(ctx *woody.Ctx) error {
	return ctx.SendString("It was not a GET request!")
}
```

:::tip
app.Use will handle requests from any route, and any method. In the example above, it will only skip if the method is GET.
:::