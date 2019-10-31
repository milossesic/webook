package util_test

import (
	"net/http/httptest"

	"git.corp.adobe.com/dc/notifications_load_test/util"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("json.go functions", func() {
	It("Test json response", func() {
		type User struct {
			First string `json:"first"`
			Last  string `json:"last"`
		}
		res := httptest.NewRecorder()
		util.JSON(res, &User{First: "foo", Last: "bar"}, 200)
		Expect(res.Code).To(Equal(200))
		Expect(res.Body.String()).To(Equal(`{"first":"foo","last":"bar"}`))
		Expect(res.HeaderMap["Content-Type"][0]).To(Equal("application/json"))
	})

	It("Json marshal error", func() {
		res := httptest.NewRecorder()
		util.JSON(res, make(chan int), 200)
		Expect(res.Code).To(Equal(500))
		Expect(res.Body.String()).To(
			Equal(`{"msg":"Internal Server Error","code":500,"cause":"Json marshalling failed"}`))
		Expect(res.HeaderMap["Content-Type"][0]).To(Equal("application/json"))
	})
})
