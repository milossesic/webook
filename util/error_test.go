package util_test

import (
	"net/http/httptest"

	"git.corp.adobe.com/dc/notifications_load_test/util"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("error.go functions", func() {
	It("Test error response", func() {
		res := httptest.NewRecorder()
		util.Error(res, "sample-request-id", 401, "some error")
		Expect(res.Code).To(Equal(401))
		Expect(res.Body.String()).To(Equal(
			`{"msg":"Unauthorized, Request ID: sample-request-id","code":401,"cause":"some error"}`))
		Expect(res.HeaderMap["Content-Type"][0]).To(Equal("application/json"))
	})
})
