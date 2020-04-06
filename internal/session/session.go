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

func NewSession(connection net.Conn, requestFilter requestfiltration.Filter) Session {
	return &session{client: clientReader{connection: connection, inputBuffer: bufio.NewReader(connection)}, requestFilter: requestFilter}
}

type session struct {
	client        clientReader
	requestFilter requestfiltration.Filter
	timeout       int
}

type clientReader struct {
	connection  net.Conn
	inputBuffer *bufio.Reader
}

func (this *session) Run(timeout int) {
	this.timeout = timeout
	go this.readRequest()
}

func (this *session) closeConnection() {
	this.client.connection.Close()
}

func (this *session) setReadTimeout() {
	this.client.connection.SetReadDeadline(time.Now().Add(time.Second * time.Duration(this.timeout)))
}

func (this *session) resetReadTimeout() {
	this.client.connection.SetReadDeadline(time.Time{})
}

func (this *session) readRequest() {
	this.setReadTimeout()

	var request, err = http.ReadRequest(this.client.inputBuffer)
	if err != nil {
		this.closeConnection()
		return
	}

	this.onReadRequest(request)
}

func (this *session) onReadRequest(request *http.Request) {
	this.resetReadTimeout()

	if false == this.requestFilter.FilterRequest(request) {
		this.closeConnection()
		return
	}

	var serverConnection, err = net.Dial("tcp4", this.requestFilter.GetLastPassAddress())
	if err != nil {
		this.closeConnection()
		return
	}

	go this.readRequest()
	go this.respond(serverConnection, request)
}

func (this *session) respond(serverConnection net.Conn, request *http.Request) {
	if nil != request.Write(serverConnection) {
		this.closeConnection()
		serverConnection.Close()
		return
	}

	io.Copy(this.client.connection, serverConnection)
	serverConnection.Close()
}
