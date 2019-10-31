package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"time"

	"git.corp.adobe.com/dc/notifications_load_test/config"
	"git.corp.adobe.com/dc/notifications_load_test/logger"
	"git.corp.adobe.com/dc/notifications_load_test/middleware"
	"git.corp.adobe.com/dc/notifications_load_test/util"
	"github.com/julienschmidt/httprouter"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func setupRouter(response_type string) *httprouter.Router {
	router := httprouter.New()
	handler := func(w http.ResponseWriter, r *http.Request) *util.AppError {
		if response_type == "valid_response" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))
		} else if response_type == "request_timeout" {
			time.Sleep(config.GetConfig().RequestTimeout + 1*time.Second)
			w.WriteHeader(http.StatusRequestTimeout)
			w.Write([]byte("Request timeout"))
		} else if response_type == "panic" {
			panic("some problem")
		}
		return nil
	}
	pingHandle := middleware.Context("/ping", logger.New(), handler)
	router.GET("/ping", pingHandle)
	return router
}

var _ = Describe("context.go functions", func() {

	It("Valid response", func() {
		router := setupRouter("valid_response")
		ts := httptest.NewServer(router)
		defer ts.Close()
		req, err := http.NewRequest("GET", ts.URL+"/ping", nil)
		if err != nil {
			Fail("Error creating a new request")
		}
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)
		Expect(rr.Code).To(Equal(200))
	})

	It("Request timeout", func() {
		router := setupRouter("request_timeout")
		ts := httptest.NewServer(router)
		defer ts.Close()
		req, err := http.NewRequest("GET", ts.URL+"/ping", nil)
		if err != nil {
			Fail("Error creating a new request")
		}
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)
		Expect(rr.Code).To(Equal(408))
		Expect(rr.Header().Get("Content-Type")).To(Equal("application/json"))
	})

	It("Generate panic", func() {
		router := setupRouter("panic")
		ts := httptest.NewServer(router)
		defer ts.Close()
		req, err := http.NewRequest("GET", ts.URL+"/ping", nil)
		if err != nil {
			Fail("Error creating a new request")
		}
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)
		Expect(rr.Code).To(Equal(500))
		Expect(rr.Header().Get("Content-Type")).To(Equal("application/json"))
	})
})
