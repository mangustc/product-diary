package middleware

import (
	"bytes"
	"io"
	"net/http"

	"github.com/bmg-c/product-diary/logger"
)

type wrappedWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *wrappedWriter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
	w.statusCode = statusCode
}

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		wrapped := &wrappedWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		bodyTextBytes, _ := io.ReadAll(r.Body)
		r.Body.Close()
		r.Body = io.NopCloser(bytes.NewBuffer(bodyTextBytes))
		bodyText := string(bodyTextBytes)

		next.ServeHTTP(wrapped, r)

		logger.Info.Println(wrapped.statusCode, r.Method, r.URL.Path, bodyText)
	})
}
