---
id: adaptor
title: Adaptor
---

Converter for net/http handlers to/from Woody request handlers, special thanks to [@arsmn](https://github.com/arsmn)!

## Signatures
| Name | Signature | Description
| :--- | :--- | :---
| HTTPHandler | `HTTPHandler(h http.Handler) woody.Handler` | http.Handler -> woody.Handler
| HTTPHandlerFunc | `HTTPHandlerFunc(h http.HandlerFunc) woody.Handler` | http.HandlerFunc -> woody.Handler
| HTTPMiddleware | `HTTPHandlerFunc(mw func(http.Handler) http.Handler) woody.Handler` | func(http.Handler) http.Handler -> woody.Handler
| WoodyHandler | `WoodyHandler(h woody.Handler) http.Handler` | woody.Handler -> http.Handler
| WoodyHandlerFunc | `WoodyHandlerFunc(h woody.Handler) http.HandlerFunc` | woody.Handler -> http.HandlerFunc
| WoodyApp | `WoodyApp(app *woody.App) http.HandlerFunc` | Woody app -> http.HandlerFunc
| CopyContextToWoodyContex | `CopyContextToWoodyContext(context interface{}, requestContext *fasthttp.RequestCtx)` | context.Context -> fasthttp.RequestCtx

## Examples

### net/http to Woody
```go
package main

import (
	"fmt"
	"net/http"

	"github.com/ximispot/v2/middleware/adaptor"
	"github.com/ximispot/woody"
)

func main() {
	// New woody app
	app := woody.New()

	// http.Handler -> woody.Handler
	app.Get("/", adaptor.HTTPHandler(handler(greet)))

	// http.HandlerFunc -> woody.Handler
	app.Get("/func", adaptor.HTTPHandlerFunc(greet))

	// Listen on port 3000
	app.Listen(":3000")
}

func handler(f http.HandlerFunc) http.Handler {
	return http.HandlerFunc(f)
}

func greet(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello World!")
}
```

### net/http middleware to Woody
```go
package main

import (
	"log"
	"net/http"

	"github.com/ximispot/v2/middleware/adaptor"
	"github.com/ximispot/woody"
)

func main() {
	// New woody app
	app := woody.New()

	// http middleware -> woody.Handler
	app.Use(adaptor.HTTPMiddleware(logMiddleware))

	// Listen on port 3000
	app.Listen(":3000")
}

func logMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("log middleware")
		next.ServeHTTP(w, r)
	})
}
```

### Woody Handler to net/http
```go
package main

import (
	"net/http"

	"github.com/ximispot/v2/middleware/adaptor"
	"github.com/ximispot/woody"
)

func main() {
	// woody.Handler -> http.Handler
	http.Handle("/", adaptor.WoodyHandler(greet))

  	// woody.Handler -> http.HandlerFunc
	http.HandleFunc("/func", adaptor.WoodyHandlerFunc(greet))

	// Listen on port 3000
	http.ListenAndServe(":3000", nil)
}

func greet(c *woody.Ctx) error {
	return c.SendString("Hello World!")
}
```

### Woody App to net/http
```go
package main

import (
	"github.com/ximispot/v2/middleware/adaptor"
	"github.com/ximispot/woody"
	"net/http"
)
func main() {
	app := woody.New()

	app.Get("/greet", greet)

	// Listen on port 3000
	http.ListenAndServe(":3000", adaptor.WoodyApp(app))
}

func greet(c *woody.Ctx) error {
	return c.SendString("Hello World!")
}
```
