package session

import (
	"fmt"
	"testing"

	"github.com/ximispot/woody"
	"github.com/ximispot/woody/utils"

	"github.com/valyala/fasthttp"
)

// go test -run TestStore_getSessionID
func TestStore_getSessionID(t *testing.T) {
	t.Parallel()
	expectedID := "test-session-id"

	// woody instance
	app := woody.New()

	t.Run("from cookie", func(t *testing.T) {
		t.Parallel()
		// session store
		store := New()
		// woody context
		ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
		defer app.ReleaseCtx(ctx)
		// set cookie
		ctx.Request().Header.SetCookie(store.sessionName, expectedID)

		utils.AssertEqual(t, expectedID, store.getSessionID(ctx))
	})

	t.Run("from header", func(t *testing.T) {
		t.Parallel()
		// session store
		store := New(Config{
			KeyLookup: "header:session_id",
		})
		// woody context
		ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
		defer app.ReleaseCtx(ctx)
		// set header
		ctx.Request().Header.Set(store.sessionName, expectedID)

		utils.AssertEqual(t, expectedID, store.getSessionID(ctx))
	})

	t.Run("from url query", func(t *testing.T) {
		t.Parallel()
		// session store
		store := New(Config{
			KeyLookup: "query:session_id",
		})
		// woody context
		ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
		defer app.ReleaseCtx(ctx)
		// set url parameter
		ctx.Request().SetRequestURI(fmt.Sprintf("/path?%s=%s", store.sessionName, expectedID))

		utils.AssertEqual(t, expectedID, store.getSessionID(ctx))
	})
}

// go test -run TestStore_Get
// Regression: https://github.com/gowoody/woody/issues/1408
func TestStore_Get(t *testing.T) {
	t.Parallel()
	unexpectedID := "test-session-id"
	// woody instance
	app := woody.New()
	t.Run("session should persisted even session is invalid", func(t *testing.T) {
		t.Parallel()
		// session store
		store := New()
		// woody context
		ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
		defer app.ReleaseCtx(ctx)
		// set cookie
		ctx.Request().Header.SetCookie(store.sessionName, unexpectedID)

		acquiredSession, err := store.Get(ctx)
		utils.AssertEqual(t, err, nil)

		if acquiredSession.ID() != unexpectedID {
			t.Fatal("server should not accept the unexpectedID which is not in the store")
		}
	})
}
