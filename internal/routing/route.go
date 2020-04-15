package routing

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
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
		Args    []struct {
			Name  string
			Value string
		}
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
	var (
		jobExec = fmt.Sprintf("%s/core/startjob", os.Getenv("cis_base_dir"))
		jobArgs = []string{
			fmt.Sprintf("%s/%s", rj.Run.Project, rj.Run.Job),
		}
	)

	if 0 != len(rj.Run.Args) {
		jobArgs = append(jobArgs, "--params")
		for _, arg := range rj.Run.Args {
			jobArgs = append(jobArgs, arg.Name)
			jobArgs = append(jobArgs, arg.Value)
		}
	}

	var startjobCmd = exec.Command(jobExec, jobArgs...)

	go runJob(startjobCmd, clientConn)
}

func runJob(startjobCmd *exec.Cmd, clientConn net.Conn) {
	startjobCmd.Run()
	clientConn.Write([]byte("HTTP/1.1 200 OK\r\nContent-Length: 11\r\n\r\nJob started"))
}
