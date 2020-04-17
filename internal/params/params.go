package params

type Params struct {
	CISBaseDir     string
	RouterConfig   string
	ProxyHost      string
	ProxyPort      string
	SessionTimeout int
}

func (p *Params) GetMessagesAboutUnsetParams() (messages []string) {
	if "" == p.CISBaseDir {
		messages = append(messages, "CIS base directory (-cis-dir) is not set")
	}
	if "" == p.RouterConfig {
		messages = append(messages, "Router config path (-router-conf) is not set")
	}
	if "" == p.ProxyHost {
		messages = append(messages, "Host (-host) is not set")
	}
	if "" == p.ProxyPort {
		messages = append(messages, "Port (-port) is not set")
	}
	if 0 == p.SessionTimeout {
		messages = append(messages, "Session timeout (-timeout) is not set")
	}

	return messages
}
