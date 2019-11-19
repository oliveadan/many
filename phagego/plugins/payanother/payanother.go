package payanother

import (
	"github.com/pkg/errors"
	"net/http"
	"fmt"
)

type SendParam struct {
	MertNo      string
	EncryptKey  string
	ConfValue   string // 通道的其他配置，必须用json格式
	OrderNo     string
	Amount      int
	BankCode    string
	BusType     int // 0:对私；1:对公
	AccountName string
	CardNo      string
	Attach      string
	NotifyURL   string
}

type SendResp struct {
	ThirdOrderNo string
	PayAmount    int
	Status       int // 1：成功；2：失败；3：处理中
	Msg          string
}

func Send(name string, param *SendParam) (*SendResp, error) {
	fmt.Println(param)
	if v, ok := GetBankMap(name)[param.BankCode]; ok {
		param.BankCode = v
	}
	switch name {
	case PayAnother_huitao:
		return new(HuiTaoPay).Pay(param)
	}
	return nil, errors.New("代付通道不存在")
}

func Notice(name string, w http.ResponseWriter, r *http.Request, getMd5Key func(args ... string) string, busCallback func(args ... string) bool) error {
	switch name {
	case PayAnother_huitao:
		return new(HuiTaoPay).Notice(w, r, getMd5Key, busCallback)
	}

	return errors.New("Channel not found")
}

type QueryBalanceParam struct {
	MertNo      string
	EncryptKey  string
	ConfValue   string // 通道的其他配置，必须用json格式
}

func QueryBalance(name string, param *QueryBalanceParam) (int, error) {
	switch name {
	case PayAnother_huitao:
		return new(HuiTaoPay).QueryBalance(param)
	}
	return 0, errors.New("该通道不支持余额查询")
}
