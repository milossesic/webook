package api

import (
	//"bytes"
	"net/http"
	"fmt"
	"net/http/httputil"
	"time"

	//"time"

	"git.corp.adobe.com/dc/notifications_load_test/logger"
	"git.corp.adobe.com/dc/notifications_load_test/middleware"
	"git.corp.adobe.com/dc/notifications_load_test/util"
)

// Handler struct
type Handler struct {
}

// NewHandler method creates new instance of Handler
func NewHandler() *Handler {
	return &Handler{}
}

// pingHandler returns pong
//
// swagger:operation GET /ping ping
// ---
// summary: Gets pong
//
// responses:
//  "200":
//    description: "Returns pong"
//    schema:
//      type: string
//      example: pong
func (h *Handler) pingHandler(w http.ResponseWriter, r *http.Request) *util.AppError {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte("pong"))
	return nil
}

// myfirstapiHandler returns myfirstapi
//
// swagger:operation GET /notificationsloadtest/myfirstapi myfirstapi
// ---
// summary: Gets myfirstapi
//
// responses:
//  "200":
//    description: "Returns myfirstapi message"
//    schema:
//      "$ref": "#/definitions/myfirstapi"
func (h *Handler) notificationHandler(w http.ResponseWriter, r *http.Request) *util.AppError {
	//res := notification{Msg: ""}
	requestDump, err := httputil.DumpRequest(r, true)
	if err != nil {
		fmt.Println(err)
	}
	log := logger.New()
	logger := log.New("foo", "bar", "req")
	logger.Info("Request", "info", fmt.Sprintf("reqDump: %s", string(requestDump)))
	time.Sleep(2 * time.Second)
	fmt.Println(string(requestDump))
	util.JSON(w, nil, 204)
	return nil
}

func (h *Handler) notificationRegisterer(w http.ResponseWriter, r *http.Request) *util.AppError {
	hr := "X-AdobeSign-ClientId"
	cid := r.Header.Get(hr)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set(hr, cid)
	w.WriteHeader(200)
	w.Write([]byte(cid))
	return nil
}

func (h *Handler) notificationHandler2(w http.ResponseWriter, r *http.Request) *util.AppError {
	//res := notification{Msg: ""}
	requestDump, err := httputil.DumpRequest(r, true)
	if err != nil {
		fmt.Println(err)
	}
	log := logger.New()
	logger := log.New("foo", "bar", "req")
	logger.Info("Request", "info", fmt.Sprintf("reqDump: %s", string(requestDump)))
	//time.Sleep(3 * time.Second)
	fmt.Println(string(requestDump))
	util.JSON(w, nil, 204)
	return nil
}

func (h *Handler) notificationRegisterer2(w http.ResponseWriter, r *http.Request) *util.AppError {
	hr := "X-AdobeSign-ClientId"
	cid := r.Header.Get(hr)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set(hr, cid)
	w.WriteHeader(204)
	w.Write([]byte(cid))
	return nil
}

// errorHandler returns internal server error
//
// swagger:operation GET /notificationsloadtest/error error
// ---
// summary: Gets internal server error
//
// responses:
//  "500":
//    "$ref": "#/definitions/internalError"
func (h *Handler) errorHandler(w http.ResponseWriter, r *http.Request) *util.AppError {
	res := util.AppError{Err: "Internal Server Error", Cause: "Some Error", HTTPCode: 500}
	util.Error(w, middleware.GetRequestID(r.Context()), res.HTTPCode, res.Cause)
	return &res
}
