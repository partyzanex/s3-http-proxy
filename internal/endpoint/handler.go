package endpoint

import (
	"bytes"
	log "github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

var (
	pathSeparator = []byte("/")
	resizePath    = []byte("resize/")
)

const partsLen = 3

func (e *Endpoint) RequestHandler() fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		defer func() {
			if r := recover(); r != nil {
				log.WithFields(errorFields(ctx)).
					Errorf("recovered: %s", r)
			}
		}()

		defer func() {
			log.WithFields(log.Fields{
				"uri":    string(ctx.RequestURI()),
				"method": string(ctx.Method()),
				"status": ctx.Response.StatusCode(),
				"ip":     ctx.RemoteIP(),
			}).Info()
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
			log.Error("bad request")
			ctx.Response.SetStatusCode(fasthttp.StatusBadRequest)
			return
		}

		switch {
		case bytes.HasPrefix(path, resizePath):
			ctx.Response.SetStatusCode(fasthttp.StatusNotFound)
			return
		default:
			err := e.getFile(ctx, parts[1], bytes.Join(parts[2:], pathSeparator))
			if err != nil {
				log.Error(err)
			}

			return
		}
	}
}

func errorFields(ctx *fasthttp.RequestCtx) log.Fields {
	return log.Fields{
		"uri":    string(ctx.RequestURI()),
		"method": string(ctx.Method()),
		"ip":     ctx.RemoteIP(),
	}
}
