// Swagger spec generator
//
// Everything you need to know about the service can be found here
//
//     Schemes: http
//     Host: localhost:8080
//     BasePath: /
//     Version: 1.0.0
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Security:
//     - Bearer: []
//
// swagger:meta
package main

import (
	"net/http"
	"os"
	"time"

	"git.corp.adobe.com/dc/notifications_load_test/api"
	"git.corp.adobe.com/dc/notifications_load_test/config"
	"git.corp.adobe.com/dc/notifications_load_test/logger"
)

func main() {

	var cfg = config.GetConfig()

	// default is infinite which can leak file descriptors and goroutines
	http.DefaultClient = &http.Client{
		Timeout: 20 * time.Minute,
	}

	log := logger.New()

	router, err := api.NewRouter(nil, nil).CreateAllRoutes()
	if err != nil {
		log.Error("Router error", "Error", err.Error())
		os.Exit(1)
	}

	log.Info("Listening", "port", cfg.AppPort)

	if err := http.ListenAndServe(":"+cfg.AppPort, router); err != nil {
		log.Error("Failed to start service", "Error", err.Error())
	}
}
