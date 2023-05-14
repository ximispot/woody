package requestid

import (
	"net/http/httptest"
	"testing"

	"github.com/ximispot/woody"
	"github.com/ximispot/woody/utils"
)

// go test -run Test_RequestID
func Test_RequestID(t *testing.T) {
	t.Parallel()
	app := woody.New()

	app.Use(New())

	app.Get("/", func(c *woody.Ctx) error {
		return c.SendString("Hello, World 👋!")
	})

	resp, err := app.Test(httptest.NewRequest(woody.MethodGet, "/", nil))
	utils.AssertEqual(t, nil, err)
	utils.AssertEqual(t, woody.StatusOK, resp.StatusCode)

	reqid := resp.Header.Get(woody.HeaderXRequestID)
	utils.AssertEqual(t, 36, len(reqid))

	req := httptest.NewRequest(woody.MethodGet, "/", nil)
	req.Header.Add(woody.HeaderXRequestID, reqid)

	resp, err = app.Test(req)
	utils.AssertEqual(t, nil, err)
	utils.AssertEqual(t, woody.StatusOK, resp.StatusCode)
	utils.AssertEqual(t, reqid, resp.Header.Get(woody.HeaderXRequestID))
}

// go test -run Test_RequestID_Next
func Test_RequestID_Next(t *testing.T) {
	t.Parallel()
	app := woody.New()
	app.Use(New(Config{
		Next: func(_ *woody.Ctx) bool {
			return true
		},
	}))

	resp, err := app.Test(httptest.NewRequest(woody.MethodGet, "/", nil))
	utils.AssertEqual(t, nil, err)
	utils.AssertEqual(t, resp.Header.Get(woody.HeaderXRequestID), "")
	utils.AssertEqual(t, woody.StatusNotFound, resp.StatusCode)
}

// go test -run Test_RequestID_Locals
func Test_RequestID_Locals(t *testing.T) {
	t.Parallel()
	reqID := "ThisIsARequestId"
	type ContextKey int
	const requestContextKey ContextKey = iota

	app := woody.New()
	app.Use(New(Config{
		Generator: func() string {
			return reqID
		},
		ContextKey: requestContextKey,
	}))

	var ctxVal string

	app.Use(func(c *woody.Ctx) error {
		ctxVal = c.Locals(requestContextKey).(string) //nolint:forcetypeassert,errcheck // We always store a string in here
		return c.Next()
	})

	_, err := app.Test(httptest.NewRequest(woody.MethodGet, "/", nil))
	utils.AssertEqual(t, nil, err)
	utils.AssertEqual(t, reqID, ctxVal)
}
