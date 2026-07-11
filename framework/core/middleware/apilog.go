package middleware

import (
	"log"
	"net/http"
	"time"
)

type APILogger struct {
	Logger *log.Logger
}

func NewAPILogger(logger *log.Logger) *APILogger {
	return &APILogger{Logger: logger}
}

func (a *APILogger) Log(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		wrw := &apiWriter{ResponseWriter: w, code: 200}
		next.ServeHTTP(wrw, r)
		a.Logger.Printf("[API] %s %s %d %s", r.Method, r.URL.Path, wrw.code, time.Since(start))
	})
}

type apiWriter struct {
	http.ResponseWriter
	code int
}

func (rw *apiWriter) WriteHeader(c int) {
	rw.code = c
	rw.ResponseWriter.WriteHeader(c)
}
