package logger_test

import (
	"bytes"
	"strings"

	"git.corp.adobe.com/dc/notifications_load_test/config"
	"git.corp.adobe.com/dc/notifications_load_test/logger"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	log15 "gopkg.in/inconshreveable/log15.v2"
)

var _ = Describe("logger.go functions", func() {

	var (
		log    log15.Logger
		cfg    config.AppConfig
		logFmt log15.Format
		buf    bytes.Buffer
	)

	BeforeEach(func() {
		log = logger.New()
		buf = bytes.Buffer{}
		logFmt = log15.LogfmtFormat()
	})

	Context("Test logger creation", func() {
		It("Test successful creation of logger", func() {
			Expect(log).NotTo(BeNil())
		})
	})

	Context("Test correct filter is being used", func() {
		It("Will print logs with lvl critical when loglevel is critical", func() {
			cfg = config.AppConfig{LogLevel: "critical"}
			logger.AddStackTraceLogging(log, &buf, logFmt, &cfg)
			log.Crit("Test")
			Expect(strings.Contains(buf.String(), "lvl=crit")).To(Equal(true))
			log.Error("Test")
			Expect(strings.Contains(buf.String(), "lvl=eror")).To(Equal(false))
		})

		It("Will print logs with lvl error and critical when loglevel is error", func() {
			cfg = config.AppConfig{LogLevel: "error"}
			logger.AddStackTraceLogging(log, &buf, logFmt, &cfg)
			log.Crit("Test")
			Expect(strings.Contains(buf.String(), "lvl=crit")).To(Equal(true))
			log.Error("Test")
			Expect(strings.Contains(buf.String(), "lvl=eror")).To(Equal(true))
			log.Info("Test")
			Expect(strings.Contains(buf.String(), "lvl=info")).To(Equal(false))
		})

		It("Will print logs with lvl error, critical and warn when loglevel is warning", func() {
			cfg = config.AppConfig{LogLevel: "warning"}
			logger.AddStackTraceLogging(log, &buf, logFmt, &cfg)
			log.Crit("Test")
			Expect(strings.Contains(buf.String(), "lvl=crit")).To(Equal(true))
			log.Error("Test")
			Expect(strings.Contains(buf.String(), "lvl=eror")).To(Equal(true))
			log.Warn("Test")
			Expect(strings.Contains(buf.String(), "lvl=warn")).To(Equal(true))
			log.Info("Test")
			Expect(strings.Contains(buf.String(), "lvl=info")).To(Equal(false))
		})

		It("Will print logs with lvl error, critical, warn and info when loglevel is info", func() {
			cfg = config.AppConfig{LogLevel: "info"}
			logger.AddStackTraceLogging(log, &buf, logFmt, &cfg)
			log.Crit("Test")
			Expect(strings.Contains(buf.String(), "lvl=crit")).To(Equal(true))
			log.Error("Test")
			Expect(strings.Contains(buf.String(), "lvl=eror")).To(Equal(true))
			log.Warn("Test")
			Expect(strings.Contains(buf.String(), "lvl=warn")).To(Equal(true))
			log.Info("Test")
			Expect(strings.Contains(buf.String(), "lvl=info")).To(Equal(true))
			log.Debug("Test")
			Expect(strings.Contains(buf.String(), "lvl=debug")).To(Equal(false))
		})

		It("Will print all logs when loglevel is debug", func() {
			cfg = config.AppConfig{LogLevel: "debug"}
			logger.AddStackTraceLogging(log, &buf, logFmt, &cfg)
			log.Crit("Test")
			Expect(strings.Contains(buf.String(), "lvl=crit")).To(Equal(true))
			log.Error("Test")
			Expect(strings.Contains(buf.String(), "lvl=eror")).To(Equal(true))
			log.Warn("Test")
			Expect(strings.Contains(buf.String(), "lvl=warn")).To(Equal(true))
			log.Info("Test")
			Expect(strings.Contains(buf.String(), "lvl=info")).To(Equal(true))
			log.Debug("Test")
			Expect(strings.Contains(buf.String(), "lvl=dbug")).To(Equal(true))
		})
	})
})
