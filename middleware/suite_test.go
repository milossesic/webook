package middleware_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestApi(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Middleware Test Suite")
}
