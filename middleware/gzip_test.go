package middleware_test

import (
	"net/http"
	"net/http/httptest"

	"git.corp.adobe.com/dc/notifications_load_test/logger"
	"git.corp.adobe.com/dc/notifications_load_test/middleware"
	"git.corp.adobe.com/dc/notifications_load_test/util"
	"github.com/julienschmidt/httprouter"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func setupGzipRouter() *httprouter.Router {
	router := httprouter.New()
	handler := middleware.Gzip(
		func(w http.ResponseWriter, r *http.Request) *util.AppError {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))
			return nil
		})
	pingHandle := middleware.Context("/ping", logger.New(), handler)
	router.GET("/ping", pingHandle)
	return router
}

var _ = Describe("gzip.go functions", func() {

	router := setupGzipRouter()
	ts := httptest.NewServer(router)
	defer ts.Close()

	It("Test gzip", func() {
		req, err := http.NewRequest("GET", ts.URL+"/ping", nil)
		req.Header.Set("Accept-Encoding", "gzip")
		Expect(err).ToNot(HaveOccurred())
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)
		Expect(rr.Code).To(Equal(200))
		Expect(rr.Header().Get("Content-Encoding")).To(Equal("gzip"))
	})
})
