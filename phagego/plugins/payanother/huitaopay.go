package payanother

import (
	"net/http"
	"crypto/tls"
	"time"
	"io/ioutil"
	"fmt"
	"encoding/json"
	"github.com/pkg/errors"
	"strconv"
	"strings"
	"net/url"
)

type HuiTaoPay struct {
}

func (a *HuiTaoPay) Pay(p *SendParam) (*SendResp, error) {
	reqUrl := "https://pay.huitaopay.com/gwapi/transfer"

	r := map[string]string{}
	r["Version"] = "V2.5.1"
	r["ServerName"] = "hcTransferPay"
	r["MerNo"] = p.MertNo
	r["ReqTime"] = time.Now().Format("2006-01-02 15:04:05")
	r["SignType"] = "MD5"

	r["TransId"] = p.OrderNo
	r["Amount"] = strconv.FormatFloat(float64(p.Amount)/100, 'f', 2, 64)
	r["BankCode"] = p.BankCode
	if p.BusType == 0 {
		r["BusType"] = "PRV"
	} else {
		r["BusType"] = "PUB"
	}
	r["AccountName"] = p.AccountName
	r["CardNo"] = p.CardNo
	r["ReturnRemark"] = p.Attach
	r["NotifyURL"] = p.NotifyURL

	r["SignInfo"] = GenMd5Sign(p.EncryptKey, &r)
	// 收款人姓名url编码
	r["AccountName"] = url.QueryEscape(r["AccountName"])

	postData := MapToUrl(r, false)
	req, err := http.NewRequest(http.MethodPost, reqUrl, strings.NewReader(postData))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		return nil, err
	}

	config := &tls.Config{InsecureSkipVerify: true}
	tr := &http.Transport{TLSClientConfig: config}
	client := &http.Client{
		Transport:     tr,
		Timeout:       15 * time.Second,
		CheckRedirect: nil,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	var rp map[string]interface{}
	if err := json.Unmarshal(body, &rp); err != nil {
		return nil, err
	}

	data := rp["RespContent"]
	if data == nil {
		return nil, errors.New(fmt.Sprintf("%v", rp["ResMsg"]))
	}
	if v, ok := data.(map[string]interface{}); ok {

		var m2 = make(map[string]string)
		for k, v1 := range v {
			m2[k] = v1.(string)
		}

		rpSign := GenMd5Sign(p.EncryptKey, &m2, "SecureCode")
		if rpSign != m2["SecureCode"] {
			return nil, errors.New("签名验证失败")
		}

		sr := SendResp{}
		switch m2["TransStatus"] {
		case "fail":
			sr.Status = 2
			sr.Msg = "代付失败"
		case "ok":
			sr.Status = 1
			sr.Msg = "代付成功"
		case "transing":
			sr.Status = 3
			sr.Msg = "银行处理中"
		}
		return &sr, nil
	}
	return nil, errors.New("返回数据格式错误")
}

func (a *HuiTaoPay) Notice(w http.ResponseWriter, r *http.Request, getMd5Key func(args ... string) string, busCallback func(args ... string) bool) error {
	err := r.ParseForm()
	if err != nil {
		return err
	}
	md5key := getMd5Key(r.FormValue("merchantId"))
	resMap := make(map[string]string)
	for k, v := range r.Form {
		resMap[k] = v[0]
	}
	sign := GenMd5Sign(md5key, &resMap)
	if sign != r.FormValue("sign") {
		return errors.New("验签失败")
	}
	status := "0"
	if r.FormValue("tradeStatus") == "success" {
		status = "1"
	} else if r.FormValue("tradeStatus") == "failure" {
		status = "2"
	}
	if isOk := busCallback(status, r.FormValue("merOrderNo")); isOk {
		w.Write([]byte("SUCCESS"))
	}

	return nil
}

/*
 * 参数：MertNo,EncryptKey
 */
func (a *HuiTaoPay) QueryBalance(p *QueryBalanceParam) (int, error) {
	reqUrl := "https://pay.huitaopay.com/gwapi/query"

	r := map[string]string{}
	r["Version"] = "V2.5.1"
	r["ServerName"] = "hcBalanceQuery"
	r["MerNo"] = p.MertNo
	r["ReqTime"] = time.Now().Format("2006-01-02 15:04:05")
	r["SignType"] = "MD5"

	r["Currency"] = "CNY"

	r["SignInfo"] = GenMd5Sign(p.EncryptKey, &r)

	postData := MapToUrl(r, false)
	req, err := http.NewRequest(http.MethodPost, reqUrl, strings.NewReader(postData))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		return 0, err
	}

	config := &tls.Config{InsecureSkipVerify: true}
	tr := &http.Transport{TLSClientConfig: config}
	client := &http.Client{
		Transport:     tr,
		Timeout:       15 * time.Second,
		CheckRedirect: nil,
	}
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	var rp map[string]interface{}
	if err := json.Unmarshal(body, &rp); err != nil {
		return 0, err
	}

	data := rp["RespContent"]
	if data == nil {
		return 0, errors.New(fmt.Sprintf("%v", rp["ResMsg"]))
	}
	if v, ok := data.(map[string]interface{}); ok {
		var m2 = make(map[string]string)
		for k, v1 := range v {
			m2[k] = v1.(string)
		}

		rpSign := GenMd5Sign(p.EncryptKey, &m2, "SignInfo")
		if rpSign != m2["SignInfo"] {
			return 0, errors.New("签名验证失败")
		}
		balance, err := strconv.ParseFloat(m2["Balance"], 32)
		if err != nil {
			return 0, err
		}
		return int(balance * 1000 / 10), nil
	}
	return 0, errors.New("返回数据格式错误")
}
