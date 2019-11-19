package platform

import (
	"fmt"
	"net/http/cookiejar"
)

type DepositAdapter interface {
	Send(p *Order) error
	StartAndGC(p *DepositReq) error
}

type Order struct {
	Id          int64
	Account     string
	DepositType int // 存入类型  4:人工存提;5:优惠活动;6:返水;7:补发派彩;99:其他
	Amount      float32
	PortalMemo  string // 前台备注，显示于会员端
	Memo        string // 后台备注
	// 以下是提交后返回的状态信息
	Status int8   // 提交状态 0:待处理；1：处理成功；2：处理失败；
	Msg    string // 提交返回信息
}

type DepositReq struct {
	ReqUrl   string // 请求域名
	Jar      *cookiejar.Jar
	Password string
}

type DepositInstance func() DepositAdapter

var depositAdapters = make(map[string]DepositInstance)

func RegisterDeposit(name string, adapter DepositInstance) {
	if adapter == nil {
		panic("platform: login adapter is nil")
	}
	if _, ok := depositAdapters[name]; ok {
		panic("platform: login called twice for adapter " + name)
	}
	depositAdapters[name] = adapter
}

func NewPlatformDeposit(adapterName string, p *DepositReq) (adapter DepositAdapter, err error) {
	instanceFunc, ok := depositAdapters[adapterName]
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
