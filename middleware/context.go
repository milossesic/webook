package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"git.corp.adobe.com/dc/notifications_load_test/config"
	"git.corp.adobe.com/dc/notifications_load_test/util"
	"github.com/julienschmidt/httprouter"
	"github.com/pborman/uuid"
	"gopkg.in/inconshreveable/log15.v2"
)

const (
	ctxKeyLog = iota
	ctxKeyParams
	ctxReqID
)

var (
	bgContext = context.Background()
)

// CtxResponseWriter is a context-aware http.ResponseWriter.
// It prevents calls to the real http.ResponseWriter that may already be closed
// due to a context timeout
type CtxResponseWriter struct {
	w   http.ResponseWriter
	ctx context.Context
}

// Header method gets header
func (recv *CtxResponseWriter) Header() http.Header {
	return recv.w.Header()
}

func (recv *CtxResponseWriter) Write(bs []byte) (int, error) {
	err := recv.ctx.Err()
	if err != nil {
		return 0, err
	}

	return recv.w.Write(bs)
}

// WriteHeader method writes header
func (recv *CtxResponseWriter) WriteHeader(i int) {
	if recv.ctx.Err() == nil {
		recv.w.WriteHeader(i)
	}
}

// Context method
func Context(route string, inputLog log15.Logger, hdl Handle) httprouter.Handle {

	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		t0 := time.Now()
		cfg := config.GetConfig()

		// Allow services to define logging keys and values to track requests
		//
		// When defining a key and value the "X-Request-Id-Key" and
		// "X-Request-Id-Value" headers should be used
		//
		// If our default key "RequestId" is sufficient then services should
		// set "X-Request-Id" for the value
		//
		reqIDKey := r.Header.Get("X-Request-Id-Key")
		if reqIDKey == "" {
			reqIDKey = "RequestId"
		}

		reqIDVal := r.Header.Get("X-Request-Id")
		if reqIDVal == "" {
			reqIDVal = r.Header.Get("X-Request-Id-Value")
		}
		if reqIDVal == "" {
			reqIDVal = uuid.NewRandom().String()
		}

		log := inputLog.New(reqIDKey, reqIDVal, "Route", route)

		ctx, cancel := context.WithTimeout(bgContext, cfg.RequestTimeout)
		defer cancel()

		ctx = WithRequestLog(ctx, log)
		ctx = WithParams(ctx, ps)
		ctx = WithRequestID(ctx, reqIDVal)

		r = r.WithContext(ctx)

		// log incoming request
		log.Info("Request", "info", fmt.Sprintf("method: %s, path: %s", r.Method, r.URL.Path))

		// doneChan is buffered because it's possible for ctx.Done() and
		// `doneChan <- struct{}{}` to happen at the same time.
		//
		// If `case _ = <-ctx.Done():` is selected first and doneChan was
		// unbuffered then `doneChan <- struct{}{}` would block forever and
		// leak goroutines.
		//
		// From the language spec:
		// If one or more of the communications can proceed, a single one that
		// can proceed is chosen via a uniform pseudo-random selection.
		// Otherwise, if there is a default case, that case is chosen. If there
		// is no default case, the "select" statement blocks until at least one
		// of the communications can proceed.
		doneChan := make(chan struct{}, 1)

		go func() {

			defer func() {
				// handle any panic in code and log response
				if err := recover(); err != nil {
					log.Error("Response", "error", fmt.Sprintf(
						"msg: panic %v, code: %d, path: %s", err, 500, r.URL.Path))
					util.Error(w, reqIDVal, 500, "panic")
					doneChan <- struct{}{}
				}
			}()

			wProxy := &CtxResponseWriter{
				w:   w,
				ctx: ctx,
			}

			// catch any error and log response
			if err := hdl(wProxy, r); err != nil {
				log.Error("Response", "error", fmt.Sprintf(
					"msg: %s, code: %d, path: %s", err.Cause, err.HTTPCode, r.URL.Path))
			} else {
				log.Info("Response", "info", fmt.Sprintf("code: %d, path: %s", 200, r.URL.Path))
			}

			doneChan <- struct{}{}
		}()

		select {
		case _ = <-ctx.Done():
			// catch any timeout and log response
			log.Error("Response", "error", fmt.Sprintf(
				"msg: request timeout, code: %d, path: %s", 408, r.URL.Path))
			util.Error(w, reqIDVal, 408, "Request timeout")
		case _ = <-doneChan:
			// nothing to do
		}

		t1 := time.Now()

		log.Info("Profile", "Duration", t1.Sub(t0).Nanoseconds()/1000000)
	}
}

// WithRequestLog method sets logger in context
func WithRequestLog(ctx context.Context, value log15.Logger) context.Context {
	return context.WithValue(ctx, ctxKeyLog, value)
}

// GetRequestLog method returns logger from context
func GetRequestLog(ctx context.Context) log15.Logger {
	value, ok := ctx.Value(ctxKeyLog).(log15.Logger)
	if ok {
		return value
	}
	packageLogger.Error("context error", "type", "ctxKeyLog")
	// attempt to track this even though the request-scoped logger is missing
	reqID := uuid.NewRandom()
	return packageLogger.New("fallback", "fallback", "RequestId", reqID)
}

// WithParams method sets param in context
func WithParams(ctx context.Context, value httprouter.Params) context.Context {
	return context.WithValue(ctx, ctxKeyParams, value)
}

// GetParams method returns param from context
func GetParams(ctx context.Context) httprouter.Params {
	value, ok := ctx.Value(ctxKeyParams).(httprouter.Params)
	if ok {
		return value
	}
	packageLogger.Error("context error", "type", "ctxKeyParams")
	return make(httprouter.Params, 0)
}

// WithRequestID method sets x-request-id in context
func WithRequestID(ctx context.Context, value string) context.Context {
	return context.WithValue(ctx, ctxReqID, value)
}

// GetRequestID gets x-request-id from context
func GetRequestID(ctx context.Context) string {
	value, ok := ctx.Value(ctxReqID).(string)
	if ok {
		return value
	}
	return ""
}
