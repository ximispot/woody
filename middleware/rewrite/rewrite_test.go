//nolint:bodyclose // Much easier to just ignore memory leaks in tests
package rewrite

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/ximispot/woody"
	"github.com/ximispot/woody/utils"
)

func Test_New(t *testing.T) {
	// Test with no config
	m := New()

	if m == nil {
		t.Error("Expected middleware to be returned, got nil")
	}

	// Test with config
	m = New(Config{
		Rules: map[string]string{
			"/old": "/new",
		},
	})

	if m == nil {
		t.Error("Expected middleware to be returned, got nil")
	}

	// Test with full config
	m = New(Config{
		Next: func(*woody.Ctx) bool {
			return true
		},
		Rules: map[string]string{
			"/old": "/new",
		},
	})

	if m == nil {
		t.Error("Expected middleware to be returned, got nil")
	}
}

func Test_Rewrite(t *testing.T) {
	// Case 1: Next function always returns true
	app := woody.New()
	app.Use(New(Config{
		Next: func(*woody.Ctx) bool {
			return true
		},
		Rules: map[string]string{
			"/old": "/new",
		},
	}))

	app.Get("/old", func(c *woody.Ctx) error {
		return c.SendString("Rewrite Successful")
	})

	req, err := http.NewRequestWithContext(context.Background(), woody.MethodGet, "/old", nil)
	utils.AssertEqual(t, err, nil)
	resp, err := app.Test(req)
	utils.AssertEqual(t, err, nil)
	body, err := io.ReadAll(resp.Body)
	utils.AssertEqual(t, err, nil)
	bodyString := string(body)

	utils.AssertEqual(t, err, nil)
	utils.AssertEqual(t, woody.StatusOK, resp.StatusCode)
	utils.AssertEqual(t, "Rewrite Successful", bodyString)

	// Case 2: Next function always returns false
	app = woody.New()
	app.Use(New(Config{
		Next: func(*woody.Ctx) bool {
			return false
		},
		Rules: map[string]string{
			"/old": "/new",
		},
	}))

	app.Get("/new", func(c *woody.Ctx) error {
		return c.SendString("Rewrite Successful")
	})

	req, err = http.NewRequestWithContext(context.Background(), woody.MethodGet, "/old", nil)
	utils.AssertEqual(t, err, nil)
	resp, err = app.Test(req)
	utils.AssertEqual(t, err, nil)
	body, err = io.ReadAll(resp.Body)
	utils.AssertEqual(t, err, nil)
	bodyString = string(body)

	utils.AssertEqual(t, err, nil)
	utils.AssertEqual(t, woody.StatusOK, resp.StatusCode)
	utils.AssertEqual(t, "Rewrite Successful", bodyString)

	// Case 3: check for captured tokens in rewrite rule
	app = woody.New()
	app.Use(New(Config{
		Rules: map[string]string{
			"/users/*/orders/*": "/user/$1/order/$2",
		},
	}))

	app.Get("/user/:userID/order/:orderID", func(c *woody.Ctx) error {
		return c.SendString(fmt.Sprintf("User ID: %s, Order ID: %s", c.Params("userID"), c.Params("orderID")))
	})

	req, err = http.NewRequestWithContext(context.Background(), woody.MethodGet, "/users/123/orders/456", nil)
	utils.AssertEqual(t, err, nil)
	resp, err = app.Test(req)
	utils.AssertEqual(t, err, nil)
	body, err = io.ReadAll(resp.Body)
	utils.AssertEqual(t, err, nil)
	bodyString = string(body)

	utils.AssertEqual(t, err, nil)
	utils.AssertEqual(t, woody.StatusOK, resp.StatusCode)
	utils.AssertEqual(t, "User ID: 123, Order ID: 456", bodyString)

	// Case 4: Send non-matching request, handled by default route
	app = woody.New()
	app.Use(New(Config{
		Rules: map[string]string{
			"/users/*/orders/*": "/user/$1/order/$2",
		},
	}))

	app.Get("/user/:userID/order/:orderID", func(c *woody.Ctx) error {
		return c.SendString(fmt.Sprintf("User ID: %s, Order ID: %s", c.Params("userID"), c.Params("orderID")))
	})

	app.Use(func(c *woody.Ctx) error {
		return c.SendStatus(woody.StatusOK)
	})

	req, err = http.NewRequestWithContext(context.Background(), woody.MethodGet, "/not-matching-any-rule", nil)
	utils.AssertEqual(t, err, nil)
	resp, err = app.Test(req)
	utils.AssertEqual(t, err, nil)
	body, err = io.ReadAll(resp.Body)
	utils.AssertEqual(t, err, nil)
	bodyString = string(body)

	utils.AssertEqual(t, err, nil)
	utils.AssertEqual(t, woody.StatusOK, resp.StatusCode)
	utils.AssertEqual(t, "OK", bodyString)

	// Case 4: Send non-matching request, with no default route
	app = woody.New()
	app.Use(New(Config{
		Rules: map[string]string{
			"/users/*/orders/*": "/user/$1/order/$2",
		},
	}))

	app.Get("/user/:userID/order/:orderID", func(c *woody.Ctx) error {
		return c.SendString(fmt.Sprintf("User ID: %s, Order ID: %s", c.Params("userID"), c.Params("orderID")))
	})

	req, err = http.NewRequestWithContext(context.Background(), woody.MethodGet, "/not-matching-any-rule", nil)
	utils.AssertEqual(t, err, nil)
	resp, err = app.Test(req)
	utils.AssertEqual(t, err, nil)
	utils.AssertEqual(t, woody.StatusNotFound, resp.StatusCode)
}
