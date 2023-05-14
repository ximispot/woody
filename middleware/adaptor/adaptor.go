package adaptor

import (
	"io"
	"net"
	"net/http"
	"reflect"
	"unsafe"

	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpadaptor"
	"github.com/ximispot/woody"
	"github.com/ximispot/woody/utils"
)

// HTTPHandlerFunc wraps net/http handler func to woody handler
func HTTPHandlerFunc(h http.HandlerFunc) woody.Handler {
	return HTTPHandler(h)
}

// HTTPHandler wraps net/http handler to woody handler
func HTTPHandler(h http.Handler) woody.Handler {
	return func(c *woody.Ctx) error {
		handler := fasthttpadaptor.NewFastHTTPHandler(h)
		handler(c.Context())
		return nil
	}
}

// CopyContextToWoodyContext copies the values of context.Context to a fasthttp.RequestCtx
func CopyContextToWoodyContext(context interface{}, requestContext *fasthttp.RequestCtx) {
	contextValues := reflect.ValueOf(context).Elem()
	contextKeys := reflect.TypeOf(context).Elem()
	if contextKeys.Kind() == reflect.Struct {
		var lastKey interface{}
		for i := 0; i < contextValues.NumField(); i++ {
			reflectValue := contextValues.Field(i)
			/* #nosec */
			reflectValue = reflect.NewAt(reflectValue.Type(), unsafe.Pointer(reflectValue.UnsafeAddr())).Elem()

			reflectField := contextKeys.Field(i)

			if reflectField.Name == "noCopy" {
				break
			} else if reflectField.Name == "Context" {
				CopyContextToWoodyContext(reflectValue.Interface(), requestContext)
			} else if reflectField.Name == "key" {
				lastKey = reflectValue.Interface()
			} else if lastKey != nil && reflectField.Name == "val" {
				requestContext.SetUserValue(lastKey, reflectValue.Interface())
			} else {
				lastKey = nil
			}
		}
	}
}

// HTTPMiddleware wraps net/http middleware to woody middleware
func HTTPMiddleware(mw func(http.Handler) http.Handler) woody.Handler {
	return func(c *woody.Ctx) error {
		var next bool
		nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next = true
			// Convert again in case request may modify by middleware
			c.Request().Header.SetMethod(r.Method)
			c.Request().SetRequestURI(r.RequestURI)
			c.Request().SetHost(r.Host)
			for key, val := range r.Header {
				for _, v := range val {
					c.Request().Header.Set(key, v)
				}
			}
			CopyContextToWoodyContext(r.Context(), c.Context())
		})

		if err := HTTPHandler(mw(nextHandler))(c); err != nil {
			return err
		}

		if next {
			return c.Next()
		}
		return nil
	}
}

// WoodyHandler wraps woody handler to net/http handler
func WoodyHandler(h woody.Handler) http.Handler {
	return WoodyHandlerFunc(h)
}

// WoodyHandlerFunc wraps woody handler to net/http handler func
func WoodyHandlerFunc(h woody.Handler) http.HandlerFunc {
	return handlerFunc(woody.New(), h)
}

// WoodyApp wraps woody app to net/http handler func
func WoodyApp(app *woody.App) http.HandlerFunc {
	return handlerFunc(app)
}

func handlerFunc(app *woody.App, h ...woody.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// New fasthttp request
		req := fasthttp.AcquireRequest()
		defer fasthttp.ReleaseRequest(req)
		// Convert net/http -> fasthttp request
		if r.Body != nil {
			body, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, utils.StatusMessage(woody.StatusInternalServerError), woody.StatusInternalServerError)
				return
			}

			req.Header.SetContentLength(len(body))
			_, err = req.BodyWriter().Write(body)
			if err != nil {
				http.Error(w, utils.StatusMessage(woody.StatusInternalServerError), woody.StatusInternalServerError)
				return
			}
		}
		req.Header.SetMethod(r.Method)
		req.SetRequestURI(r.RequestURI)
		req.SetHost(r.Host)
		for key, val := range r.Header {
			for _, v := range val {
				req.Header.Set(key, v)
			}
		}
		if _, _, err := net.SplitHostPort(r.RemoteAddr); err != nil && err.(*net.AddrError).Err == "missing port in address" { //nolint:errorlint, forcetypeassert // overlinting
			r.RemoteAddr = net.JoinHostPort(r.RemoteAddr, "80")
		}
		remoteAddr, err := net.ResolveTCPAddr("tcp", r.RemoteAddr)
		if err != nil {
			http.Error(w, utils.StatusMessage(woody.StatusInternalServerError), woody.StatusInternalServerError)
			return
		}

		// New fasthttp Ctx
		var fctx fasthttp.RequestCtx
		fctx.Init(req, remoteAddr, nil)
		if len(h) > 0 {
			// New woody Ctx
			ctx := app.AcquireCtx(&fctx)
			defer app.ReleaseCtx(ctx)
			// Execute woody Ctx
			err := h[0](ctx)
			if err != nil {
				_ = app.Config().ErrorHandler(ctx, err) //nolint:errcheck // not needed
			}
		} else {
			// Execute fasthttp Ctx though app.Handler
			app.Handler()(&fctx)
		}

		// Convert fasthttp Ctx > net/http
		fctx.Response.Header.VisitAll(func(k, v []byte) {
			w.Header().Add(string(k), string(v))
		})
		w.WriteHeader(fctx.Response.StatusCode())
		_, _ = w.Write(fctx.Response.Body()) //nolint:errcheck // not needed
	}
}
