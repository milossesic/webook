package logger_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Logger Test Suite")
}
