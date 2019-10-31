package config_test

import (
	"git.corp.adobe.com/dc/notifications_load_test/config"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("config.go functions", func() {

	BeforeEach(func() {
		config.UnsetConfig()
	})

	It("Get config method", func() {
		Expect(config.GetConfig()).ShouldNot(BeNil())
	})

	It("Read from environment variable", func() {
		Expect(config.ReadFromEnv()).ShouldNot(BeNil())
	})

	It("Read from config file", func() {
		Expect(config.ReadFromConfig("prod")).ShouldNot(BeNil())
	})

	It("Test config singleton creation", func() {

		var cc1, cc2 *config.AppConfig

		ch1 := make(chan *config.AppConfig)
		ch2 := make(chan *config.AppConfig)
		go func(c chan *config.AppConfig) {
			c <- config.GetConfig()
		}(ch1)
		go func(c chan *config.AppConfig) {
			c <- config.GetConfig()
		}(ch2)
		cc1 = <-ch1
		cc2 = <-ch2

		// make sure the two point to the same place in memory
		Expect(cc1).To(BeIdenticalTo(cc2))
	})
})
