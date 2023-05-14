---
id: templates
title: 📝 Templates
description: Woody supports server-side template engines.
sidebar_position: 3
---

import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';

## Template interfaces

Woody provides a Views interface to provide your own template engine:

<Tabs>
<TabItem value="views" label="Views">

```go
type Views interface {
    Load() error
    Render(io.Writer, string, interface{}, ...string) error
}
```
</TabItem>
</Tabs>

`Views` interface contains a `Load` and `Render` method, `Load` is executed by Woody on app initialization to load/parse the templates.

```go
// Pass engine to Woody's Views Engine
app := woody.New(woody.Config{
    Views: engine,
    // Views Layout is the global layout for all template render until override on Render function.
    ViewsLayout: "layouts/main"
})
```

The `Render` method is linked to the [**ctx.Render\(\)**](../api/ctx.md#render) function that accepts a template name and binding data. It will use global layout if layout is not being defined in `Render` function.
If the Woody config option `PassLocalsToViews` is enabled, then all locals set using `ctx.Locals(key, value)` will be passed to the template.

```go
app.Get("/", func(c *woody.Ctx) error {
    return c.Render("index", woody.Map{
        "hello": "world",
    });
})
```

## Engines

Woody team maintains [templates](https://github.com/gowoody/template) package that provides wrappers for multiple template engines:

* [html](https://github.com/gowoody/template/tree/master/html)
* [ace](https://github.com/gowoody/template/tree/master/ace)
* [amber](https://github.com/gowoody/template/tree/master/amber)
* [django](https://github.com/gowoody/template/tree/master/django)
* [handlebars](https://github.com/gowoody/template/tree/master/handlebars)
* [jet](https://github.com/gowoody/template/tree/master/jet)
* [mustache](https://github.com/gowoody/template/tree/master/mustache)
* [pug](https://github.com/gowoody/template/tree/master/pug)

<Tabs>
<TabItem value="example" label="Example">

```go
package main

import (
    "log"
    "github.com/gowoody/woody/v2"
    "github.com/gowoody/template/html"
)

func main() {
    // Initialize standard Go html template engine
    engine := html.New("./views", ".html")
    // If you want other engine, just replace with following
    // Create a new engine with django
	// engine := django.New("./views", ".django")

    app := woody.New(woody.Config{
        Views: engine,
    })
    app.Get("/", func(c *woody.Ctx) error {
        // Render index template
        return c.Render("index", woody.Map{
            "Title": "Hello, World!",
        })
    })

    log.Fatal(app.Listen(":3000"))
}
```
</TabItem>
<TabItem value="index" label="views/index.html">

```markup
<!DOCTYPE html>
<body>
    <h1>{{.Title}}</h1>
</body>
</html>
```
</TabItem>
</Tabs>
