package platform

import (
	"fmt"
	"net/http/cookiejar"
	"time"
)

type ExportAdapter interface {
	Export(p ExportParam) (error, []ExportData, string)
	StartAndGC(p *ExportReq) error
}

type ExportParam struct {
	DataType  string
	Account   string
	StartTime time.Time
	EndTime   time.Time
}

type ExportData struct {
	Id       string
	Account  string
	Amount   float32
	DataTime time.Time
	Remark   string
}

type ExportReq struct {
	ReqUrl   string // 请求域名
	Jar      *cookiejar.Jar
	Password string
}

type ExportInstance func() ExportAdapter

var ExportAdapters = make(map[string]ExportInstance)

func RegisterExport(name string, adapter ExportInstance) {
	if adapter == nil {
		panic("platform: login adapter is nil")
	}
	if _, ok := ExportAdapters[name]; ok {
		panic("platform: login called twice for adapter " + name)
	}
	ExportAdapters[name] = adapter
}

func NewPlatformExport(adapterName string, p *ExportReq) (adapter ExportAdapter, err error) {
	instanceFunc, ok := ExportAdapters[adapterName]
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
