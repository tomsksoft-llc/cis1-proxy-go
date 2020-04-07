package main

import (
	"flag"
	"fmt"
	"net"
	"strconv"

	"github.com/tomsksoft-llc/cis1-proxy-go/internal/proxy"
)

func main() {
	var (
		host           = flag.String("a", "", "Proxy host")
		port           = flag.Int("p", 0, "Proxy port")
		configPath     = flag.String("c", "", "Config file path")
		sessionTimeout = flag.Int("t", 60, "Session timeout (sec)")
	)

	flag.Parse()

	var areAllParamsSet = true
	if "" == *host {
		fmt.Println("Host value (-a) is not set")
		areAllParamsSet = false
	}
	if 0 == *port {
		fmt.Println("Port value (-p) is not set")
		areAllParamsSet = false
	}
	if "" == *configPath {
		fmt.Println("Config path (-c) is not set")
		areAllParamsSet = false
	}

	if true == areAllParamsSet {
		var (
			err   error
			proxy = proxy.NewProxy()
		)

		err = proxy.Configure(*configPath)
		if nil != err {
			fmt.Printf("Config: %s\r\n", err.Error())
			return
		}

		var address = net.JoinHostPort(*host, strconv.Itoa(*port))
		err = proxy.Listen(address)
		if nil != err {
			fmt.Printf("Listen: %s\r\n", err.Error())
			return
		}

		proxy.Run(*sessionTimeout)
	}
}
