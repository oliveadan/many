package Export

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
	"strconv"
	"time"
)

type DafaExport struct {
	reqUrl   string // 请求域名
	password string
	Jar      *cookiejar.Jar
}

func NewDafaExport() platform.ExportAdapter {
	return &DafaExport{}
}

func (a *DafaExport) Export(param platform.ExportParam) (error, []platform.ExportData, string) {
	datas := make([]platform.ExportData, 0)
	var err error
	var dataTmp []platform.ExportData
	var total int64
	if param.DataType == "1" {
		err, dataTmp, total = a.queryLottery(param, 1)
	} else if param.DataType == "2" {
		err, dataTmp, total = a.queryCmsGame(param, 1)
	} else {
		return errors.New("未知导出类型"), nil, ""
	}

	if err != nil {
		return err, nil, ""
	}
	datas = append(datas, dataTmp...)
	if total > 300 {
		p := total / 300
		y := total % 300
		if y > 0 {
			p += 1
		}
		var i int64 = 2
		for ; i <= p; i++ {
			fmt.Println(i)
			if param.DataType == "1" {
				err, dataTmp, _ = a.queryLottery(param, i)
			} else if param.DataType == "2" {
				err, dataTmp, _ = a.queryCmsGame(param, i)
			}
			if err != nil {
				return err, nil, ""
			}
			datas = append(datas, dataTmp...)
		}
	}
	return nil, datas, ""
}

func (a *DafaExport) queryCmsGame(param platform.ExportParam, page int64) (error, []platform.ExportData, int64) {

	uris := fmt.Sprintf("userName=%s&openState=&gameCode=&issue=&roomNumber=&rountType=&recordCode=&startTime=%s&endTime=%s&pageNum=%d&pageSize=300&", param.Account, param.StartTime.Format("2006-01-02")+"%2000%3A00%3A00", param.EndTime.Format("2006-01-02")+"%2023%3A59%3A59", page)
	resp, bodys, err := a.HttpRequest(http.MethodGet, "/v1/game/getCmsGameBetDataList?"+uris, nil)
	if err != nil {
		return errors.New("访问地址/v1/game/getCmsGameBetDataList"), nil, 0
	}
	//fmt.Println(bodys)
	if resp.StatusCode != http.StatusOK {
		return errors.New("返回状态异常"), nil, 0
	}
	//fmt.Println(bodys)
	if bodys != "" {
		var m map[string]interface{}
		if err := json.Unmarshal([]byte(bodys), &m); err != nil {
			return err, nil, 0
		}
		v := m["code"]
		if vv, ok := v.(int); ok && vv != 1 {
			v = m["msg"]
			if msg, ok := v.(string); ok {
				return errors.New(msg), nil, 0
			} else {
				return errors.New("返回信息未知"), nil, 0
			}
		}
		if vv, ok := m["data"].(map[string]interface{}); !ok {
			return errors.New("数据解析失败，没有data项"), nil, 0
		} else if fmt.Sprintf("%v", vv["total"]) == "0" {
			return errors.New("没有数据需要导出"), nil, 0
		} else if vv2, ok := vv["cmsGameDataList"].([]interface{}); !ok {
			return errors.New("数据解析失败，没有cmsGameDataList项"), nil, 0
		} else {
			var data = make([]platform.ExportData, 0)
			for _, lv := range vv2 {
				switch lv.(type) {
				case map[string]interface{}:
					lvm := lv.(map[string]interface{})
					amount, err1 := strconv.ParseFloat(fmt.Sprintf("%v", lvm["bettingAmount"]), 32)
					if err1 != nil {
						return errors.New("数据解析失败，金额格式错误"), nil, 0
					}
					dime, err1 := time.ParseInLocation("2006-01-02 15:04:05", fmt.Sprintf("%v", lvm["addTime"]), time.Local)
					if err1 != nil {
						return errors.New("数据解析失败，日期格式错误"), nil, 0
					}
					//fmt.Println(reflect.TypeOf(vv["betInfoList"]))
					da := platform.ExportData{
						Id:       fmt.Sprintf("%.0f", lvm["id"]),
						Account:  fmt.Sprintf("%v", lvm["userName"]),
						DataTime: dime,
						Amount:   float32(amount),
						Remark:   fmt.Sprintf("棋牌名称:%v;场次:%v;房号:%v;局号:%v;结果:%v;投注内容:%v;派送奖金:%v;投注单号:%v;终端:%v", lvm["gameName"], lvm["roundType"], lvm["roomNumber"], lvm["issue"], lvm["openType"], lvm["bettingNumber"], lvm["state"], lvm["recordCode"], lvm["sourceName"]),
					}
					data = append(data, da)
				}
			}
			//fmt.Println(reflect.TypeOf(vv["total"]))
			var totalFloat int64
			switch vv["total"].(type) {
			case float64:
				totalFloat = int64(vv["total"].(float64))
			case float32:
				totalFloat = int64(vv["total"].(float32))
			case int64:
				totalFloat = vv["total"].(int64)
			case int32:
				totalFloat = int64(vv["total"].(int32))
			}

			return nil, data, totalFloat
		}
	}
	return errors.New("返回数据未知"), nil, 0
}

func (a *DafaExport) queryLottery(param platform.ExportParam, page int64) (error, []platform.ExportData, int64) {
	uris := fmt.Sprintf("userName=%s&lotteryCode=&openState=&startTime=%s&endTime=%s&issue=&pageNum=%d&pageSize=300&", param.Account, param.StartTime.Format("2006-01-02"), param.EndTime.Format("2006-01-02"), page)
	resp, bodys, err := a.HttpRequest(http.MethodGet, "/v1/betting/getBetDataList?"+uris, nil)
	if err != nil {
		return errors.New("访问地址/v1/betting/getBetDataList异常"), nil, 0
	}
	//fmt.Println(bodys)
	if resp.StatusCode != http.StatusOK {
		return errors.New("返回状态异常"), nil, 0
	}
	//fmt.Println(bodys)
	if bodys != "" {
		var m map[string]interface{}
		if err := json.Unmarshal([]byte(bodys), &m); err != nil {
			return err, nil, 0
		}
		v := m["code"]
		if vv, ok := v.(int); ok && vv != 1 {
			v = m["msg"]
			if msg, ok := v.(string); ok {
				return errors.New(msg), nil, 0
			} else {
				return errors.New("返回信息未知"), nil, 0
			}
		}
		if vv, ok := m["data"].(map[string]interface{}); !ok {
			return errors.New("数据解析失败，没有data项"), nil, 0
		} else if fmt.Sprintf("%v", vv["total"]) == "0" {
			return errors.New("没有数据需要导出"), nil, 0
		} else if vv2, ok := vv["betInfoList"].([]interface{}); !ok {
			return errors.New("数据解析失败，没有betInfoList项"), nil, 0
		} else {
			var data = make([]platform.ExportData, 0)
			for _, lv := range vv2 {
				switch lv.(type) {
				case map[string]interface{}:
					lvm := lv.(map[string]interface{})
					amount, err1 := strconv.ParseFloat(fmt.Sprintf("%v", lvm["betMoney"]), 32)
					if err1 != nil {
						return errors.New("数据解析失败，金额格式错误"), nil, 0
					}
					dime, err1 := time.ParseInLocation("2006-01-02 15:04:05", fmt.Sprintf("%v", lvm["gmtCreated"]), time.Local)
					if err1 != nil {
						return errors.New("数据解析失败，日期格式错误"), nil, 0
					}
					//fmt.Println(reflect.TypeOf(vv["betInfoList"]))
					da := platform.ExportData{
						Id:       fmt.Sprintf("%.0f", lvm["id"]),
						Account:  fmt.Sprintf("%v", lvm["userName"]),
						DataTime: dime,
						Amount:   float32(amount),
						Remark:   fmt.Sprintf("彩种:%v;玩法:%v;期号:%v;投注内容:%v;注数:%v;倍数:%v;状态:%v;投注类型:%v;终端:%v", lvm["lotteryName"], lvm["playName"], lvm["issue"], lvm["betNumber"], lvm["betCount"], lvm["graduation"], lvm["bonusOrState"], lvm["playType"], lvm["sourceName"]),
					}
					data = append(data, da)
				}
			}
			//fmt.Println(reflect.TypeOf(vv["total"]))
			var totalFloat int64
			switch vv["total"].(type) {
			case float64:
				totalFloat = int64(vv["total"].(float64))
			case float32:
				totalFloat = int64(vv["total"].(float32))
			case int64:
				totalFloat = vv["total"].(int64)
			case int32:
				totalFloat = int64(vv["total"].(int32))
			}

			return nil, data, totalFloat
		}
	}
	return errors.New("返回数据未知"), nil, 0
}

func (a *DafaExport) HttpRequest(method string, url string, reqbody io.Reader) (resp *http.Response, bodys string, err error) {
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
	//fmt.Println(a.reqUrl + url)
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
func (a *DafaExport) StartAndGC(p *platform.ExportReq) error {
	a.reqUrl = strings.TrimRight(strings.TrimSpace(p.ReqUrl), "/")
	a.Jar = p.Jar
	a.password = p.Password

	if !strings.HasPrefix(strings.ToLower(a.reqUrl), "http") {
		a.reqUrl = "http://" + a.reqUrl
	}
	return nil
}

func init() {
	fmt.Println("init dafa Export")
	platform.RegisterExport("dafa", NewDafaExport)
}
