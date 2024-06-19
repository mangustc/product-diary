package middleware

import (
	"bytes"
	"io"
	"net/http"

	"github.com/bmg-c/product-diary/logger"
)

type wrappedWriter struct {
	http.ResponseWriter
	statusCode  int
	wroteHeader bool
}

func (w *wrappedWriter) WriteHeader(statusCode int) {
	if w.wroteHeader {
		return
	}
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
	w.wroteHeader = true
}

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		wrapped := &wrappedWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
			wroteHeader:    false,
		}

		bodyTextBytes, _ := io.ReadAll(r.Body)
		r.Body.Close()
		r.Body = io.NopCloser(bytes.NewBuffer(bodyTextBytes))
		bodyText := string(bodyTextBytes)

		next.ServeHTTP(wrapped, r)

		logger.Info.Printf("|%d| %s %s\t%s", wrapped.statusCode, r.Method, r.URL.Path, bodyText)
	})
}
