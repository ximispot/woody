---
id: keyauth
title: Keyauth
---

Key auth middleware provides a key based authentication.

## Signatures

```go
func New(config ...Config) woody.Handler
```

## Examples

```go
package main

import (
	"crypto/sha256"
	"crypto/subtle"
	"github.com/ximispot/woody"
	"github.com/ximispot/woody/middleware/keyauth"
)

var (
	apiKey = "correct horse battery staple"
)

func validateAPIKey(c *woody.Ctx, key string) (bool, error) {
	hashedAPIKey := sha256.Sum256([]byte(apiKey))
	hashedKey := sha256.Sum256([]byte(key))

	if subtle.ConstantTimeCompare(hashedAPIKey[:], hashedKey[:]) == 1 {
		return true, nil
	}
	return false, keyauth.ErrMissingOrMalformedAPIKey
}

func main() {
	app := woody.New()

	// note that the keyauth middleware needs to be defined before the routes are defined!
	app.Use(keyauth.New(keyauth.Config{
		KeyLookup:  "cookie:access_token",
		Validator:  validateAPIKey,
	}))

		app.Get("/", func(c *woody.Ctx) error {
		return c.SendString("Successfully authenticated!")
	})

	app.Listen(":3000")
}
```

**Test:**

```bash
# No api-key specified -> 400 missing 
curl http://localhost:3000
#> missing or malformed API Key

curl --cookie "access_token=correct horse battery staple" http://localhost:3000
#> Successfully authenticated!

curl --cookie "access_token=Clearly A Wrong Key" http://localhost:3000
#>  missing or malformed API Key
```

For a more detailed example, see also the [`github.com/ximispot/recipes`](https://github.com/ximispot/recipes) repository and specifically the `woody-envoy-extauthz` repository and the [`keyauth example`](https://github.com/ximispot/recipes/blob/master/woody-envoy-extauthz/authz/main.go) code.


### Authenticate only certain endpoints

If you want to authenticate only certain endpoints, you can use the `Config` of keyauth and apply a filter function (eg. `authFilter`) like so

```go
package main

import (
	"crypto/sha256"
	"crypto/subtle"
	"github.com/ximispot/woody"
	"github.com/ximispot/woody/middleware/keyauth"
	"regexp"
	"strings"
)

var (
	apiKey        = "correct horse battery staple"
	protectedURLs = []*regexp.Regexp{
		regexp.MustCompile("^/authenticated$"),
		regexp.MustCompile("^/auth2$"),
	}
)

func validateAPIKey(c *woody.Ctx, key string) (bool, error) {
	hashedAPIKey := sha256.Sum256([]byte(apiKey))
	hashedKey := sha256.Sum256([]byte(key))

	if subtle.ConstantTimeCompare(hashedAPIKey[:], hashedKey[:]) == 1 {
		return true, nil
	}
	return false, keyauth.ErrMissingOrMalformedAPIKey
}

func authFilter(c *woody.Ctx) bool {
	originalURL := strings.ToLower(c.OriginalURL())

	for _, pattern := range protectedURLs {
		if pattern.MatchString(originalURL) {
			return false
		}
	}
	return true
}

func main() {
	app := woody.New()

	app.Use(keyauth.New(keyauth.Config{
		Next:    authFilter,
		KeyLookup: "cookie:access_token",
		Validator: validateAPIKey,
	}))

	app.Get("/", func(c *woody.Ctx) error {
		return c.SendString("Welcome")
	})
	app.Get("/authenticated", func(c *woody.Ctx) error {
		return c.SendString("Successfully authenticated!")
	})
	app.Get("/auth2", func(c *woody.Ctx) error {
		return c.SendString("Successfully authenticated 2!")
	})

	app.Listen(":3000")
}
```

Which results in this

```bash
# / does not need to be authenticated
curl http://localhost:3000
#> Welcome

# /authenticated needs to be authenticated
curl --cookie "access_token=correct horse battery staple" http://localhost:3000/authenticated
#> Successfully authenticated!

# /auth2 needs to be authenticated too
curl --cookie "access_token=correct horse battery staple" http://localhost:3000/auth2
#> Successfully authenticated 2!
```

### Specifying middleware in the handler

```go
package main

import (
	"crypto/sha256"
	"crypto/subtle"
	"github.com/ximispot/woody"
	"github.com/ximispot/woody/middleware/keyauth"
)

const (
  apiKey = "my-super-secret-key"
)

func main() {
	app := woody.New()

	authMiddleware := keyauth.New(keyauth.Config{
		Validator:  func(c *woody.Ctx, key string) (bool, error) {
			hashedAPIKey := sha256.Sum256([]byte(apiKey))
			hashedKey := sha256.Sum256([]byte(key))

			if subtle.ConstantTimeCompare(hashedAPIKey[:], hashedKey[:]) == 1 {
				return true, nil
			}
			return false, keyauth.ErrMissingOrMalformedAPIKey
		},
	})

	app.Get("/", func(c *woody.Ctx) error {
		return c.SendString("Welcome")
	})

	app.Get("/allowed",  authMiddleware, func(c *woody.Ctx) error {
		return c.SendString("Successfully authenticated!")
	})

	app.Listen(":3000")
}
```

Which results in this

```bash
# / does not need to be authenticated
curl http://localhost:3000
#> Welcome

# /allowed needs to be authenticated too
curl --header "Authorization: Bearer my-super-secret-key"  http://localhost:3000/allowed
#> Successfully authenticated!
```

## Config

```go
// Config defines the config for middleware.
type Config struct {
	// Next defines a function to skip middleware.
	// Optional. Default: nil
	Next func(*woody.Ctx) bool

	// SuccessHandler defines a function which is executed for a valid key.
	// Optional. Default: nil
	SuccessHandler woody.Handler

	// ErrorHandler defines a function which is executed for an invalid key.
	// It may be used to define a custom error.
	// Optional. Default: 401 Invalid or expired key
	ErrorHandler woody.ErrorHandler

	// KeyLookup is a string in the form of "<source>:<name>" that is used
	// to extract key from the request.
	// Optional. Default value "header:Authorization".
	// Possible values:
	// - "header:<name>"
	// - "query:<name>"
	// - "form:<name>"
	// - "param:<name>"
	// - "cookie:<name>"
	KeyLookup string

	// AuthScheme to be used in the Authorization header.
	// Optional. Default value "Bearer".
	AuthScheme string

	// Validator is a function to validate key.
	Validator func(*woody.Ctx, string) (bool, error)

	// Context key to store the bearertoken from the token into context.
	// Optional. Default: "token".
	ContextKey string
}
```

## Default Config

```go
var ConfigDefault = Config{
	SuccessHandler: func(c *woody.Ctx) error {
		return c.Next()
	},
	ErrorHandler: func(c *woody.Ctx, err error) error {
		if err == ErrMissingOrMalformedAPIKey {
			return c.Status(woody.StatusUnauthorized).SendString(err.Error())
		}
		return c.Status(woody.StatusUnauthorized).SendString("Invalid or expired API Key")
	},
	KeyLookup:  "header:" + woody.HeaderAuthorization,
	AuthScheme: "Bearer",
	ContextKey: "token",
}
```
