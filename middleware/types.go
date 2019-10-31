package middleware

import (
	"net/http"

	"git.corp.adobe.com/dc/notifications_load_test/util"
)

// Handle function type
type Handle func(w http.ResponseWriter, r *http.Request) *util.AppError
