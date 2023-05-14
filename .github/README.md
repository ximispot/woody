## Quickstart

This framework based on Woody framework. Fork written for experimental use.

[Workbook](https://ximispot.github.io/woody.io/)

```go
package main

import "github.com/ximispot/woody"

func main() {
    app := woody.New()

    app.Get("/", func(c *woody.Ctx) error {
        return c.SendString("Hello, World!")
    })

    app.Listen(":3000")
}
```

## ⚙️ Installation

Current build tested on version `1.20.4`

Install woody with

```bash
go get -u github.com/ximispot/woody
```

#### **Basic Routing**

```go
func main() {
    app := woody.New()

    // GET /api/register
    app.Get("/api/*", func(c *woody.Ctx) error {
        msg := fmt.Sprintf("✋ %s", c.Params("*"))
        return c.SendString(msg) // => ✋ register
    })

    // GET /flights/LAX-SFO
    app.Get("/flights/:from-:to", func(c *woody.Ctx) error {
        msg := fmt.Sprintf("💸 From: %s, To: %s", c.Params("from"), c.Params("to"))
        return c.SendString(msg) // => 💸 From: LAX, To: SFO
    })

    // GET /dictionary.txt
    app.Get("/:file.:ext", func(c *woody.Ctx) error {
        msg := fmt.Sprintf("📃 %s.%s", c.Params("file"), c.Params("ext"))
        return c.SendString(msg) // => 📃 dictionary.txt
    })

    // GET /john/75
    app.Get("/:name/:age/:gender?", func(c *woody.Ctx) error {
        msg := fmt.Sprintf("👴 %s is %s years old", c.Params("name"), c.Params("age"))
        return c.SendString(msg) // => 👴 john is 75 years old
    })

    // GET /john
    app.Get("/:name", func(c *woody.Ctx) error {
        msg := fmt.Sprintf("Hello, %s 👋!", c.Params("name"))
        return c.SendString(msg) // => Hello john 👋!
    })

    log.Fatal(app.Listen(":3000"))
}

```

#### 📖 Route Naming

```go
func main() {
    app := woody.New()

    // GET /api/register
    app.Get("/api/*", func(c *woody.Ctx) error {
        msg := fmt.Sprintf("✋ %s", c.Params("*"))
        return c.SendString(msg) // => ✋ register
    }).Name("api")

    data, _ := json.MarshalIndent(app.GetRoute("api"), "", "  ")
    fmt.Print(string(data))
    // Prints:
    // {
    //    "method": "GET",
    //    "name": "api",
    //    "path": "/api/*",
    //    "params": [
    //      "*1"
    //    ]
    // }


    log.Fatal(app.Listen(":3000"))
}

```

#### 📖 [**Serving Static Files**](https://docs.ximispot.io/api/app#static)

```go
func main() {
    app := woody.New()

    app.Static("/", "./public")
    // => http://localhost:3000/js/script.js
    // => http://localhost:3000/css/style.css

    app.Static("/prefix", "./public")
    // => http://localhost:3000/prefix/js/script.js
    // => http://localhost:3000/prefix/css/style.css

    app.Static("*", "./public/index.html")
    // => http://localhost:3000/any/path/shows/index/html

    log.Fatal(app.Listen(":3000"))
}

```

#### 📖 [**Middleware & Next**](https://docs.ximispot.io/api/ctx#next)

```go
func main() {
    app := woody.New()

    // Match any route
    app.Use(func(c *woody.Ctx) error {
        fmt.Println("🥇 First handler")
        return c.Next()
    })

    // Match all routes starting with /api
    app.Use("/api", func(c *woody.Ctx) error {
        fmt.Println("🥈 Second handler")
        return c.Next()
    })

    // GET /api/list
    app.Get("/api/list", func(c *woody.Ctx) error {
        fmt.Println("🥉 Last handler")
        return c.SendString("Hello, World 👋!")
    })

    log.Fatal(app.Listen(":3000"))
}

```

<details>
  <summary>📚 Show more code examples</summary>

### Views engines

📖 [Config](https://docs.ximispot.io/api/woody#config)
📖 [Engines](https://github.com/ximispot/template)
📖 [Render](https://docs.ximispot.io/api/ctx#render)

Woody defaults to the [html/template](https://pkg.go.dev/html/template/) when no view engine is set.

If you want to execute partials or use a different engine like [amber](https://github.com/eknkc/amber), [handlebars](https://github.com/aymerick/raymond), [mustache](https://github.com/cbroglie/mustache) or [pug](https://github.com/Joker/jade) etc..

Checkout our [Template](https://github.com/ximispot/template) package that support multiple view engines.

```go
package main

import (
    "github.com/ximispot/woody"
    "github.com/ximispot/template/pug"
)

func main() {
    // You can setup Views engine before initiation app:
    app := woody.New(woody.Config{
        Views: pug.New("./views", ".pug"),
    })

    // And now, you can call template `./views/home.pug` like this:
    app.Get("/", func(c *woody.Ctx) error {
        return c.Render("home", woody.Map{
            "title": "Homepage",
            "year":  1999,
        })
    })

    log.Fatal(app.Listen(":3000"))
}
```

### Grouping routes into chains

📖 [Group](https://docs.ximispot.io/api/app#group)

```go
func middleware(c *woody.Ctx) error {
    fmt.Println("Don't mind me!")
    return c.Next()
}

func handler(c *woody.Ctx) error {
    return c.SendString(c.Path())
}

func main() {
    app := woody.New()

    // Root API route
    api := app.Group("/api", middleware) // /api

    // API v1 routes
    v1 := api.Group("/v1", middleware) // /api/v1
    v1.Get("/list", handler)           // /api/v1/list
    v1.Get("/user", handler)           // /api/v1/user

    // API v2 routes
    v2 := api.Group("/v2", middleware) // /api/v2
    v2.Get("/list", handler)           // /api/v2/list
    v2.Get("/user", handler)           // /api/v2/user

    // ...
}

```

### Middleware logger

📖 [Logger](https://docs.ximispot.io/api/middleware/logger)

```go
package main

import (
    "log"

    "github.com/ximispot/woody"
    "github.com/ximispot/woody/middleware/logger"
)

func main() {
    app := woody.New()

    app.Use(logger.New())

    // ...

    log.Fatal(app.Listen(":3000"))
}
```

### Cross-Origin Resource Sharing (CORS)

📖 [CORS](https://docs.ximispot.io/api/middleware/cors)

```go
import (
    "log"

    "github.com/ximispot/woody"
    "github.com/ximispot/woody/middleware/cors"
)

func main() {
    app := woody.New()

    app.Use(cors.New())

    // ...

    log.Fatal(app.Listen(":3000"))
}
```

Check CORS by passing any domain in `Origin` header:

```bash
curl -H "Origin: http://example.com" --verbose http://localhost:3000
```

### Custom 404 response

📖 [HTTP Methods](https://docs.ximispot.io/api/ctx#status)

```go
func main() {
    app := woody.New()

    app.Static("/", "./public")

    app.Get("/demo", func(c *woody.Ctx) error {
        return c.SendString("This is a demo!")
    })

    app.Post("/register", func(c *woody.Ctx) error {
        return c.SendString("Welcome!")
    })

    // Last middleware to match anything
    app.Use(func(c *woody.Ctx) error {
        return c.SendStatus(404)
        // => 404 "Not Found"
    })

    log.Fatal(app.Listen(":3000"))
}
```

### JSON Response

📖 [JSON](https://docs.ximispot.io/api/ctx#json)

```go
type User struct {
    Name string `json:"name"`
    Age  int    `json:"age"`
}

func main() {
    app := woody.New()

    app.Get("/user", func(c *woody.Ctx) error {
        return c.JSON(&User{"John", 20})
        // => {"name":"John", "age":20}
    })

    app.Get("/json", func(c *woody.Ctx) error {
        return c.JSON(woody.Map{
            "success": true,
            "message": "Hi John!",
        })
        // => {"success":true, "message":"Hi John!"}
    })

    log.Fatal(app.Listen(":3000"))
}
```

### WebSocket Upgrade

📖 [Websocket](https://github.com/ximispot/websocket)

```go
import (
    "github.com/ximispot/woody"
    "github.com/ximispot/woody/middleware/websocket"
)

func main() {
  app := woody.New()

  app.Get("/ws", websocket.New(func(c *websocket.Conn) {
    for {
      mt, msg, err := c.ReadMessage()
      if err != nil {
        log.Println("read:", err)
        break
      }
      log.Printf("recv: %s", msg)
      err = c.WriteMessage(mt, msg)
      if err != nil {
        log.Println("write:", err)
        break
      }
    }
  }))

  log.Fatal(app.Listen(":3000"))
  // ws://localhost:3000/ws
}
```

### Server-Sent Events

📖 [More Info](https://developer.mozilla.org/en-US/docs/Web/API/Server-sent_events/Using_server-sent_events)

```go
import (
    "github.com/ximispot/woody"
    "github.com/valyala/fasthttp"
)

func main() {
  app := woody.New()

  app.Get("/sse", func(c *woody.Ctx) error {
    c.Set("Content-Type", "text/event-stream")
    c.Set("Cache-Control", "no-cache")
    c.Set("Connection", "keep-alive")
    c.Set("Transfer-Encoding", "chunked")

    c.Context().SetBodyStreamWriter(fasthttp.StreamWriter(func(w *bufio.Writer) {
      fmt.Println("WRITER")
      var i int

      for {
        i++
        msg := fmt.Sprintf("%d - the time is %v", i, time.Now())
        fmt.Fprintf(w, "data: Message: %s\n\n", msg)
        fmt.Println(msg)

        w.Flush()
        time.Sleep(5 * time.Second)
      }
    }))

    return nil
  })

  log.Fatal(app.Listen(":3000"))
}
```

### Recover middleware

📖 [Recover](https://docs.ximispot.io/api/middleware/recover)

```go
import (
    "github.com/ximispot/woody"
    "github.com/ximispot/woody/middleware/recover"
)

func main() {
    app := woody.New()

    app.Use(recover.New())

    app.Get("/", func(c *woody.Ctx) error {
        panic("normally this would crash your app")
    })

    log.Fatal(app.Listen(":3000"))
}
```

</details>

### Using Trusted Proxy

📖 [Config](https://docs.ximispot.io/api/woody#config)

```go
import (
    "github.com/ximispot/woody"
    "github.com/ximispot/woody/middleware/recover"
)

func main() {
    app := woody.New(woody.Config{
        EnableTrustedProxyCheck: true,
        TrustedProxies: []string{"0.0.0.0", "1.1.1.1/30"}, // IP address or IP address range
        ProxyHeader: woody.HeaderXForwardedFor},
    })

    // ...

    log.Fatal(app.Listen(":3000"))
}
```

</details>


## License

Check the License.md file