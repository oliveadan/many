package login

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"phagego/plugins/platform"
	"net/http/cookiejar"
	"strings"
	"encoding/json"
	"bytes"
)

type GpkLogin struct {
	reqUrl      string // 后台域名，只要 如 http://baidu.com
	Jar         *cookiejar.Jar
	username    string
	password    string
}

func NewGpkLogin() platform.LoginAdapter {
	return &GpkLogin{}
}

func (a *GpkLogin) Login() (*cookiejar.Jar, error) {
	params := map[string]interface{}{
		"account": a.username,
		"password":         a.password,
	}
	jsonParams, _ := json.Marshal(params)
	resp, dd, err := a.HttpRequest(http.MethodPost, "/Account/ValidateAccount", bytes.NewBuffer(jsonParams), func(req *http.Request) {
		req.Header.Set("Accept", "application/json, text/plain, */*")
		req.Header.Set("Content-Type", "application/json;charset=utf-8")
		req.Header.Set("X-Requested-With", "XMLHttpRequest")
	})
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("返回状态异常")
	}
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(dd), &m); err != nil {
		return nil, err
	}
	if m["IsSuccess"] != true {
		return nil, errors.New("登录失败")
	}
	return a.Jar, nil
}

func (a *GpkLogin) CheckLogon() error {
	resp, _, err := a.HttpRequest(http.MethodPost, "/Home/GetOnlineMembersCount", nil, func(req *http.Request) {
		req.Header.Set("Accept", "application/json, text/plain, */*")
		req.Header.Set("Content-Type", "application/json;charset=utf-8")
		req.Header.Set("X-Requested-With", "XMLHttpRequest")
	})
	if err == nil && resp.StatusCode == 200 {
		return nil
	}
	return errors.New("登录掉线，请重新登录")
}

func (a *GpkLogin) HttpRequest(method string, url string, reqbody io.Reader, setHeader func(reqt *http.Request)) (resp *http.Response, bodys string, err error) {
	client := &http.Client{
		CheckRedirect: nil,
		Jar: a.Jar,
	}
	fmt.Println(a.reqUrl+url)
	fmt.Println(a.Jar)
	var req *http.Request
	req, err = http.NewRequest(method, a.reqUrl+url, reqbody)
	if err != nil {
		return
	}
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/62.0.3202.94 Safari/537.36")

	if setHeader != nil {
		setHeader(req)
	}
	client.CheckRedirect = func(req1 *http.Request, via []*http.Request) error {
		return errors.New("url been Redirect")
	}
	resp, err = client.Do(req)
	if err != nil {
		return
	} else {
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		bodys = string(body)
		if strings.Contains(bodys, "帐号或密码错误") {
			return resp, bodys, errors.New("帐号或密码错误")
		}
	}
	return resp, bodys, nil
}

func (a *GpkLogin) StartAndGC(p *platform.LoginReq) error {
	a.reqUrl = p.ReqUrl
	a.Jar = p.Jar
	a.username = p.Username
	a.password = p.Password
	if a.Jar == nil {
		a.Jar, _ = cookiejar.New(nil)
	}
	if !strings.HasPrefix(strings.ToLower(a.reqUrl), "http") {
		a.reqUrl = "http://" + a.reqUrl
	}
	return nil
}


func init() {
	fmt.Println("init gpk login")
	platform.RegisterLogin("gpk", NewGpkLogin)
}
