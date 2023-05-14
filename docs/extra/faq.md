---
id: faq
title: 🤔 FAQ
description: >-
  List of frequently asked questions. Feel free to open an issue to add your
  question to this page.
sidebar_position: 1
---

## How should I structure my application?

There is no definitive answer to this question. The answer depends on the scale of your application and the team that is involved. To be as flexible as possible, Woody makes no assumptions in terms of structure.

Routes and other application-specific logic can live in as many files as you wish, in any directory structure you prefer. View the following examples for inspiration:

* [gowoody/boilerplate](https://github.com/gowoody/boilerplate)
* [thomasvvugt/woody-boilerplate](https://github.com/thomasvvugt/woody-boilerplate)
* [Youtube - Building a REST API using Gorm and Woody](https://www.youtube.com/watch?v=Iq2qT0fRhAA)
* [embedmode/woodyseed](https://github.com/embedmode/woodyseed)

## How do I handle custom 404 responses?

If you're using v2.32.0 or later, all you need to do is to implement a custom error handler. See below, or see a more detailed explanation at [Error Handling](../guide/error-handling.md#custom-error-handler). 

If you're using v2.31.0 or earlier, the error handler will not capture 404 errors. Instead, you need to add a middleware function at the very bottom of the stack \(below all other functions\) to handle a 404 response:

```go title="Example"
app.Use(func(c *woody.Ctx) error {
    return c.Status(woody.StatusNotFound).SendString("Sorry can't find that!")
})
```

## How can i use live reload ?

[Air](https://github.com/cosmtrek/air) is a handy tool that automatically restarts your Go applications whenever the source code changes, making your development process faster and more efficient.

To use Air in a Woody project, follow these steps:

1. Install Air by downloading the appropriate binary for your operating system from the GitHub release page or by building the tool directly from source.
2. Create a configuration file for Air in your project directory. This file can be named, for example, .air.toml or air.conf. Here's a sample configuration file that works with Woody:
```toml
# .air.toml
root = "."
tmp_dir = "tmp"
[build]
  cmd = "go build -o ./tmp/main ."
  bin = "./tmp/main"
  delay = 1000 # ms
  exclude_dir = ["assets", "tmp", "vendor"]
  include_ext = ["go", "tpl", "tmpl", "html"]
  exclude_regex = ["_test\\.go"]
```
3. Start your Woody application using Air by running the following command in the terminal:
```sh
air
```

As you make changes to your source code, Air will detect them and automatically restart the application.

A complete example demonstrating the use of Air with Woody can be found in the [Woody Recipes repository](https://github.com/gowoody/recipes/tree/master/air). This example shows how to configure and use Air in a Woody project to create an efficient development environment.


## How do I set up an error handler?

To override the default error handler, you can override the default when providing a [Config](../api/woody.md#config) when initiating a new [Woody instance](../api/woody.md#new).

```go title="Example"
app := woody.New(woody.Config{
    ErrorHandler: func(c *woody.Ctx, err error) error {
        return c.Status(woody.StatusInternalServerError).SendString(err.Error())
    },
})
```

We have a dedicated page explaining how error handling works in Woody, see [Error Handling](../guide/error-handling.md).

## Which template engines does Woody support?

Woody currently supports 8 template engines in our [gowoody/template](https://github.com/gowoody/template) middleware:

* [Ace](https://github.com/yosssi/ace)
* [Amber](https://github.com/eknkc/amber)
* [Django](https://github.com/flosch/pongo2)
* [Handlebars](https://github.com/aymerick/raymond)
* [HTML](https://pkg.go.dev/html/template/)
* [Jet](https://github.com/CloudyKit/jet)
* [Mustache](https://github.com/cbroglie/mustache)
* [Pug](https://github.com/Joker/jade)

To learn more about using Templates in Woody, see [Templates](../guide/templates.md).

## Does Woody have a community chat?

Yes, we have our own [Discord ](https://gowoody.io/discord)server, where we hang out. We have different rooms for every subject.  
If you have questions or just want to have a chat, feel free to join us via this **&gt;** [**invite link**](https://gowoody.io/discord) **&lt;**.

![](/img/support-discord.png)

## Does woody support sub domain routing ?

Yes we do, here are some examples: 
This example works v2
```go
package main

import (
	"log"

	"github.com/gowoody/woody/v2"
	"github.com/ximispot/woody/middleware/logger"
)

type Host struct {
	Woody *woody.App
}

func main() {
	// Hosts
	hosts := map[string]*Host{}
	//-----
	// API
	//-----
	api := woody.New()
	api.Use(logger.New(logger.Config{
		Format: "[${ip}]:${port} ${status} - ${method} ${path}\n",
	}))
	hosts["api.localhost:3000"] = &Host{api}
	api.Get("/", func(c *woody.Ctx) error {
		return c.SendString("API")
	})
	//------
	// Blog
	//------
	blog := woody.New()
	blog.Use(logger.New(logger.Config{
		Format: "[${ip}]:${port} ${status} - ${method} ${path}\n",
	}))
	hosts["blog.localhost:3000"] = &Host{blog}
	blog.Get("/", func(c *woody.Ctx) error {
		return c.SendString("Blog")
	})
	//---------
	// Website
	//---------
	site := woody.New()
	site.Use(logger.New(logger.Config{
		Format: "[${ip}]:${port} ${status} - ${method} ${path}\n",
	}))

	hosts["localhost:3000"] = &Host{site}
	site.Get("/", func(c *woody.Ctx) error {
		return c.SendString("Website")
	})
	// Server
	app := woody.New()
	app.Use(func(c *woody.Ctx) error {
		host := hosts[c.Hostname()]
		if host == nil {
			return c.SendStatus(woody.StatusNotFound)
		} else {
			host.Woody.Handler()(c.Context())
			return nil
		}
	})
	log.Fatal(app.Listen(":3000"))
}
```
If more information is needed, please refer to this issue [#750](https://github.com/gowoody/woody/issues/750)
