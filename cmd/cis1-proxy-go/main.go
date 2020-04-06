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

	var isAllSet = true
	if "" == *host {
		fmt.Println("Host value (-a) is not set")
		isAllSet = false
	}
	if 0 == *port {
		fmt.Println("Port value (-p) is not set")
		isAllSet = false
	}
	if "" == *configPath {
		fmt.Println("Config path (-c) is not set")
		isAllSet = false
	}

	if true == isAllSet {
		var proxy = proxy.NewProxy()

		if nil != proxy.Configure(*configPath) {
			return
		}

		var address = net.JoinHostPort(*host, strconv.Itoa(*port))
		if nil != proxy.Listen(address) {
			return
		}

		proxy.Run(*sessionTimeout)
	}
}
