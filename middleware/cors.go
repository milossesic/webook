package middleware

import (
	"net/http"
	"time"

	"git.corp.adobe.com/dc/notifications_load_test/util"
	"github.com/martini-contrib/cors"
)

// Cors method
func Cors(hdl Handle) Handle {
	corsOpts := &cors.Options{
		AllowOrigins: []string{
			"*",
		},
		AllowCredentials: true,
		AllowMethods:     []string{"POST", "GET", "PUT", "DELETE"},
		AllowHeaders: []string{
			"Accept",
			"Content-Type",
			"If-None-Match",
		},
		MaxAge: 1 * time.Hour,
		ExposeHeaders: []string{
			"Content-Type",
			"Cache-Control",
			"Etag",
			"Location",
		},
	}

	corsHdl := cors.Allow(corsOpts)

	return func(w http.ResponseWriter, r *http.Request) *util.AppError {

		corsHdl(w, r)

		if r.Method != "OPTIONS" {
			return hdl(w, r)
		}

		return nil
	}
}
