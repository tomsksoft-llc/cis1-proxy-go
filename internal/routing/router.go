package routing

import (
	"encoding/json"
	"net/http"
	"strings"
)

type Router interface {
	FindRoute(req *http.Request) Route
	ParseConfig(b []byte) error
}

func NewRouter() Router {
	return new(router)
}

type router struct {
	routes []Route
}

func (r *router) FindRoute(req *http.Request) Route {
	for _, route := range r.routes {
		if 0 == strings.Index(req.URL.Path, route.GetTarget()) {
			return route
		}
	}

	return nil
}

func (r *router) ParseConfig(b []byte) error {
	var err error

	var tempRouteHTTPs []routeHTTP
	err = json.Unmarshal(b, &tempRouteHTTPs)
	if nil != err {
		return err
	}

	for _, route := range tempRouteHTTPs {
		if "" != route.Location {
			var newRoute = new(routeHTTP)
			*newRoute = route
			r.routes = append(r.routes, newRoute)
		}
	}

	var tempRouteJobs []routeJob
	err = json.Unmarshal(b, &tempRouteJobs)
	if nil != err {
		return err
	}

	for _, route := range tempRouteJobs {
		if "" != route.Job {
			var newRoute = new(routeJob)
			*newRoute = route
			r.routes = append(r.routes, newRoute)
		}
	}

	return nil
}
