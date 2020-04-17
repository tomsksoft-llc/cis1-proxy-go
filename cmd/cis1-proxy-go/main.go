package main

import (
	"fmt"
	"net"
	"os"

	"github.com/tomsksoft-llc/cis1-proxy-go/internal/params"
	"github.com/tomsksoft-llc/cis1-proxy-go/internal/proxy"
)

func main() {
	var (
		params params.Params
		err    error
	)

	err = params.ParseIni("proxy_config.ini")
	if nil != err {
		fmt.Println("Init:", err.Error())
		return
	}

	params.ParseCommandLine()

	var messages = params.GetMessagesAboutUnsetParams()
	if 0 != len(messages) {
		for _, message := range messages {
			fmt.Println(message)
		}
	}

	os.Setenv("cis_base_dir", params.CISBaseDir)

	var proxy = proxy.NewProxy()

	err = proxy.ConfigureRouter(params.RouterConfig)
	if nil != err {
		fmt.Printf("Router config: %s\r\n", err.Error())
		return
	}

	err = proxy.Listen(net.JoinHostPort(params.ProxyHost, params.ProxyPort))
	if nil != err {
		fmt.Printf("Listen: %s\r\n", err.Error())
		return
	}

	proxy.Run(params.SessionTimeout)
}
