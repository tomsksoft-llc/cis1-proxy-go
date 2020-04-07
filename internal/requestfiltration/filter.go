package requestfiltration

import (
	"net"
	"net/http"
	"strings"
)

type Filter interface {
	Filter(req *http.Request) (passAddress string, isPassed bool)
	AddRoutes(routes []Route)
}

func NewFilter() Filter {
	return &filter{routes: make(map[string]RoutePass)}
}

type filter struct {
	routes map[string]RoutePass
}

func (f *filter) Filter(req *http.Request) (passAddress string, isPassed bool) {
	for location, _ := range f.routes {

		if 0 == strings.Index(req.URL.Path, location) {
			f.modifyRequest(req, location)

			passAddress = net.JoinHostPort(f.routes[location].Host, f.routes[location].Port)
			return passAddress, true
		}
	}

	return "", false
}

func (f *filter) AddRoutes(routes []Route) {
	for _, route := range routes {
		f.routes[route.Location] = route.Pass
	}
}

func (f *filter) modifyRequest(req *http.Request, location string) {
	req.URL.Path = strings.Replace(req.URL.Path, location, f.routes[location].Path, 1)
	req.Host = f.routes[location].Host
	req.Header.Add("Connection", "close")
}
