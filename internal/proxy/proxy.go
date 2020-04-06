package proxy

import (
	"io/ioutil"
	"net"
	"sync"

	"github.com/tomsksoft-llc/cis1-proxy-go/internal/requestfiltration"
	"github.com/tomsksoft-llc/cis1-proxy-go/internal/session"
)

type Proxy interface {
	Configure(configPath string) error
	Listen(address string) error
	Run(sessionTimeout int)
}

func NewProxy() Proxy {
	return &proxy{requestFilter: requestfiltration.NewFilter()}
}

type proxy struct {
	listener       net.Listener
	sessionTimeout int
	requestFilter  requestfiltration.Filter
}

func (this *proxy) Configure(configPath string) error {
	var configData, err = ioutil.ReadFile(configPath)
	if err != nil {
		return err
	}

	var routes []requestfiltration.Route
	routes, err = requestfiltration.ParseRoutes(configData)
	if err != nil {
		return err
	}

	for _, route := range routes {
		this.requestFilter.AddRoute(route)
	}

	return nil
}

func (this *proxy) Listen(address string) error {
	var err error
	this.listener, err = net.Listen("tcp4", address)

	return err
}

func (this *proxy) Run(sessionTimeout int) { // todo: correct interruption (maybe?)
	this.sessionTimeout = sessionTimeout

	var waitGroup sync.WaitGroup
	waitGroup.Add(1)
	go this.accept()
	waitGroup.Wait()
}

func (this *proxy) accept() {
	var connection, err = this.listener.Accept()
	if err == nil {
		this.onAccept(connection)
	}

	go this.accept()
}

func (this *proxy) onAccept(connection net.Conn) {
	var newSession = session.NewSession(connection, this.requestFilter)
	newSession.Run(this.sessionTimeout)
}
