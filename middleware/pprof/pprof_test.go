package pprof

import (
	"bytes"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/ximispot/woody"
	"github.com/ximispot/woody/utils"
)

func Test_Non_Pprof_Path(t *testing.T) {
	t.Parallel()
	app := woody.New(woody.Config{DisableStartupMessage: true})

	app.Use(New())

	app.Get("/", func(c *woody.Ctx) error {
		return c.SendString("escaped")
	})

	resp, err := app.Test(httptest.NewRequest(woody.MethodGet, "/", nil))
	utils.AssertEqual(t, nil, err)
	utils.AssertEqual(t, 200, resp.StatusCode)

	b, err := io.ReadAll(resp.Body)
	utils.AssertEqual(t, nil, err)
	utils.AssertEqual(t, "escaped", string(b))
}

func Test_Non_Pprof_Path_WithPrefix(t *testing.T) {
	t.Parallel()
	app := woody.New(woody.Config{DisableStartupMessage: true})

	app.Use(New(Config{Prefix: "/federated-woody"}))

	app.Get("/", func(c *woody.Ctx) error {
		return c.SendString("escaped")
	})

	resp, err := app.Test(httptest.NewRequest(woody.MethodGet, "/", nil))
	utils.AssertEqual(t, nil, err)
	utils.AssertEqual(t, 200, resp.StatusCode)

	b, err := io.ReadAll(resp.Body)
	utils.AssertEqual(t, nil, err)
	utils.AssertEqual(t, "escaped", string(b))
}

func Test_Pprof_Index(t *testing.T) {
	t.Parallel()
	app := woody.New(woody.Config{DisableStartupMessage: true})

	app.Use(New())

	app.Get("/", func(c *woody.Ctx) error {
		return c.SendString("escaped")
	})

	resp, err := app.Test(httptest.NewRequest(woody.MethodGet, "/debug/pprof/", nil))
	utils.AssertEqual(t, nil, err)
	utils.AssertEqual(t, 200, resp.StatusCode)
	utils.AssertEqual(t, woody.MIMETextHTMLCharsetUTF8, resp.Header.Get(woody.HeaderContentType))

	b, err := io.ReadAll(resp.Body)
	utils.AssertEqual(t, nil, err)
	utils.AssertEqual(t, true, bytes.Contains(b, []byte("<title>/debug/pprof/</title>")))
}

func Test_Pprof_Index_WithPrefix(t *testing.T) {
	t.Parallel()
	app := woody.New(woody.Config{DisableStartupMessage: true})

	app.Use(New(Config{Prefix: "/federated-woody"}))

	app.Get("/", func(c *woody.Ctx) error {
		return c.SendString("escaped")
	})

	resp, err := app.Test(httptest.NewRequest(woody.MethodGet, "/federated-woody/debug/pprof/", nil))
	utils.AssertEqual(t, nil, err)
	utils.AssertEqual(t, 200, resp.StatusCode)
	utils.AssertEqual(t, woody.MIMETextHTMLCharsetUTF8, resp.Header.Get(woody.HeaderContentType))

	b, err := io.ReadAll(resp.Body)
	utils.AssertEqual(t, nil, err)
	utils.AssertEqual(t, true, bytes.Contains(b, []byte("<title>/debug/pprof/</title>")))
}

func Test_Pprof_Subs(t *testing.T) {
	t.Parallel()
	app := woody.New(woody.Config{DisableStartupMessage: true})

	app.Use(New())

	app.Get("/", func(c *woody.Ctx) error {
		return c.SendString("escaped")
	})

	subs := []string{
		"cmdline", "profile", "symbol", "trace", "allocs", "block",
		"goroutine", "heap", "mutex", "threadcreate",
	}

	for _, sub := range subs {
		t.Run(sub, func(t *testing.T) {
			t.Parallel()
			target := "/debug/pprof/" + sub
			if sub == "profile" {
				target += "?seconds=1"
			}
			resp, err := app.Test(httptest.NewRequest(woody.MethodGet, target, nil), 5000)
			utils.AssertEqual(t, nil, err)
			utils.AssertEqual(t, 200, resp.StatusCode)
		})
	}
}

func Test_Pprof_Subs_WithPrefix(t *testing.T) {
	t.Parallel()
	app := woody.New(woody.Config{DisableStartupMessage: true})

	app.Use(New(Config{Prefix: "/federated-woody"}))

	app.Get("/", func(c *woody.Ctx) error {
		return c.SendString("escaped")
	})

	subs := []string{
		"cmdline", "profile", "symbol", "trace", "allocs", "block",
		"goroutine", "heap", "mutex", "threadcreate",
	}

	for _, sub := range subs {
		t.Run(sub, func(t *testing.T) {
			t.Parallel()
			target := "/federated-woody/debug/pprof/" + sub
			if sub == "profile" {
				target += "?seconds=1"
			}
			resp, err := app.Test(httptest.NewRequest(woody.MethodGet, target, nil), 5000)
			utils.AssertEqual(t, nil, err)
			utils.AssertEqual(t, 200, resp.StatusCode)
		})
	}
}

func Test_Pprof_Other(t *testing.T) {
	t.Parallel()
	app := woody.New(woody.Config{DisableStartupMessage: true})

	app.Use(New())

	app.Get("/", func(c *woody.Ctx) error {
		return c.SendString("escaped")
	})

	resp, err := app.Test(httptest.NewRequest(woody.MethodGet, "/debug/pprof/302", nil))
	utils.AssertEqual(t, nil, err)
	utils.AssertEqual(t, 302, resp.StatusCode)
}

func Test_Pprof_Other_WithPrefix(t *testing.T) {
	t.Parallel()
	app := woody.New(woody.Config{DisableStartupMessage: true})

	app.Use(New(Config{Prefix: "/federated-woody"}))

	app.Get("/", func(c *woody.Ctx) error {
		return c.SendString("escaped")
	})

	resp, err := app.Test(httptest.NewRequest(woody.MethodGet, "/federated-woody/debug/pprof/302", nil))
	utils.AssertEqual(t, nil, err)
	utils.AssertEqual(t, 302, resp.StatusCode)
}

// go test -run Test_Pprof_Next
func Test_Pprof_Next(t *testing.T) {
	t.Parallel()
	app := woody.New()

	app.Use(New(Config{
		Next: func(_ *woody.Ctx) bool {
			return true
		},
	}))

	resp, err := app.Test(httptest.NewRequest(woody.MethodGet, "/debug/pprof/", nil))
	utils.AssertEqual(t, nil, err)
	utils.AssertEqual(t, 404, resp.StatusCode)
}

// go test -run Test_Pprof_Next_WithPrefix
func Test_Pprof_Next_WithPrefix(t *testing.T) {
	t.Parallel()
	app := woody.New()

	app.Use(New(Config{
		Next: func(_ *woody.Ctx) bool {
			return true
		},
		Prefix: "/federated-woody",
	}))

	resp, err := app.Test(httptest.NewRequest(woody.MethodGet, "/federated-woody/debug/pprof/", nil))
	utils.AssertEqual(t, nil, err)
	utils.AssertEqual(t, 404, resp.StatusCode)
}
