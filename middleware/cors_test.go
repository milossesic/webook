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

func setupCorsRouter() *httprouter.Router {
	router := httprouter.New()
	handler := middleware.Cors(
		func(w http.ResponseWriter, r *http.Request) *util.AppError {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))
			return nil
		})
	pingHandle := middleware.Context("/ping", logger.New(), handler)
	router.GET("/ping", pingHandle)
	return router
}

var _ = Describe("cors.go functions", func() {

	router := setupCorsRouter()
	ts := httptest.NewServer(router)
	defer ts.Close()

	It("Test cors", func() {
		req, err := http.NewRequest("GET", ts.URL+"/ping", nil)
		Expect(err).ToNot(HaveOccurred())
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)
		Expect(rr.Code).To(Equal(200))
		Expect(rr.Header().Get("Access-Control-Allow-Headers")).To(Equal(
			"Accept,Content-Type,If-None-Match"))
		Expect(rr.Header().Get("Access-Control-Expose-Headers")).To(Equal(
			"Content-Type,Cache-Control,Etag,Location"))
		Expect(rr.Header().Get("Access-Control-Max-Age")).To(Equal("3600"))
		Expect(rr.Header().Get("Access-Control-Allow-Origin")).To(Equal(""))
		Expect(rr.Header().Get("Access-Control-Allow-Credentials")).To(Equal("true"))
		Expect(rr.Header().Get("Access-Control-Allow-Methods")).To(Equal("POST,GET,PUT,DELETE"))
	})
})
