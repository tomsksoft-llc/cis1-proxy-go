package proxy

import (
	"io/ioutil"
	"net"
	"sync"

	"github.com/tomsksoft-llc/cis1-proxy-go/internal/routing"
	"github.com/tomsksoft-llc/cis1-proxy-go/internal/session"
)

type Proxy interface {
	ConfigureRouter(configPath string) error
	Listen(address string) error
	Run(sessionTimeout int)
}

func NewProxy() Proxy {
	return &proxy{router: routing.NewRouter()}
}

type proxy struct {
	listener       net.Listener
	sessionTimeout int
	router         routing.Router
}

func (p *proxy) ConfigureRouter(configPath string) error {
	var b, err = ioutil.ReadFile(configPath)
	if nil != err {
		return err
	}

	return p.router.ParseConfig(b)
}

func (p *proxy) Listen(address string) error {
	var err error
	p.listener, err = net.Listen("tcp4", address)

	return err
}

func (p *proxy) Run(sessionTimeout int) { // todo: correct interruption (maybe?)
	p.sessionTimeout = sessionTimeout

	var wg sync.WaitGroup
	wg.Add(1)
	go p.accept()
	wg.Wait()
}

func (p *proxy) accept() {
	var conn, err = p.listener.Accept()
	if nil == err {
		p.onAccept(conn)
	}

	go p.accept()
}

func (p *proxy) onAccept(conn net.Conn) {
	var s = session.NewSession(conn, p.router)
	s.Run(p.sessionTimeout)
}
