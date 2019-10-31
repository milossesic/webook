package api

import (
	"git.corp.adobe.com/dc/notifications_load_test/logger"
	"git.corp.adobe.com/dc/notifications_load_test/middleware"
	"github.com/julienschmidt/httprouter"
)

// Router struct
type Router struct {
	router  *httprouter.Router
	handler *Handler
}

// NewRouter method creates new Router instance
func NewRouter(router *httprouter.Router, handler *Handler) *Router {

	if router == nil {
		router = httprouter.New()
	}

	if handler == nil {
		handler = NewHandler()
	}

	return &Router{
		router:  router,
		handler: handler,
	}
}

// CreateAllRoutes method defines middleware and all the possible routes
func (r *Router) CreateAllRoutes() (*httprouter.Router, error) {

	middlewares := []func(middleware.Handle) middleware.Handle{
		// enable cors
		middleware.Cors,
		// enable gzip response compression
		middleware.Gzip,
	}

	// ping
	createRoute(r.router.GET, "/ping", r.handler.pingHandler, middlewares...)

	// notificationsloadtest/myfirstapi
	createRoute(r.router.POST, "/notification", r.handler.notificationHandler, middlewares...)
	createRoute(r.router.GET, "/notification", r.handler.notificationRegisterer, middlewares...)
	createRoute(r.router.POST, "/notification2", r.handler.notificationHandler2, middlewares...)
	createRoute(r.router.GET, "/notification2", r.handler.notificationRegisterer2, middlewares...)

	// notificationsloadtest/error
	createRoute(r.router.GET, "/notificationsloadtest/error", r.handler.errorHandler, middlewares...)

	return r.router, nil
}

// createRoute method sends request to middleware followed by actual method
func createRoute(method func(path string, handle httprouter.Handle),
	path string,
	handler middleware.Handle,
	middlewares ...func(middleware.Handle) middleware.Handle) {

	for _, m := range middlewares {
		handler = m(handler)
	}

	log := logger.New()

	routeHandle := middleware.Context(path, log, handler)

	method(path, routeHandle)
}
