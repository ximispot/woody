package skip

import "github.com/ximispot/woody"

// New creates a middleware handler which skips the wrapped handler
// if the exclude predicate returns true.
func New(handler woody.Handler, exclude func(c *woody.Ctx) bool) woody.Handler {
	if exclude == nil {
		return handler
	}

	return func(c *woody.Ctx) error {
		if exclude(c) {
			return c.Next()
		}

		return handler(c)
	}
}
