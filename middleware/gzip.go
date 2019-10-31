package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"

	"git.corp.adobe.com/dc/notifications_load_test/util"
)

type gzipResponseWriter struct {
	wroteHeader bool
	gzipWriter  io.Writer
	w           http.ResponseWriter
}

// Header method gets header
func (recv *gzipResponseWriter) Header() http.Header {
	return recv.w.Header()
}

func (recv *gzipResponseWriter) Write(b []byte) (int, error) {
	return recv.gzipWriter.Write(b)
}

// WriteHeader method writes header
func (recv *gzipResponseWriter) WriteHeader(status int) {
	if !recv.wroteHeader {
		recv.wroteHeader = true
		recv.w.Header().Set("Content-Encoding", "gzip")
	}
	recv.w.WriteHeader(status)
}

// Gzip method
func Gzip(hdl Handle) Handle {

	return func(w http.ResponseWriter, r *http.Request) *util.AppError {

		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			return hdl(w, r)
		}

		gzipWriter := gzip.NewWriter(w)
		defer gzipWriter.Close()

		gzrw := &gzipResponseWriter{
			gzipWriter: gzipWriter,
			w:          w,
		}

		return hdl(gzrw, r)
	}
}
