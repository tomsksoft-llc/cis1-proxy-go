package params

import "flag"

func (p *Params) ParseCommandLine() {
	flag.StringVar(&p.CISBaseDir, "cis-dir", p.CISBaseDir, "CIS base directory")
	flag.StringVar(&p.RouterConfig, "router-conf", p.RouterConfig, "Router config path")
	flag.StringVar(&p.ProxyHost, "host", p.ProxyHost, "Proxy host")
	flag.StringVar(&p.ProxyPort, "port", p.ProxyPort, "Proxy port")
	flag.IntVar(&p.SessionTimeout, "timeout", p.SessionTimeout, "Session timeout (seconds)")

	flag.Parse()
}
