package platform

import (
	"fmt"
	"net/http/cookiejar"
)

type LoginAdapter interface {
	Login() (*cookiejar.Jar, error)
	CheckLogon() error
	StartAndGC(p *LoginReq) error
}

type LoginReq struct {
	ReqUrl      string // 后台域名，只要 如 http://baidu.com
	Jar         *cookiejar.Jar
	Username    string
	Password    string
}

type LoginInstance func() LoginAdapter

var loginAdapters = make(map[string]LoginInstance)

func RegisterLogin(name string, adapter LoginInstance) {
	if adapter == nil {
		panic("platform: login adapter is nil")
	}
	if _, ok := loginAdapters[name]; ok {
		panic("platform: login called twice for adapter " + name)
	}
	loginAdapters[name] = adapter
}

func NewPlatformLogin(adapterName string, p *LoginReq) (adapter LoginAdapter, err error) {
	instanceFunc, ok := loginAdapters[adapterName]
	if !ok {
		err = fmt.Errorf("platform: unknown adapter name %q (forgot to import?)", adapterName)
		return
	}
	adapter = instanceFunc()
	err = adapter.StartAndGC(p)
	if err != nil {
		adapter = nil
	}
	return
}
