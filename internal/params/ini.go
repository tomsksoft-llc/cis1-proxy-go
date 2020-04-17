package params

import (
	"io/ioutil"

	"github.com/go-ini/ini"
)

func (p *Params) ParseIni(INIPath string) error {
	var (
		b   []byte
		err error
	)

	b, err = ioutil.ReadFile(INIPath)
	if nil != err {
		return err
	}

	var config *ini.File
	config, err = ini.Load(b)
	if nil != err {
		return err
	}

	err = p.parsePathsSection(config)
	if nil != err {
		return err
	}
	err = p.parseProxySection(config)
	if nil != err {
		return err
	}
	err = p.parseSessionSection(config)
	if nil != err {
		return err
	}

	return nil
}

func (p *Params) parsePathsSection(config *ini.File) error {
	var (
		section *ini.Section
		err     error
	)

	section, err = config.GetSection("paths")
	if nil != err {
		return err
	}

	var key *ini.Key

	key, err = section.GetKey("cis_base_dir")
	if nil != err {
		return err
	}
	p.CISBaseDir = key.String()

	key, err = section.GetKey("router_config")
	if nil != err {
		return err
	}
	p.RouterConfig = key.String()

	return nil
}

func (p *Params) parseProxySection(config *ini.File) error {
	var (
		section *ini.Section
		err     error
	)

	section, err = config.GetSection("proxy")
	if nil != err {
		return err
	}

	var key *ini.Key

	key, err = section.GetKey("host")
	if nil != err {
		return err
	}
	p.ProxyHost = key.String()

	key, err = section.GetKey("port")
	if nil != err {
		return err
	}
	p.ProxyPort = key.String()

	return nil
}

func (p *Params) parseSessionSection(config *ini.File) error {
	var (
		section *ini.Section
		err     error
	)

	section, err = config.GetSection("session")
	if nil != err {
		return err
	}

	var key *ini.Key

	key, err = section.GetKey("timeout")
	if nil != err {
		return err
	}
	p.SessionTimeout, err = key.Int()
	if nil != err {
		return err
	}

	return nil
}
