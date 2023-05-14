//nolint:bodyclose // Much easier to just ignore memory leaks in tests
package redirect

import (
	"context"
	"net/http"
	"testing"

	"github.com/ximispot/woody"
	"github.com/ximispot/woody/utils"
)

func Test_Redirect(t *testing.T) {
	app := *woody.New()

	app.Use(New(Config{
		Rules: map[string]string{
			"/default": "google.com",
		},
		StatusCode: woody.StatusMovedPermanently,
	}))
	app.Use(New(Config{
		Rules: map[string]string{
			"/default/*": "woody.wiki",
		},
		StatusCode: woody.StatusTemporaryRedirect,
	}))
	app.Use(New(Config{
		Rules: map[string]string{
			"/redirect/*": "$1",
		},
		StatusCode: woody.StatusSeeOther,
	}))
	app.Use(New(Config{
		Rules: map[string]string{
			"/pattern/*": "golang.org",
		},
		StatusCode: woody.StatusFound,
	}))

	app.Use(New(Config{
		Rules: map[string]string{
			"/": "/swagger",
		},
		StatusCode: woody.StatusMovedPermanently,
	}))

	app.Get("/api/*", func(c *woody.Ctx) error {
		return c.SendString("API")
	})

	app.Get("/new", func(c *woody.Ctx) error {
		return c.SendString("Hello, World!")
	})

	tests := []struct {
		name       string
		url        string
		redirectTo string
		statusCode int
	}{
		{
			name:       "should be returns status StatusFound without a wildcard",
			url:        "/default",
			redirectTo: "google.com",
			statusCode: woody.StatusMovedPermanently,
		},
		{
			name:       "should be returns status StatusTemporaryRedirect  using wildcard",
			url:        "/default/xyz",
			redirectTo: "woody.wiki",
			statusCode: woody.StatusTemporaryRedirect,
		},
		{
			name:       "should be returns status StatusSeeOther without set redirectTo to use the default",
			url:        "/redirect/github.com/ximispot/redirect",
			redirectTo: "github.com/ximispot/redirect",
			statusCode: woody.StatusSeeOther,
		},
		{
			name:       "should return the status code default",
			url:        "/pattern/xyz",
			redirectTo: "golang.org",
			statusCode: woody.StatusFound,
		},
		{
			name:       "access URL without rule",
			url:        "/new",
			statusCode: woody.StatusOK,
		},
		{
			name:       "redirect to swagger route",
			url:        "/",
			redirectTo: "/swagger",
			statusCode: woody.StatusMovedPermanently,
		},
		{
			name:       "no redirect to swagger route",
			url:        "/api/",
			statusCode: woody.StatusOK,
		},
		{
			name:       "no redirect to swagger route #2",
			url:        "/api/test",
			statusCode: woody.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequestWithContext(context.Background(), woody.MethodGet, tt.url, nil)
			utils.AssertEqual(t, err, nil)
			req.Header.Set("Location", "github.com/ximispot/redirect")
			resp, err := app.Test(req)

			utils.AssertEqual(t, err, nil)
			utils.AssertEqual(t, tt.statusCode, resp.StatusCode)
			utils.AssertEqual(t, tt.redirectTo, resp.Header.Get("Location"))
		})
	}
}

func Test_Next(t *testing.T) {
	// Case 1 : Next function always returns true
	app := *woody.New()
	app.Use(New(Config{
		Next: func(*woody.Ctx) bool {
			return true
		},
		Rules: map[string]string{
			"/default": "google.com",
		},
		StatusCode: woody.StatusMovedPermanently,
	}))

	app.Use(func(c *woody.Ctx) error {
		return c.SendStatus(woody.StatusOK)
	})

	req, err := http.NewRequestWithContext(context.Background(), woody.MethodGet, "/default", nil)
	utils.AssertEqual(t, err, nil)
	resp, err := app.Test(req)
	utils.AssertEqual(t, err, nil)

	utils.AssertEqual(t, woody.StatusOK, resp.StatusCode)

	// Case 2 : Next function always returns false
	app = *woody.New()
	app.Use(New(Config{
		Next: func(*woody.Ctx) bool {
			return false
		},
		Rules: map[string]string{
			"/default": "google.com",
		},
		StatusCode: woody.StatusMovedPermanently,
	}))

	req, err = http.NewRequestWithContext(context.Background(), woody.MethodGet, "/default", nil)
	utils.AssertEqual(t, err, nil)
	resp, err = app.Test(req)
	utils.AssertEqual(t, err, nil)

	utils.AssertEqual(t, woody.StatusMovedPermanently, resp.StatusCode)
	utils.AssertEqual(t, "google.com", resp.Header.Get("Location"))
}

func Test_NoRules(t *testing.T) {
	// Case 1: No rules with default route defined
	app := *woody.New()

	app.Use(New(Config{
		StatusCode: woody.StatusMovedPermanently,
	}))

	app.Use(func(c *woody.Ctx) error {
		return c.SendStatus(woody.StatusOK)
	})

	req, err := http.NewRequestWithContext(context.Background(), woody.MethodGet, "/default", nil)
	utils.AssertEqual(t, err, nil)
	resp, err := app.Test(req)
	utils.AssertEqual(t, err, nil)
	utils.AssertEqual(t, woody.StatusOK, resp.StatusCode)

	// Case 2: No rules and no default route defined
	app = *woody.New()

	app.Use(New(Config{
		StatusCode: woody.StatusMovedPermanently,
	}))

	req, err = http.NewRequestWithContext(context.Background(), woody.MethodGet, "/default", nil)
	utils.AssertEqual(t, err, nil)
	resp, err = app.Test(req)
	utils.AssertEqual(t, err, nil)
	utils.AssertEqual(t, woody.StatusNotFound, resp.StatusCode)
}

func Test_DefaultConfig(t *testing.T) {
	// Case 1: Default config and no default route
	app := *woody.New()

	app.Use(New())

	req, err := http.NewRequestWithContext(context.Background(), woody.MethodGet, "/default", nil)
	utils.AssertEqual(t, err, nil)
	resp, err := app.Test(req)

	utils.AssertEqual(t, err, nil)
	utils.AssertEqual(t, woody.StatusNotFound, resp.StatusCode)

	// Case 2: Default config and default route
	app = *woody.New()

	app.Use(New())
	app.Use(func(c *woody.Ctx) error {
		return c.SendStatus(woody.StatusOK)
	})

	req, err = http.NewRequestWithContext(context.Background(), woody.MethodGet, "/default", nil)
	utils.AssertEqual(t, err, nil)
	resp, err = app.Test(req)

	utils.AssertEqual(t, err, nil)
	utils.AssertEqual(t, woody.StatusOK, resp.StatusCode)
}

func Test_RegexRules(t *testing.T) {
	// Case 1: Rules regex is empty
	app := *woody.New()
	app.Use(New(Config{
		Rules:      map[string]string{},
		StatusCode: woody.StatusMovedPermanently,
	}))

	app.Use(func(c *woody.Ctx) error {
		return c.SendStatus(woody.StatusOK)
	})

	req, err := http.NewRequestWithContext(context.Background(), woody.MethodGet, "/default", nil)
	utils.AssertEqual(t, err, nil)
	resp, err := app.Test(req)

	utils.AssertEqual(t, err, nil)
	utils.AssertEqual(t, woody.StatusOK, resp.StatusCode)

	// Case 2: Rules regex map contains valid regex and well-formed replacement URLs
	app = *woody.New()
	app.Use(New(Config{
		Rules: map[string]string{
			"/default": "google.com",
		},
		StatusCode: woody.StatusMovedPermanently,
	}))

	app.Use(func(c *woody.Ctx) error {
		return c.SendStatus(woody.StatusOK)
	})

	req, err = http.NewRequestWithContext(context.Background(), woody.MethodGet, "/default", nil)
	utils.AssertEqual(t, err, nil)
	resp, err = app.Test(req)

	utils.AssertEqual(t, err, nil)
	utils.AssertEqual(t, woody.StatusMovedPermanently, resp.StatusCode)
	utils.AssertEqual(t, "google.com", resp.Header.Get("Location"))

	// Case 3: Test invalid regex throws panic
	defer func() {
		if r := recover(); r != nil {
			t.Log("Recovered from invalid regex: ", r)
		}
	}()

	app = *woody.New()
	app.Use(New(Config{
		Rules: map[string]string{
			"(": "google.com",
		},
		StatusCode: woody.StatusMovedPermanently,
	}))
	t.Error("Expected panic, got nil")
}
