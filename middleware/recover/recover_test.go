package recover //nolint:predeclared // TODO: Rename to some non-builtin

import (
	"net/http/httptest"
	"testing"

	"github.com/ximispot/woody"
	"github.com/ximispot/woody/utils"
)

// go test -run Test_Recover
func Test_Recover(t *testing.T) {
	t.Parallel()
	app := woody.New(woody.Config{
		ErrorHandler: func(c *woody.Ctx, err error) error {
			utils.AssertEqual(t, "Hi, I'm an error!", err.Error())
			return c.SendStatus(woody.StatusTeapot)
		},
	})

	app.Use(New())

	app.Get("/panic", func(c *woody.Ctx) error {
		panic("Hi, I'm an error!")
	})

	resp, err := app.Test(httptest.NewRequest(woody.MethodGet, "/panic", nil))
	utils.AssertEqual(t, nil, err)
	utils.AssertEqual(t, woody.StatusTeapot, resp.StatusCode)
}

// go test -run Test_Recover_Next
func Test_Recover_Next(t *testing.T) {
	t.Parallel()
	app := woody.New()
	app.Use(New(Config{
		Next: func(_ *woody.Ctx) bool {
			return true
		},
	}))

	resp, err := app.Test(httptest.NewRequest(woody.MethodGet, "/", nil))
	utils.AssertEqual(t, nil, err)
	utils.AssertEqual(t, woody.StatusNotFound, resp.StatusCode)
}

func Test_Recover_EnableStackTrace(t *testing.T) {
	t.Parallel()
	app := woody.New()
	app.Use(New(Config{
		EnableStackTrace: true,
	}))

	app.Get("/panic", func(c *woody.Ctx) error {
		panic("Hi, I'm an error!")
	})

	resp, err := app.Test(httptest.NewRequest(woody.MethodGet, "/panic", nil))
	utils.AssertEqual(t, nil, err)
	utils.AssertEqual(t, woody.StatusInternalServerError, resp.StatusCode)
}
