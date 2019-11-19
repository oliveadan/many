package regist

import (
	"net/http"
	"net/http/cookiejar"
	"crypto/tls"
	"time"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"strings"
	"errors"
	"strconv"
)
/*
  BOS平台
 */
func (a *PlatformRegister) boss(param *RegParam) (int, string)  {
	// 验证
	if _, err := strconv.ParseUint(param.WithdrawPass, 10, 64); err != nil {
		return 101, "取款密码必须为数字"
	} else if len(param.WithdrawPass) != 4 {
		return 101, "取款密码必须为4位数字"
	}

	var cookieJar, _ = cookiejar.New(nil)

	var req *http.Request
	var resp *http.Response

	var err error

	var isAgent bool
	// 判断是否代理地址
	if !strings.Contains(a.ReqUrl, "signup") {
		isAgent = true
	}
	if strings.Contains(a.ReqUrl, "?f=") && strings.Count(a.ReqUrl, "::") == 2 {
		if param.Upline == "" {
			a.ReqUrl = strings.Replace(a.ReqUrl, "::", "", 2)
		} else {
			ss := strings.Split(a.ReqUrl, "::")
			ss[1] = param.Upline
			a.ReqUrl = strings.Join(ss, "")
		}
	}
	req, err = http.NewRequest(http.MethodGet, a.ReqUrl, nil)
	if err != nil {
		return 500, "服务器内部错误"
	}
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/69.0.3497.100 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Close = true

	config := &tls.Config{InsecureSkipVerify: true}
	tr := &http.Transport{TLSClientConfig: config}
	client := &http.Client{
		Transport: tr,
		Timeout:   15 * time.Second,
		CheckRedirect: nil,
		Jar:           cookieJar,
	}
	resp, err = client.Do(req) // 请求注册网页，获取 cookie

	if err != nil {
		return 501, "请求失败"
	} else if isAgent {
		a.ReqUrl = "https://" + req.Host + "/signup"
		req, err = http.NewRequest(http.MethodGet, a.ReqUrl, nil)
		if err != nil {
			return 500, "服务器内部错误"
		}
		req.Header.Set("Cache-Control", "no-cache")
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/69.0.3497.100 Safari/537.36")
		req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8")
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.Close = true

		config := &tls.Config{InsecureSkipVerify: true}
		tr := &http.Transport{TLSClientConfig: config}
		client := &http.Client{
			Transport: tr,
			Timeout:   15 * time.Second,
			CheckRedirect: nil,
			Jar:           cookieJar,
		}
		resp, err = client.Do(req) // 请求注册网页，获取 cookie
	}

	dom, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		return 501, "请求失败"
	}
	slt := dom.Find("input[name=_token]").First()
	token, _ := slt.Attr("value")
	//fmt.Println("token=", token)

	var r http.Request
	r.ParseForm()
	r.Form.Add("_token", token)
	r.Form.Add("upline", param.Upline)
	r.Form.Add("username", param.Account)
	r.Form.Add("password", param.Password)
	r.Form.Add("password_confirm", param.Password)
	r.Form.Add("withdrawal_code", param.WithdrawPass)
	r.Form.Add("fullname", param.RealName)
	r.Form.Add("qq", param.QqNo)
	r.Form.Add("email", param.Email)
	r.Form.Add("tel", param.Mobile)
	r.Form.Add("wechat", param.WxNo)
	//r.Form.Add("agree", "Y")
	//r.Form.Add("OK2", "立即注册")

	fmt.Println(a.ReqUrl)
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
		Transport: tr,
		Timeout:   15 * time.Second,
		CheckRedirect: func(req1 *http.Request, via []*http.Request) error {
			return errors.New("请求被重定向")
		},
		Jar:           cookieJar,
	}

	resp, err = client1.Do(req)
	if err != nil {
		fmt.Println("异常：", err.Error())
		if resp != nil && resp.StatusCode == http.StatusFound { //status code 302
			fmt.Println("重定向地址：", resp.Header.Get("Location"))
			return 302, "请求重定向"
		} else if resp != nil {
			return resp.StatusCode, "请求异常"
		} else {
			return 501, "请求失败"
		}
	} else {
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		bodys := string(body)
		d1 := []byte(bodys)
		ioutil.WriteFile("test.txt", d1, 0644)
		if strings.Contains(bodys, "注册成功") {
			return 1, ""
		}
		if strings.Contains(bodys,"alert(") {
			msg := strings.Split(bodys, "\"")
			if len(msg) > 2 {
				return 103, msg[1]
			}
 			return 103, "请检查信息"
		}

		return 102, "验证不通过或用户名不可用，请检查注册信息或更换用户名"
	}
}
