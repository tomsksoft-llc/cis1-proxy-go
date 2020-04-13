package routing

import (
	"io"
	"net"
	"net/http"
	"strings"
)

type Route interface {
	GetTarget() string
	Process(req *http.Request, clientConn net.Conn)
}

type routeHTTP struct {
	Location string
	Pass     struct {
		Host string
		Port string
		Path string
	}
}

type routeJob struct {
	Job string
	Run struct {
		Project string
		Job     string
		Args    []string
	}
}

func (rh *routeHTTP) GetTarget() string {
	return rh.Location
}

func (rj *routeJob) GetTarget() string {
	return rj.Job
}

func (rh *routeHTTP) Process(req *http.Request, clientConn net.Conn) {
	var serverConn, err = net.Dial("tcp4", net.JoinHostPort(rh.Pass.Host, rh.Pass.Port))
	if nil != err {
		return
	}

	req.URL.Path = strings.Replace(req.URL.Path, rh.Location, rh.Pass.Path, 1)
	req.Host = rh.Pass.Host
	req.Header.Add("Connection", "close")

	go passRequest(req, clientConn, serverConn)
}

func passRequest(req *http.Request, clientConn net.Conn, serverConn net.Conn) {
	var err = req.Write(serverConn)
	if nil != err {
		serverConn.Close()
		return
	}

	go passResponse(clientConn, serverConn)
}

func passResponse(clientConn net.Conn, serverConn net.Conn) {
	io.Copy(clientConn, serverConn)
	serverConn.Close()
}

func (rj *routeJob) Process(req *http.Request, clientConn net.Conn) {
	clientConn.Write([]byte("HTTP/1.1 404 Not Found\r\nContent-Length: 7\r\n\r\nJob run"))
}
