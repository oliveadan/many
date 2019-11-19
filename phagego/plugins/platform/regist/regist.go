package regist

import (
	"fmt"
	"net/http/cookiejar"
)

const (
	PlatformTypeTBK  = "TBK"
	PlatformTypeBoss = "BOS"
	PlatformTypeDafa = "dafa"
)

type PlatformRegister struct {
	PlatformType string // 平台类型
	PlatformName string // 平台名称 （个别平台有用，比如大发）
	// 请求头信息
	ReqMethod string // 请求方法 POST,GET
	ReqUrl    string // 请求完整地址
	Jar       *cookiejar.Jar
}

type RegParam struct {
	Upline       string // 推广代码
	Account      string // 账号
	Password     string // 密码
	RePassword   string // 确认密码
	RealName     string // 真实姓名
	Mobile       string // 手机号
	WithdrawPass string // 取款密码
	WxNo         string // 微信号
	QqNo         string // QQ号
	Email        string // 邮箱
	Birthday     string // 生日
	Question     string // 提示问题
	Answer       string // 提示答案
}

func (a *PlatformRegister) Regist(param *RegParam) (bool, string) {
	if a.PlatformType == "" {
		return false, "平台类型未配置"
	}
	if a.ReqMethod == "" {
		return false, "请求方式未配置"
	}
	if a.ReqUrl == "" {
		return false, "注册网址未配置"
	}
	if a.Jar == nil {
		a.Jar, _ = cookiejar.New(nil)
	}

	var code int
	var msg string
	switch a.PlatformType {
	case PlatformTypeTBK:
		code, msg = a.tbk(param)
	case PlatformTypeBoss:
		code, msg = a.boss(param)
	case PlatformTypeDafa:
		code, msg = a.dafa(param)
	default:
		code = 2
		msg = "平台类型配置错误"
	}
	fmt.Println("返回状态：", code, " 信息：", msg)
	return code == 1, msg
}

func (a *PlatformRegister) GetCaptcha() {

}
