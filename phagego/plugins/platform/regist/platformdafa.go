package regist

import (
	"net/http"
	"strings"
	"fmt"
	"io/ioutil"
	"encoding/json"
	"phagego/common/utils"
)
/*
  大发平台
 */

func (a *PlatformRegister) dafa(param *RegParam) (int, string) {
	var cookieJar = a.Jar

	var req *http.Request
	var resp *http.Response

	var err error

	if strings.Count(a.ReqUrl, "::") == 2 {
		ss := strings.Split(a.ReqUrl, "::")
		if param.Upline == "" {
			a.ReqUrl = strings.Replace(a.ReqUrl, "::", "", 2)
			param.Upline = ss[1]
		} else {
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
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Close = true

	client := &http.Client{
		CheckRedirect: nil,
		Jar:           cookieJar,
	}
	resp, err = client.Do(req) // 请求注册网页，获取 cookie
	//fmt.Println(resp.StatusCode)
	if err != nil {
		return 501, "请求失败"
	}

	// 获取验证码
	var r2 http.Request
	r2.ParseForm()
	r2.Form.Add("Action", "GetImageCode")
	r2.Form.Add("SourceName", "PC")

	a.ReqUrl = fmt.Sprintf("http://%s/tools/ssc_ajax.ashx?A=%s&S=%s", req.Host, "GetImageCode", a.PlatformName)
	req, err = http.NewRequest(http.MethodPost, a.ReqUrl, strings.NewReader(r2.Form.Encode()))
	if err != nil {
		return 500, "服务器内部错误"
	}
	req.Header.Set("Origin", "http://" + req.Host)
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/69.0.3497.100 Safari/537.36")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Close = true

	client = &http.Client{
		CheckRedirect: nil,
		Jar:           cookieJar,
	}
	resp2, err := client.Do(req)
	if err != nil {
		return 501, "请求失败"
	}

	defer resp2.Body.Close()
	body2, _ := ioutil.ReadAll(resp2.Body)
	bodys1 := string(body2)
	d1 := []byte(bodys1)
	ioutil.WriteFile("test3.txt", d1, 0644)

	// 账号验证
	var r1 http.Request
	r1.ParseForm()
	r1.Form.Add("Action", "CheckUser")
	r1.Form.Add("UserName", param.Account)
	r1.Form.Add("SourceName", "PC")

	a.ReqUrl = fmt.Sprintf("http://%s/tools/ssc_ajax.ashx?A=%s&S=%s&U=%s", req.Host, "CheckUser", a.PlatformName, param.Account)
	req, err = http.NewRequest(http.MethodPost, a.ReqUrl, strings.NewReader(r1.Form.Encode()))
	if err != nil {
		return 500, "服务器内部错误"
	}
	req.Header.Set("Origin", "http://" + req.Host)
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/69.0.3497.100 Safari/537.36")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Close = true

	client = &http.Client{
		CheckRedirect: nil,
		Jar:           cookieJar,
	}
	resp1, err := client.Do(req)
	if err != nil {
		return 501, "请求失败"
	}

	defer resp1.Body.Close()
	body, _ := ioutil.ReadAll(resp1.Body)
	var m = make(map[string]interface{})
	if err = json.Unmarshal(body, &m); err != nil {
		return 501, "请求失败"
	}
	if fmt.Sprintf("%v", m["Code"]) != "0" {
		if m["StrCode"] != nil {
			return 101, fmt.Sprintf("%v", m["StrCode"])
		} else if fmt.Sprintf("%v", m["Exist"]) == "true" {
			return 101, "账号已存在"
		}
	}

	var r http.Request
	r.ParseForm()
	r.Form.Add("Action", "Register")
	r.Form.Add("InvitationCode", param.Upline)
	r.Form.Add("UserName", param.Account)
	r.Form.Add("Password", utils.Md5(param.Account + utils.Md5(param.Password)))
	r.Form.Add("Type", "Hash")
	r.Form.Add("SourceName", "PC")

	a.ReqUrl = fmt.Sprintf("http://%s/tools/ssc_ajax.ashx?A=%s&S=%s&U=%s", req.Host, "Register", a.PlatformName, param.Account)
	fmt.Println(a.ReqUrl)
	req, err = http.NewRequest(a.ReqMethod, a.ReqUrl, strings.NewReader(r.Form.Encode()))
	if err != nil {
		return 500, "服务器内部错误"
	}
	req.Header.Set("Origin", "http://" + req.Host)
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/69.0.3497.100 Safari/537.36")
	req.Header.Set("Accept", "*/*")
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
		//bodys := string(body)
		//d1 := []byte(bodys)
		//ioutil.WriteFile("test.txt", d1, 0644)

		if err = json.Unmarshal(body, &m); err != nil {
			return 501, "请求失败"
		}
		if fmt.Sprintf("%v", m["Code"]) == "0" && fmt.Sprintf("%v", m["StrCode"]) == "注册成功" {
			return 1, ""
		} else if m["StrCode"] != nil && fmt.Sprintf("%v", m["StrCode"]) != "" {
			return 101, fmt.Sprintf("%v", m["StrCode"])
		}

		return 200, "验证不通过，请检查"
	}
}
