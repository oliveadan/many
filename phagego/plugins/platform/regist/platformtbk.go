package regist

import (
	"net/http/cookiejar"
	"net/http"
	"strings"
	"fmt"
	"errors"
	"io/ioutil"
)
/*
  TBK平台
 */
func (a *PlatformRegister) tbk(param *RegParam) (int, string) {
	if len(param.WithdrawPass) < 6 || len(param.WithdrawPass) > 12 {
		return 101, "取款密码长度必须在6~12个字符"
	}

	var cookieJar, _ = cookiejar.New(nil)

	var req *http.Request
	var resp *http.Response

	var err error

	var isAgent bool
	// 判断是否代理地址
	if !strings.Contains(a.ReqUrl, "register") {
		isAgent = true
	}
	if strings.Count(a.ReqUrl, "::") == 2 {
		if param.Upline == "" {
			a.ReqUrl = strings.Replace(a.ReqUrl, "::", "", 2)
		} else {
			ss := strings.Split(a.ReqUrl, "::")
			ss[1] = param.Upline
			a.ReqUrl = strings.Join(ss, "")
		}
	}
	//fmt.Println(a.ReqUrl)
	req, err = http.NewRequest(http.MethodGet, a.ReqUrl, nil)
	if err != nil {
		return 500, "服务器内部错误"
	}
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/69.0.3497.100 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Close = true

	client := &http.Client{
		CheckRedirect: func(req1 *http.Request, via []*http.Request) error {
			return errors.New("请求被重定向")
		},
		Jar:           cookieJar,
	}
	resp, err = client.Do(req) // 请求注册网页，获取 cookie
	//fmt.Println(resp.StatusCode)
	if err != nil {
		if resp.StatusCode == http.StatusFound {
			u, _ := resp.Location()
			a.ReqUrl = u.String()
		} else {
			return 501, "请求失败"
		}
	} else if isAgent {
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		bodys := string(body)
		if strings.Contains(bodys, "无效的参数") {
			return 101, "推荐人ID不存在"
		}
	}

	var r http.Request
	r.ParseForm()
	r.Form.Add("username", param.Account)
	r.Form.Add("password", param.Password)
	r.Form.Add("repassword", param.Password)
	r.Form.Add("realname", param.RealName)
	r.Form.Add("tel", param.Mobile)
	r.Form.Add("weixin", param.WxNo)
	r.Form.Add("qkmm", param.WithdrawPass)
	r.Form.Add("question", param.Question)
	r.Form.Add("answer", param.Answer)
	r.Form.Add("qq", param.QqNo)

	req, err = http.NewRequest(a.ReqMethod, a.ReqUrl, strings.NewReader(r.Form.Encode()))
	if err != nil {
		return 500, "服务器内部错误"
	}
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/69.0.3497.100 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Close = true

	client1 := &http.Client{
		Transport: nil,
		CheckRedirect: nil,
		Jar:     cookieJar,
		Timeout: 0,
	}

	resp, err = client1.Do(req)
	if err != nil {
		fmt.Println("异常：", err.Error())
		if resp != nil {
			return resp.StatusCode, "请求异常"
		} else {
			return 501, "请求失败"
		}
	} else {
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		bodys := string(body)
		//d1 := []byte(bodys)
		//ioutil.WriteFile("test.txt", d1, 0644)
		if strings.Contains(bodys, param.Account) && strings.Contains(bodys, "/cn/logout") {
			return 1, ""
		} else if strings.Contains(bodys, "该名称已经被使用,请更换") {
			return 101, "该名称已经被使用,请更换"
		} else if strings.Contains(bodys, "指定时间内超过注册数") {
			return 101, "指定时间内超过注册数"
		}
		return 200, "验证不通过，请检查"
	}
}
