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
	"phagego/common/utils"
)

type DafaLogin struct {
	reqUrl      string // 后台域名，只要 如 http://baidu.com
	Jar         *cookiejar.Jar
	username    string
	password    string
}

func NewDafaLogin() platform.LoginAdapter {
	return &DafaLogin{}
}

func (a *DafaLogin) Login() (*cookiejar.Jar, error) {

	var r http.Request
	r.ParseForm()
	r.Form.Add("managerName", a.username)
	r.Form.Add("password", utils.Md5(a.username + utils.Md5(a.password)))

	resp, dd, err := a.HttpRequest(http.MethodPost, "/v1/management/manager/login", strings.NewReader(r.Form.Encode()), func(req *http.Request) {
		req.Header.Set("Accept", "application/json, text/plain, */*")
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	})
	fmt.Println(dd)
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
	if fmt.Sprintf("%v%v", m["code"], m["msg"]) != "1登录成功" {
		return nil, errors.New("登录失败")
	}
	return a.Jar, nil
}

func (a *DafaLogin) CheckLogon() error {
	resp, _, err := a.HttpRequest(http.MethodGet, "/v1/management/manager/xAnnouncementList?pageNum=1&pageSize=20&", nil, func(req *http.Request) {
		req.Header.Set("Accept", "application/json, text/plain, */*")
	})
	if err == nil && resp.StatusCode == 200 {
		return nil
	}
	return errors.New("登录掉线，请重新登录")
}

func (a *DafaLogin) HttpRequest(method string, url string, reqbody io.Reader, setHeader func(reqt *http.Request)) (resp *http.Response, bodys string, err error) {
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

func (a *DafaLogin) StartAndGC(p *platform.LoginReq) error {
	a.reqUrl = strings.TrimRight(strings.TrimSpace(p.ReqUrl), "/")
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
	fmt.Println("init Dafa login")
	platform.RegisterLogin("dafa", NewDafaLogin)
}
