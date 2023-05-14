---
id: envvar
title: EnvVar
---

EnvVar middleware for [Woody](https://github.com/gowoody/woody) that can be used to expose environment variables with various options.

## Signatures

```go
func New(config ...Config) woody.Handler
```

## Examples

Import the middleware package that is part of the Woody web framework

```go
import (
  "github.com/gowoody/woody/v2"
  "github.com/ximispot/woody/middleware/envvar"
)
```

After you initiate your Woody app, you can use the following possibilities:

```go
// Initialize default config
app.Use("/expose/envvars", envvar.New())

// Or extend your config for customization
app.Use("/expose/envvars", envvar.New(
	envvar.Config{
		ExportVars:  map[string]string{"testKey": "", "testDefaultKey": "testDefaultVal"},
		ExcludeVars: map[string]string{"excludeKey": ""},
	}),
)
```

:::note
You will need to provide a path to use the envvar middleware.
:::

## Response

Http response contract:
```
{
  "vars": {
    "someEnvVariable": "someValue",
    "anotherEnvVariable": "anotherValue",
  }
}

```

## Config

```go
// Config defines the config for middleware.
type Config struct {
    // ExportVars specifies the environment variables that should export
    ExportVars map[string]string
    // ExcludeVars specifies the environment variables that should not export
    ExcludeVars map[string]string
}

```

## Default Config

```go
Config{}
```
