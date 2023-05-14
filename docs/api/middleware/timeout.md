---
id: timeout
title: Timeout
---

There exist two distinct implementations of timeout middleware [Woody](https://github.com/gowoody/woody).

**New**

Wraps a `woody.Handler` with a timeout. If the handler takes longer than the given duration to return, the timeout error is set and forwarded to the centralized [ErrorHandler](https://docs.gowoody.io/error-handling).

:::caution
This has been deprecated since it raises race conditions.
:::

**NewWithContext**

As a `woody.Handler` wrapper, it creates a context with `context.WithTimeout` and pass it in `UserContext`. 
 
If the context passed executions (eg. DB ops, Http calls) takes longer than the given duration to return, the timeout error is set and forwarded to the centralized `ErrorHandler`.


It does not cancel long running executions. Underlying executions must handle timeout by using `context.Context` parameter.

## Signatures

```go
func New(handler woody.Handler, timeout time.Duration, timeoutErrors ...error) woody.Handler
func NewWithContext(handler woody.Handler, timeout time.Duration, timeoutErrors ...error) woody.Handler
```

## Examples

Import the middleware package that is part of the Woody web framework

```go
import (
  "github.com/gowoody/woody/v2"
  "github.com/ximispot/woody/middleware/timeout"
)
```

After you initiate your Woody app, you can use the following possibilities:

```go
func main() {
	app := woody.New()

	h := func(c *woody.Ctx) error {
		sleepTime, _ := time.ParseDuration(c.Params("sleepTime") + "ms")
		if err := sleepWithContext(c.UserContext(), sleepTime); err != nil {
			return fmt.Errorf("%w: execution error", err)
		}
		return nil
	}

	app.Get("/foo/:sleepTime", timeout.New(h, 2*time.Second))
	log.Fatal(app.Listen(":3000"))
}

func sleepWithContext(ctx context.Context, d time.Duration) error {
	timer := time.NewTimer(d)

	select {
	case <-ctx.Done():
		if !timer.Stop() {
			<-timer.C
		}
		return context.DeadlineExceeded
	case <-timer.C:
	}
	return nil
}
```

Test http 200 with curl:

```bash
curl --location -I --request GET 'http://localhost:3000/foo/1000' 
```

Test http 408 with curl:

```bash
curl --location -I --request GET 'http://localhost:3000/foo/3000' 
```

Use with custom error:

```go
var ErrFooTimeOut = errors.New("foo context canceled")

func main() {
	app := woody.New()
	h := func(c *woody.Ctx) error {
		sleepTime, _ := time.ParseDuration(c.Params("sleepTime") + "ms")
		if err := sleepWithContextWithCustomError(c.UserContext(), sleepTime); err != nil {
			return fmt.Errorf("%w: execution error", err)
		}
		return nil
	}

	app.Get("/foo/:sleepTime", timeout.NewWithContext(h, 2*time.Second, ErrFooTimeOut))
	log.Fatal(app.Listen(":3000"))
}

func sleepWithContextWithCustomError(ctx context.Context, d time.Duration) error {
	timer := time.NewTimer(d)
	select {
	case <-ctx.Done():
		if !timer.Stop() {
			<-timer.C
		}
		return ErrFooTimeOut
	case <-timer.C:
	}
	return nil
}
```

Sample usage with a DB call:

```go
func main() {
	app := woody.New()
	db, _ := gorm.Open(postgres.Open("postgres://localhost/foodb"), &gorm.Config{})

	handler := func(ctx *woody.Ctx) error {
		tran := db.WithContext(ctx.UserContext()).Begin()
		
		if tran = tran.Exec("SELECT pg_sleep(50)"); tran.Error != nil {
			return tran.Error
		}
		
		if tran = tran.Commit(); tran.Error != nil {
			return tran.Error
		}

		return nil
	}

	app.Get("/foo", timeout.NewWithContext(handler, 10*time.Second))
	log.Fatal(app.Listen(":3000"))
}
```
