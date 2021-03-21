package endpoint

import (
	"bytes"

	"github.com/rs/zerolog/log"
	"github.com/valyala/fasthttp"
)

const partsLen = 3

func (e *Endpoint) RequestHandler() fasthttp.RequestHandler {
	var (
		pathSeparator = []byte("/")
		resizePath    = []byte("resize/")
		getter        = e.getFileGetter()
	)

	return func(ctx *fasthttp.RequestCtx) {
		logger := log.With().Fields(errorFields(ctx)).Logger()

		defer func() {
			if r := recover(); r != nil {
				logger.Printf("recovered: %s", r)
			}
		}()

		defer func() {
			logger.Print()
		}()

		path := ctx.RequestURI()

		// match /
		if bytes.Equal(path, pathSeparator) {
			ctx.Response.SetStatusCode(fasthttp.StatusOK)

			return
		}

		// match /:bucket/:path
		parts := bytes.Split(path, pathSeparator)
		count := len(parts)

		if count < partsLen {
			logger.Print("bad request")
			ctx.Response.SetStatusCode(fasthttp.StatusBadRequest)

			return
		}

		switch {
		case bytes.HasPrefix(path, resizePath):
			ctx.Response.SetStatusCode(fasthttp.StatusNotFound)

			return
		default:
			err := getter(ctx, parts[1], bytes.Join(parts[2:], pathSeparator))
			if err != nil {
				logger.Err(err)
			}

			return
		}
	}
}

func errorFields(ctx *fasthttp.RequestCtx) map[string]interface{} {
	return map[string]interface{}{
		"uri":    string(ctx.RequestURI()),
		"method": string(ctx.Method()),
		"ip":     ctx.RemoteIP(),
		"status": ctx.Response.StatusCode(),
	}
}
