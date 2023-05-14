---
id: grouping
title: 🎭 Grouping
sidebar_position: 2
---

:::info
In general, the Group functionality in Woody behaves similarly to ExpressJS. Groups are declared virtually and all routes declared within the group are flattened into a single list with a prefix, which is then checked by the framework in the order it was declared. This means that the behavior of Group in Woody is identical to that of ExpressJS.
:::

## Paths

Like **Routing**, groups can also have paths that belong to a cluster.

```go
func main() {
  app := woody.New()

  api := app.Group("/api", middleware) // /api

  v1 := api.Group("/v1", middleware)   // /api/v1
  v1.Get("/list", handler)             // /api/v1/list
  v1.Get("/user", handler)             // /api/v1/user

  v2 := api.Group("/v2", middleware)   // /api/v2
  v2.Get("/list", handler)             // /api/v2/list
  v2.Get("/user", handler)             // /api/v2/user

  log.Fatal(app.Listen(":3000"))
}
```

A **Group** of paths can have an optional handler.

```go
func main() {
  app := woody.New()

  api := app.Group("/api")      // /api

  v1 := api.Group("/v1")        // /api/v1
  v1.Get("/list", handler)      // /api/v1/list
  v1.Get("/user", handler)      // /api/v1/user

  v2 := api.Group("/v2")        // /api/v2
  v2.Get("/list", handler)      // /api/v2/list
  v2.Get("/user", handler)      // /api/v2/user

  log.Fatal(app.Listen(":3000"))
}
```

:::caution
Running **/api**, **/v1** or **/v2** will result in **404** error, make sure you have the errors set.
:::

## Group Handlers

Group handlers can also be used as a routing path but they must have **Next** added to them so that the flow can continue.

```go
func main() {
    app := woody.New()

    handler := func(c *woody.Ctx) error {
        return c.SendStatus(woody.StatusOK)
    }
    api := app.Group("/api") // /api

    v1 := api.Group("/v1", func(c *woody.Ctx) error { // middleware for /api/v1
        c.Set("Version", "v1")
        return c.Next()
    })
    v1.Get("/list", handler) // /api/v1/list
    v1.Get("/user", handler) // /api/v1/user

    log.Fatal(app.Listen(":3000"))
}
```
