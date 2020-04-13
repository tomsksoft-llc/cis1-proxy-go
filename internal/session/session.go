package session

import (
	"bufio"
	"net"
	"net/http"
	"time"

	"github.com/tomsksoft-llc/cis1-proxy-go/internal/routing"
)

type Session interface {
	Run(timeout int)
}

func NewSession(clientConn net.Conn, router routing.Router) Session {
	return &session{
		client: clientReader{
			conn:        clientConn,
			inputBuffer: bufio.NewReader(clientConn),
		},
		router: router,
	}
}

type session struct {
	client  clientReader
	router  routing.Router
	timeout int
}

type clientReader struct {
	conn        net.Conn
	inputBuffer *bufio.Reader
}

func (s *session) Run(timeout int) {
	s.timeout = timeout
	go s.readRequest()
}

func (s *session) setTimeout() {
	s.client.conn.SetReadDeadline(
		time.Now().Add(
			time.Second * time.Duration(s.timeout),
		),
	)
}

func (s *session) resetTimeout() {
	s.client.conn.SetReadDeadline(time.Time{})
}

func (s *session) readRequest() {
	s.setTimeout()

	var req, err = http.ReadRequest(s.client.inputBuffer)
	if nil != err {
		s.client.conn.Close()
		return
	}

	s.onReadRequest(req)
}

func (s *session) onReadRequest(req *http.Request) {
	s.resetTimeout()

	var route = s.router.FindRoute(req)
	if nil == route {
		s.client.conn.Write([]byte("HTTP/1.1 404 Not Found\r\nContent-Length: 13\r\n\r\n404 Not Found"))
	} else {
		route.Process(req, s.client.conn)
	}

	go s.readRequest()
}
