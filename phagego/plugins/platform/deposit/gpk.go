package deposit

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"phagego/plugins/platform"
	"net/http/cookiejar"
	"time"
	"encoding/json"
	"strings"
)

type GpkDeposit struct {
	reqUrl   string // 请求域名
	password string
	Jar      *cookiejar.Jar
}

func NewGpkDeposit() platform.DepositAdapter {
	return &GpkDeposit{}
}

func (a *GpkDeposit) Send(order *platform.Order) error {
	var failStatus int8 = 2
	if order.DepositType != 4 && order.DepositType != 5 && order.DepositType != 6 && order.DepositType != 7 {
		order.Status = failStatus
		order.Msg = "存款类型不存在"
		return errors.New("存款类型不存在")
	}
	_, token, err := a.HttpRequest(http.MethodPost, "/Member/DepositToken", nil)
	if err != nil {
		order.Status = failStatus
		order.Msg = "请求token失败"
		return err
	}
	fmt.Println(token)
	isReal := false
	if order.DepositType == 4 {
		isReal = true
	}
	params := map[string]interface{}{
		"AccountsString": order.Account,
		"Amount":         order.Amount,
		"AuditType":      "None",
		"DepositToken":   strings.Trim(token, "\""),
		"IsReal":         isReal,
		"PortalMemo":     order.PortalMemo,
		"Memo":           order.Memo,
		"Password":       a.password,
		"TimeStamp":      time.Now().UnixNano() / 1000000,
		"Type":           order.DepositType,
	}
	jsonParams, _ := json.Marshal(params)
	fmt.Println(string(jsonParams))
	resp, bodys, err := a.HttpRequest(http.MethodPost, "/Member/DepositSubmit", bytes.NewBuffer(jsonParams))
	if err != nil {
		order.Status = failStatus
		order.Msg = "提交存款失败"
		return err
	}
	fmt.Println(bodys)
	if resp.StatusCode != http.StatusOK {
		order.Status = failStatus
		order.Msg = "提交存款返回状态异常"
		return errors.New("请求返回状态异常")
	}
	if bodys != "" {
		if bodys != "true" {
			var m map[string]interface{}
			if err := json.Unmarshal([]byte(bodys), &m); err != nil {
				order.Status = failStatus
				order.Msg = "提交存款返回数据格式异常"
				return err
			}
			v := m["IsSuccess"]
			if vv, ok := v.(bool); ok && !vv {
				order.Status = failStatus
				v = m["ErrorMessage"]
				if msg, ok := v.(string); ok {
					order.Msg = msg
				} else {
					order.Msg = "提交存款返回信息未知"
				}
				return errors.New("存款失败")
			}
		}
	}
	order.Status = 1 // 成功
	return nil
}

func (a *GpkDeposit) HttpRequest(method string, url string, reqbody io.Reader) (resp *http.Response, bodys string, err error) {
	client := &http.Client{
		CheckRedirect: func(req1 *http.Request, via []*http.Request) error {
			return errors.New("url been Redirect")
		},
		Jar: a.Jar,
	}
	var req *http.Request
	req, err = http.NewRequest(method, a.reqUrl+url, reqbody)
	if err != nil {
		return
	}
	fmt.Println(a.reqUrl+url)
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/62.0.3202.94 Safari/537.36")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Content-Type", "application/json;charset=utf-8")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Close = true

	resp, err = client.Do(req)
	if err != nil {
		return
	} else {
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		bodys = string(body)
	}
	return resp, bodys, nil
}
func (a *GpkDeposit) StartAndGC(p *platform.DepositReq) error {
	a.reqUrl = p.ReqUrl
	a.Jar = p.Jar
	a.password = p.Password

	if !strings.HasPrefix(strings.ToLower(a.reqUrl), "http") {
		a.reqUrl = "http://" + a.reqUrl
	}
	return nil
}

func init() {
	fmt.Println("init gpk Deposit")
	platform.RegisterDeposit("gpk", NewGpkDeposit)
}
