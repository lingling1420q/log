package main

import (
	"context"
	"crypto/rand"
	"math/big"
	"net/http"

	"github.com/chanxuehong/log"
	"github.com/chanxuehong/log/trace"
)

// set defaults if necessary.
// normally sets when the program starts.
func init() {
	log.SetFormatter(log.TextFormatter)
	log.SetLevelString(log.DebugLevelString)
	log.SetDefaultOptions([]log.Option{
		log.WithFormatter(log.TextFormatter),
		log.WithLevelString(log.DebugLevelString),
	})
}

type mockResponseWriter struct {
	http.ResponseWriter
	header http.Header
}

func (w *mockResponseWriter) Header() http.Header {
	return w.header
}

func main() {
	// mock http.ResponseWriter and *http.Request
	w := &mockResponseWriter{
		header: make(http.Header),
	}
	r, _ := http.NewRequest(http.MethodGet, "http://www.example.com", nil)

	// handle http request
	middleware(handler)(w, r)
}

func middleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// gets traceId from request header, if not exists, generates it.
		traceId, ok := trace.FromRequest(r)
		if !ok {
			traceId = trace.NewTraceId()
			r.Header.Set(trace.TraceIdHeaderKey, traceId)
		}
		// sets traceId to response header
		w.Header().Set(trace.TraceIdHeaderKey, traceId)

		// adds traceId to request.Context
		r = r.WithContext(trace.NewContext(r.Context(), traceId))

		// adds log.Logger to request.Context
		l := log.New(log.WithTraceId(traceId))
		r = log.NewRequest(r, l)

		// call http handler
		next(w, r)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	l := log.MustFromRequest(r)
	l.Info("1.info message")
	l.Info("2.info message", "key1", 1, "key2", 2)

	l = l.WithField("key3", 3)
	l = l.WithFields("key4", 4, "key5", 5)
	l.Info("3.info message")

	ctx := log.NewContext(r.Context(), l)
	fn1(ctx)
}

func fn1(ctx context.Context) {
	l := log.MustFromContext(ctx)
	l.Info("4.info message")
	l.Info("5.info message", "key6", 6)

	l = l.WithField("key7", 7)
	ctx = log.NewContext(ctx, l)
	fn2(ctx)
}

func fn2(ctx context.Context) {
	log.InfoContext(ctx, "6.info message", "key8", 8) // shortcut
	fn3(ctx)
}

func fn3(ctx context.Context) {
	log.InfoContext(ctx, "7.info message", "key9", 9) // shortcut
}

func mian() {
	ctx := context.Background()

	v, err := rand.Int(rand.Reader, big.NewInt(int64(16)))
	if err != nil {
		log.ErrorContext(ctx, " rand.In failed", "error", err.Error())
		return
	}
	log.InfoContext(ctx, "the rand Int result", "v", v)
}
