package web

import (
	"fmt"
	"net/http"
	"time"

	"kube-job-runner/pkg/app/reporter"

	"github.com/go-chi/chi/middleware"
)

func CreateHttResponseLogger(reporter *reporter.Reporter) func(next http.Handler) http.Handler {
	return middleware.RequestLogger(&structuredLogger{reporter})
}

type structuredLogger struct {
	reporter *reporter.Reporter
}

type requestInfo struct {
	scheme string
	method string
	uri    string
}

func (l *structuredLogger) NewLogEntry(r *http.Request) middleware.LogEntry {
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	method := r.Method
	uri := r.URL.Path

	return &structuredLoggerEntry{
		reporter:    l.reporter,
		requestInfo: requestInfo{scheme, method, uri},
	}
}

type structuredLoggerEntry struct {
	reporter    *reporter.Reporter
	requestInfo requestInfo
}

func (l *structuredLoggerEntry) Write(status, bytes int, elapsed time.Duration) {
	l.reporter.Info(
		"http.respond",
		fmt.Sprintf("Responded to %s %s with %d", l.requestInfo.method, l.requestInfo.uri, status),
		map[string]interface{}{
			"scheme":           l.requestInfo.scheme,
			"method":           l.requestInfo.method,
			"uri":              l.requestInfo.uri,
			"responseStatus":   status,
			"responseDuration": elapsed.Seconds(),
		},
	)
}

func (l *structuredLoggerEntry) Panic(v interface{}, stack []byte) {
	l.reporter.Info("unexpected.error", "Unexpected panic.", map[string]interface{}{
		"stack": string(stack),
		"panic": fmt.Sprintf("%+v", v),
	})
}
