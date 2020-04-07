package session

import (
	"bufio"
	"io"
	"net"
	"net/http"
	"time"

	"github.com/tomsksoft-llc/cis1-proxy-go/internal/requestfiltration"
)

type Session interface {
	Run(timeout int)
}

func NewSession(conn net.Conn, f requestfiltration.Filter) Session {
	return &session{
		client: clientReader{
			conn:        conn,
			inputBuffer: bufio.NewReader(conn),
		},
		filter: f,
	}
}

type session struct {
	client struct {
		conn        net.Conn
		inputBuffer *bufio.Reader
	}
	filter  requestfiltration.Filter
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

func (s *session) closeConnection() {
	s.client.conn.Close()
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
		s.closeConnection()
		return
	}

	s.onReadRequest(req)
}

func (s *session) onReadRequest(req *http.Request) {
	s.resetTimeout()

	var passAddress, isPassed = s.filter.Filter(req)
	if false == isPassed {
		s.closeConnection()
		return
	}

	var serverConn, err = net.Dial("tcp4", passAddress)
	if err != nil {
		s.closeConnection()
		return
	}

	go s.readRequest()
	s.respond(serverConn, req)
}

func (s *session) respond(serverConn net.Conn, req *http.Request) {
	if nil != req.Write(serverConn) {
		s.closeConnection()
		serverConn.Close()
		return
	}

	io.Copy(s.client.conn, serverConn)
	serverConn.Close()
}
