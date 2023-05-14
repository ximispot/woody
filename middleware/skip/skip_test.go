package skip_test

import (
	"net/http/httptest"
	"testing"

	"github.com/ximispot/woody"
	"github.com/ximispot/woody/middleware/skip"
	"github.com/ximispot/woody/utils"
)

// go test -run Test_Skip
func Test_Skip(t *testing.T) {
	t.Parallel()
	app := woody.New()

	app.Use(skip.New(errTeapotHandler, func(*woody.Ctx) bool { return true }))
	app.Get("/", helloWorldHandler)

	resp, err := app.Test(httptest.NewRequest(woody.MethodGet, "/", nil))
	utils.AssertEqual(t, nil, err)
	utils.AssertEqual(t, woody.StatusOK, resp.StatusCode)
}

// go test -run Test_SkipFalse
func Test_SkipFalse(t *testing.T) {
	t.Parallel()
	app := woody.New()

	app.Use(skip.New(errTeapotHandler, func(*woody.Ctx) bool { return false }))
	app.Get("/", helloWorldHandler)

	resp, err := app.Test(httptest.NewRequest(woody.MethodGet, "/", nil))
	utils.AssertEqual(t, nil, err)
	utils.AssertEqual(t, woody.StatusTeapot, resp.StatusCode)
}

// go test -run Test_SkipNilFunc
func Test_SkipNilFunc(t *testing.T) {
	t.Parallel()
	app := woody.New()

	app.Use(skip.New(errTeapotHandler, nil))
	app.Get("/", helloWorldHandler)

	resp, err := app.Test(httptest.NewRequest(woody.MethodGet, "/", nil))
	utils.AssertEqual(t, nil, err)
	utils.AssertEqual(t, woody.StatusTeapot, resp.StatusCode)
}

func helloWorldHandler(c *woody.Ctx) error {
	return c.SendString("Hello, World 👋!")
}

func errTeapotHandler(*woody.Ctx) error {
	return woody.ErrTeapot
}
