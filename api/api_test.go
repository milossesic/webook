package api_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	"git.corp.adobe.com/dc/notifications_load_test/api"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("api.go functions", func() {

	r := api.NewRouter(nil, nil)
	router, err := r.CreateAllRoutes()
	if err != nil {
		Fail(fmt.Sprintf("Failed to create all routes: %s", err.Error()))
	}

	Context("Get ping", func() {
		It("Returns 200", func() {
			req, err := http.NewRequest("GET", "/ping", nil)
			if err != nil {
				Fail("Error creating a new request")
			}
			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)
			Expect(rr.Code).To(Equal(200))
			Expect(rr.Body.String()).To(Equal(`pong`))
		})
	})

	Context("POST notification", func() {
		It("Returns 204", func() {
			req, err := http.NewRequest("POST", "/notification", nil)
			if err != nil {
				Fail("Error creating a new request")
			}
			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)
			Expect(rr.Code).To(Equal(204))
			Expect(rr.Body.String()).To(Equal(`{"Msg":"Hello World"}`))
		})
	})

	Context("Get internal server error", func() {
		It("Returns 200", func() {
			req, err := http.NewRequest("GET", "/notificationsloadtest/error", nil)
			if err != nil {
				Fail("Error creating a new request")
			}
			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)
			Expect(rr.Code).To(Equal(500))
		})
	})
})
