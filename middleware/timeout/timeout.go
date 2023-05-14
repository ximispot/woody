package timeout

import (
	"context"
	"errors"
	"log"
	"sync"
	"time"

	"github.com/ximispot/woody"
)

var once sync.Once

// New wraps a handler and aborts the process of the handler if the timeout is reached.
//
// Deprecated: This implementation contains data race issues. Use NewWithContext instead.
// Find documentation and sample usage on https://docs.gowoody.io/api/middleware/timeout
func New(handler woody.Handler, timeout time.Duration) woody.Handler {
	once.Do(func() {
		log.Printf("[Warning] - [TIMEOUT] timeout contains data race issues, not ready for production!")
	})

	if timeout <= 0 {
		return handler
	}

	// logic is from fasthttp.TimeoutWithCodeHandler https://github.com/valyala/fasthttp/blob/master/server.go#L418
	return func(ctx *woody.Ctx) error {
		ch := make(chan struct{}, 1)

		go func() {
			defer func() {
				if err := recover(); err != nil {
					log.Printf("[Warning] - [TIMEOUT] recover error %v", err)
				}
			}()
			if err := handler(ctx); err != nil {
				log.Printf("[Warning] - [TIMEOUT] handler error %v", err)
			}
			ch <- struct{}{}
		}()

		select {
		case <-ch:
		case <-time.After(timeout):
			return woody.ErrRequestTimeout
		}

		return nil
	}
}

// NewWithContext implementation of timeout middleware. Set custom errors(context.DeadlineExceeded vs) for get woody.ErrRequestTimeout response.
func NewWithContext(h woody.Handler, t time.Duration, tErrs ...error) woody.Handler {
	return func(ctx *woody.Ctx) error {
		timeoutContext, cancel := context.WithTimeout(ctx.UserContext(), t)
		defer cancel()
		ctx.SetUserContext(timeoutContext)
		if err := h(ctx); err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				return woody.ErrRequestTimeout
			}
			for i := range tErrs {
				if errors.Is(err, tErrs[i]) {
					return woody.ErrRequestTimeout
				}
			}
			return err
		}
		return nil
	}
}
