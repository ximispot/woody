//nolint:bodyclose // Much easier to just ignore memory leaks in tests
package filesystem

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ximispot/woody"
	"github.com/ximispot/woody/utils"
)

// go test -run Test_FileSystem
func Test_FileSystem(t *testing.T) {
	t.Parallel()
	app := woody.New()

	app.Use("/test", New(Config{
		Root: http.Dir("../../.github/testdata/fs"),
	}))

	app.Use("/dir", New(Config{
		Root:   http.Dir("../../.github/testdata/fs"),
		Browse: true,
	}))

	app.Get("/", func(c *woody.Ctx) error {
		return c.SendString("Hello, World!")
	})

	app.Use("/spatest", New(Config{
		Root:         http.Dir("../../.github/testdata/fs"),
		Index:        "index.html",
		NotFoundFile: "index.html",
	}))

	app.Use("/prefix", New(Config{
		Root:       http.Dir("../../.github/testdata/fs"),
		PathPrefix: "img",
	}))

	tests := []struct {
		name         string
		url          string
		statusCode   int
		contentType  string
		modifiedTime string
	}{
		{
			name:        "Should be returns status 200 with suitable content-type",
			url:         "/test/index.html",
			statusCode:  200,
			contentType: "text/html",
		},
		{
			name:        "Should be returns status 200 with suitable content-type",
			url:         "/test",
			statusCode:  200,
			contentType: "text/html",
		},
		{
			name:        "Should be returns status 200 with suitable content-type",
			url:         "/test/css/style.css",
			statusCode:  200,
			contentType: "text/css",
		},
		{
			name:       "Should be returns status 404",
			url:        "/test/nofile.js",
			statusCode: 404,
		},
		{
			name:       "Should be returns status 404",
			url:        "/test/nofile",
			statusCode: 404,
		},
		{
			name:        "Should be returns status 200",
			url:         "/",
			statusCode:  200,
			contentType: "text/plain; charset=utf-8",
		},
		{
			name:       "Should be returns status 403",
			url:        "/test/img",
			statusCode: 403,
		},
		{
			name:        "Should list the directory contents",
			url:         "/dir/img",
			statusCode:  200,
			contentType: "text/html",
		},
		{
			name:        "Should list the directory contents",
			url:         "/dir/img/",
			statusCode:  200,
			contentType: "text/html",
		},
		{
			name:        "Should be returns status 200",
			url:         "/dir/img/woody.png",
			statusCode:  200,
			contentType: "image/png",
		},
		{
			name:        "Should be return status 200",
			url:         "/spatest/doesnotexist",
			statusCode:  200,
			contentType: "text/html",
		},
		{
			name:        "PathPrefix should be applied",
			url:         "/prefix/woody.png",
			statusCode:  200,
			contentType: "image/png",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			resp, err := app.Test(httptest.NewRequest(woody.MethodGet, tt.url, nil))
			utils.AssertEqual(t, nil, err)
			utils.AssertEqual(t, tt.statusCode, resp.StatusCode)

			if tt.contentType != "" {
				ct := resp.Header.Get("Content-Type")
				utils.AssertEqual(t, tt.contentType, ct)
			}
		})
	}
}

// go test -run Test_FileSystem_Next
func Test_FileSystem_Next(t *testing.T) {
	t.Parallel()
	app := woody.New()
	app.Use(New(Config{
		Root: http.Dir("../../.github/testdata/fs"),
		Next: func(_ *woody.Ctx) bool {
			return true
		},
	}))

	resp, err := app.Test(httptest.NewRequest(woody.MethodGet, "/", nil))
	utils.AssertEqual(t, nil, err)
	utils.AssertEqual(t, woody.StatusNotFound, resp.StatusCode)
}

func Test_FileSystem_NonGetAndHead(t *testing.T) {
	t.Parallel()
	app := woody.New()

	app.Use("/test", New(Config{
		Root: http.Dir("../../.github/testdata/fs"),
	}))

	resp, err := app.Test(httptest.NewRequest(woody.MethodPost, "/test", nil))
	utils.AssertEqual(t, nil, err)
	utils.AssertEqual(t, 404, resp.StatusCode)
}

func Test_FileSystem_Head(t *testing.T) {
	t.Parallel()
	app := woody.New()

	app.Use("/test", New(Config{
		Root: http.Dir("../../.github/testdata/fs"),
	}))

	req, err := http.NewRequestWithContext(context.Background(), woody.MethodHead, "/test", nil)
	utils.AssertEqual(t, nil, err)
	resp, err := app.Test(req)
	utils.AssertEqual(t, nil, err)
	utils.AssertEqual(t, 200, resp.StatusCode)
}

func Test_FileSystem_NoRoot(t *testing.T) {
	t.Parallel()
	defer func() {
		utils.AssertEqual(t, "filesystem: Root cannot be nil", recover())
	}()

	app := woody.New()
	app.Use(New())
	_, err := app.Test(httptest.NewRequest(woody.MethodGet, "/", nil))
	utils.AssertEqual(t, nil, err)
}

func Test_FileSystem_UsingParam(t *testing.T) {
	t.Parallel()
	app := woody.New()

	app.Use("/:path", func(c *woody.Ctx) error {
		return SendFile(c, http.Dir("../../.github/testdata/fs"), c.Params("path")+".html")
	})

	req, err := http.NewRequestWithContext(context.Background(), woody.MethodHead, "/index", nil)
	utils.AssertEqual(t, nil, err)
	resp, err := app.Test(req)
	utils.AssertEqual(t, nil, err)
	utils.AssertEqual(t, 200, resp.StatusCode)
}

func Test_FileSystem_UsingParam_NonFile(t *testing.T) {
	t.Parallel()
	app := woody.New()

	app.Use("/:path", func(c *woody.Ctx) error {
		return SendFile(c, http.Dir("../../.github/testdata/fs"), c.Params("path")+".html")
	})

	req, err := http.NewRequestWithContext(context.Background(), woody.MethodHead, "/template", nil)
	utils.AssertEqual(t, nil, err)
	resp, err := app.Test(req)
	utils.AssertEqual(t, nil, err)
	utils.AssertEqual(t, 404, resp.StatusCode)
}

func Test_FileSystem_UsingContentTypeCharset(t *testing.T) {
	t.Parallel()
	app := woody.New()
	app.Use(New(Config{
		Root:               http.Dir("../../.github/testdata/fs/index.html"),
		ContentTypeCharset: "UTF-8",
	}))

	resp, err := app.Test(httptest.NewRequest(woody.MethodGet, "/", nil))
	utils.AssertEqual(t, nil, err)
	utils.AssertEqual(t, 200, resp.StatusCode)
	utils.AssertEqual(t, "text/html; charset=UTF-8", resp.Header.Get("Content-Type"))
}
