package cors

import (
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ximispot/woody"
	"github.com/ximispot/woody/utils"

	"github.com/valyala/fasthttp"
)

func Test_CORS_Defaults(t *testing.T) {
	t.Parallel()
	app := woody.New()
	app.Use(New())

	testDefaultOrEmptyConfig(t, app)
}

func Test_CORS_Empty_Config(t *testing.T) {
	t.Parallel()
	app := woody.New()
	app.Use(New(Config{}))

	testDefaultOrEmptyConfig(t, app)
}

func testDefaultOrEmptyConfig(t *testing.T, app *woody.App) {
	t.Helper()

	h := app.Handler()

	// Test default GET response headers
	ctx := &fasthttp.RequestCtx{}
	ctx.Request.Header.SetMethod(woody.MethodGet)
	h(ctx)

	utils.AssertEqual(t, "*", string(ctx.Response.Header.Peek(woody.HeaderAccessControlAllowOrigin)))
	utils.AssertEqual(t, "", string(ctx.Response.Header.Peek(woody.HeaderAccessControlAllowCredentials)))
	utils.AssertEqual(t, "", string(ctx.Response.Header.Peek(woody.HeaderAccessControlExposeHeaders)))

	// Test default OPTIONS (preflight) response headers
	ctx = &fasthttp.RequestCtx{}
	ctx.Request.Header.SetMethod(woody.MethodOptions)
	h(ctx)

	utils.AssertEqual(t, "GET,POST,HEAD,PUT,DELETE,PATCH", string(ctx.Response.Header.Peek(woody.HeaderAccessControlAllowMethods)))
	utils.AssertEqual(t, "", string(ctx.Response.Header.Peek(woody.HeaderAccessControlAllowHeaders)))
	utils.AssertEqual(t, "", string(ctx.Response.Header.Peek(woody.HeaderAccessControlMaxAge)))
}

// go test -run -v Test_CORS_Wildcard
func Test_CORS_Wildcard(t *testing.T) {
	t.Parallel()
	// New woody instance
	app := woody.New()
	// OPTIONS (preflight) response headers when AllowOrigins is *
	app.Use(New(Config{
		AllowOrigins:     "*",
		AllowCredentials: true,
		MaxAge:           3600,
		ExposeHeaders:    "X-Request-ID",
		AllowHeaders:     "Authentication",
	}))
	// Get handler pointer
	handler := app.Handler()

	// Make request
	ctx := &fasthttp.RequestCtx{}
	ctx.Request.SetRequestURI("/")
	ctx.Request.Header.Set(woody.HeaderOrigin, "localhost")
	ctx.Request.Header.SetMethod(woody.MethodOptions)

	// Perform request
	handler(ctx)

	// Check result
	utils.AssertEqual(t, "*", string(ctx.Response.Header.Peek(woody.HeaderAccessControlAllowOrigin)))
	utils.AssertEqual(t, "true", string(ctx.Response.Header.Peek(woody.HeaderAccessControlAllowCredentials)))
	utils.AssertEqual(t, "3600", string(ctx.Response.Header.Peek(woody.HeaderAccessControlMaxAge)))
	utils.AssertEqual(t, "Authentication", string(ctx.Response.Header.Peek(woody.HeaderAccessControlAllowHeaders)))

	// Test non OPTIONS (preflight) response headers
	ctx = &fasthttp.RequestCtx{}
	ctx.Request.Header.SetMethod(woody.MethodGet)
	handler(ctx)

	utils.AssertEqual(t, "true", string(ctx.Response.Header.Peek(woody.HeaderAccessControlAllowCredentials)))
	utils.AssertEqual(t, "X-Request-ID", string(ctx.Response.Header.Peek(woody.HeaderAccessControlExposeHeaders)))
}

// go test -run -v Test_CORS_Subdomain
func Test_CORS_Subdomain(t *testing.T) {
	t.Parallel()
	// New woody instance
	app := woody.New()
	// OPTIONS (preflight) response headers when AllowOrigins is set to a subdomain
	app.Use("/", New(Config{AllowOrigins: "http://*.example.com"}))

	// Get handler pointer
	handler := app.Handler()

	// Make request with disallowed origin
	ctx := &fasthttp.RequestCtx{}
	ctx.Request.SetRequestURI("/")
	ctx.Request.Header.SetMethod(woody.MethodOptions)
	ctx.Request.Header.Set(woody.HeaderOrigin, "http://google.com")

	// Perform request
	handler(ctx)

	// Allow-Origin header should be "" because http://google.com does not satisfy http://*.example.com
	utils.AssertEqual(t, "", string(ctx.Response.Header.Peek(woody.HeaderAccessControlAllowOrigin)))

	ctx.Request.Reset()
	ctx.Response.Reset()

	// Make request with allowed origin
	ctx.Request.SetRequestURI("/")
	ctx.Request.Header.SetMethod(woody.MethodOptions)
	ctx.Request.Header.Set(woody.HeaderOrigin, "http://test.example.com")

	handler(ctx)

	utils.AssertEqual(t, "http://test.example.com", string(ctx.Response.Header.Peek(woody.HeaderAccessControlAllowOrigin)))
}

func Test_CORS_AllowOriginScheme(t *testing.T) {
	t.Parallel()
	tests := []struct {
		reqOrigin, pattern string
		shouldAllowOrigin  bool
	}{
		{
			pattern:           "http://example.com",
			reqOrigin:         "http://example.com",
			shouldAllowOrigin: true,
		},
		{
			pattern:           "https://example.com",
			reqOrigin:         "https://example.com",
			shouldAllowOrigin: true,
		},
		{
			pattern:           "http://example.com",
			reqOrigin:         "https://example.com",
			shouldAllowOrigin: false,
		},
		{
			pattern:           "http://*.example.com",
			reqOrigin:         "http://aaa.example.com",
			shouldAllowOrigin: true,
		},
		{
			pattern:           "http://*.example.com",
			reqOrigin:         "http://bbb.aaa.example.com",
			shouldAllowOrigin: true,
		},
		{
			pattern:           "http://*.aaa.example.com",
			reqOrigin:         "http://bbb.aaa.example.com",
			shouldAllowOrigin: true,
		},
		{
			pattern:           "http://*.example.com:8080",
			reqOrigin:         "http://aaa.example.com:8080",
			shouldAllowOrigin: true,
		},
		{
			pattern:           "http://example.com",
			reqOrigin:         "http://gowoody.com",
			shouldAllowOrigin: false,
		},
		{
			pattern:           "http://*.aaa.example.com",
			reqOrigin:         "http://ccc.bbb.example.com",
			shouldAllowOrigin: false,
		},
		{
			pattern: "http://*.example.com",
			reqOrigin: `http://1234567890.1234567890.1234567890.1234567890.1234567890.1234567890.1234567890.1234567890.1234567890.1234567890\
		  .1234567890.1234567890.1234567890.1234567890.1234567890.1234567890.1234567890.1234567890.1234567890.1234567890\
		  .1234567890.1234567890.1234567890.1234567890.1234567890.1234567890.1234567890.1234567890.1234567890.1234567890\
			.1234567890.1234567890.1234567890.1234567890.1234567890.1234567890.1234567890.1234567890.1234567890.1234567890.example.com`,
			shouldAllowOrigin: false,
		},
		{
			pattern:           "http://example.com",
			reqOrigin:         "http://ccc.bbb.example.com",
			shouldAllowOrigin: false,
		},
		{
			pattern:           "https://*--aaa.bbb.com",
			reqOrigin:         "https://prod-preview--aaa.bbb.com",
			shouldAllowOrigin: false,
		},
		{
			pattern:           "http://*.example.com",
			reqOrigin:         "http://ccc.bbb.example.com",
			shouldAllowOrigin: true,
		},
		{
			pattern:           "http://foo.[a-z]*.example.com",
			reqOrigin:         "http://ccc.bbb.example.com",
			shouldAllowOrigin: false,
		},
	}

	for _, tt := range tests {
		app := woody.New()
		app.Use("/", New(Config{AllowOrigins: tt.pattern}))

		handler := app.Handler()

		ctx := &fasthttp.RequestCtx{}
		ctx.Request.SetRequestURI("/")
		ctx.Request.Header.SetMethod(woody.MethodOptions)
		ctx.Request.Header.Set(woody.HeaderOrigin, tt.reqOrigin)

		handler(ctx)

		if tt.shouldAllowOrigin {
			utils.AssertEqual(t, tt.reqOrigin, string(ctx.Response.Header.Peek(woody.HeaderAccessControlAllowOrigin)))
		} else {
			utils.AssertEqual(t, "", string(ctx.Response.Header.Peek(woody.HeaderAccessControlAllowOrigin)))
		}
	}
}

// go test -run Test_CORS_Next
func Test_CORS_Next(t *testing.T) {
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

func Test_CORS_AllowOriginsAndAllowOriginsFunc(t *testing.T) {
	t.Parallel()
	// New woody instance
	app := woody.New()
	app.Use("/", New(Config{
		AllowOrigins: "http://example-1.com",
		AllowOriginsFunc: func(origin string) bool {
			return strings.Contains(origin, "example-2")
		},
	}))

	// Get handler pointer
	handler := app.Handler()

	// Make request with disallowed origin
	ctx := &fasthttp.RequestCtx{}
	ctx.Request.SetRequestURI("/")
	ctx.Request.Header.SetMethod(woody.MethodOptions)
	ctx.Request.Header.Set(woody.HeaderOrigin, "http://google.com")

	// Perform request
	handler(ctx)

	// Allow-Origin header should be "" because http://google.com does not satisfy http://example-1.com or 'strings.Contains(origin, "example-2")'
	utils.AssertEqual(t, "", string(ctx.Response.Header.Peek(woody.HeaderAccessControlAllowOrigin)))

	ctx.Request.Reset()
	ctx.Response.Reset()

	// Make request with allowed origin
	ctx.Request.SetRequestURI("/")
	ctx.Request.Header.SetMethod(woody.MethodOptions)
	ctx.Request.Header.Set(woody.HeaderOrigin, "http://example-1.com")

	handler(ctx)

	utils.AssertEqual(t, "http://example-1.com", string(ctx.Response.Header.Peek(woody.HeaderAccessControlAllowOrigin)))

	ctx.Request.Reset()
	ctx.Response.Reset()

	// Make request with allowed origin
	ctx.Request.SetRequestURI("/")
	ctx.Request.Header.SetMethod(woody.MethodOptions)
	ctx.Request.Header.Set(woody.HeaderOrigin, "http://example-2.com")

	handler(ctx)

	utils.AssertEqual(t, "http://example-2.com", string(ctx.Response.Header.Peek(woody.HeaderAccessControlAllowOrigin)))
}

func Test_CORS_AllowOriginsFunc(t *testing.T) {
	t.Parallel()
	// New woody instance
	app := woody.New()
	app.Use("/", New(Config{
		AllowOriginsFunc: func(origin string) bool {
			return strings.Contains(origin, "example-2")
		},
	}))

	// Get handler pointer
	handler := app.Handler()

	// Make request with disallowed origin
	ctx := &fasthttp.RequestCtx{}
	ctx.Request.SetRequestURI("/")
	ctx.Request.Header.SetMethod(woody.MethodOptions)
	ctx.Request.Header.Set(woody.HeaderOrigin, "http://google.com")

	// Perform request
	handler(ctx)

	// Allow-Origin header should be "*" because http://google.com does not satisfy 'strings.Contains(origin, "example-2")'
	// and AllowOrigins has not been set so the default "*" is used
	utils.AssertEqual(t, "*", string(ctx.Response.Header.Peek(woody.HeaderAccessControlAllowOrigin)))

	ctx.Request.Reset()
	ctx.Response.Reset()

	// Make request with allowed origin
	ctx.Request.SetRequestURI("/")
	ctx.Request.Header.SetMethod(woody.MethodOptions)
	ctx.Request.Header.Set(woody.HeaderOrigin, "http://example-2.com")

	handler(ctx)

	// Allow-Origin header should be "http://example-2.com"
	utils.AssertEqual(t, "http://example-2.com", string(ctx.Response.Header.Peek(woody.HeaderAccessControlAllowOrigin)))
}
