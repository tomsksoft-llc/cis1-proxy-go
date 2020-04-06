package requestfiltration

import (
	"net"
	"net/http"
	"regexp"
	"strings"
)

var (
	passPathSearch, _    = regexp.Compile(`http://[.\S]+:\d+(/[.\S]*)`)
	passHostSearch, _    = regexp.Compile(`http://([.\S]+):\d+/[.\S]*`)
	passAddressSearch, _ = regexp.Compile(`http://([.\S]+:\d+)/[.\S]*`)
	ipSearch, _          = regexp.Compile(`\d+.\d+.\d+.\d+`)
)

type Filter interface {
	FilterRequest(request *http.Request) bool
	GetLastPassAddress() string
	AddRoute(route Route)
}

func NewFilter() Filter {
	return &filter{routes: make(map[string]string)}
}

type filter struct {
	routes          map[string]string
	lastPassAddress string
}

func (this *filter) FilterRequest(request *http.Request) bool {
	for location, pass := range this.routes {

		if 0 == strings.Index(request.URL.Path, location) {
			this.modifyRequest(request, location)
			this.lastPassAddress = passAddressSearch.FindStringSubmatch(pass)[1]

			return true
		}
	}

	this.lastPassAddress = ""
	return false
}

func (this *filter) GetLastPassAddress() string {
	return this.lastPassAddress
}

func (this *filter) AddRoute(route Route) {
	// Route.Pass = resolvePassHost(Route.Pass)
	this.routes[route.Location] = route.Pass
}

func (this *filter) modifyRequest(request *http.Request, location string) {
	request.URL.Path = strings.Replace(request.URL.Path, location, passPathSearch.FindStringSubmatch(this.routes[location])[1], 1)
	request.Host = passHostSearch.FindStringSubmatch(this.routes[location])[1]
	request.Header.Add("Connection", "close")
}

func resolvePassHost(pass string) {
	var passHost = passHostSearch.FindStringSubmatch(pass)[1]

	if ipSearch.MatchString(passHost) {
		var resolveResult, _ = net.LookupAddr(passHost)
		pass = strings.Replace(pass, passHost, resolveResult[0], 1)
	}
}
