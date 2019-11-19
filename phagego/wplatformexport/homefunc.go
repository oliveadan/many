package main

import (
	"bufio"
	"io"
	"os"
	"strings"
	"phagego/plugins/platform"
	_ "phagego/plugins/platform/login"
	_ "phagego/plugins/platform/export"
	"github.com/tealeg/xlsx"
	"fmt"
	"net/http/cookiejar"
	"path/filepath"
	"time"
)

var ch2 = make(chan int)
var jar *cookiejar.Jar

func InitHomeGui(mw *MyMainWindow) {
	// 读取缓存
	fi, err := os.Open("cache.store")
	if err == nil {
		defer fi.Close()
		br := bufio.NewReader(fi)
		for {
			a, _, c := br.ReadLine()
			if c == io.EOF {
				break
			}
			arr := strings.Split(string(a), "=")
			if len(arr) != 2 {
				continue
			}
			switch arr[0] {
			case "ps":
				mw.PsTX.SetText(arr[1])
			case "pa":
				mw.PaTX.SetText(arr[1])
			case "pfcb":
				plcbm := mw.PlatformCB.Model()
				if plcbarr, ok := plcbm.([]string); ok {
					for i, v := range plcbarr {
						if v == arr[1] {
							mw.PlatformCB.SetCurrentIndex(i)
						}
					}
				}
			}
		}
	}
	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println(r)
			}
		}()
		for {
			select {
			case v := <-ch2:
				switch v {
				case 0:
					// 登录平台
					p := platform.LoginReq{
						ReqUrl:   ps,
						Username: pa,
						Password: pp,
					}
					pt, err := platform.NewPlatformLogin(pfcb, &p)
					if err != nil {
						lv.Appendln("初始化平台失败！")
						break
					}
					jar, err = pt.Login()
					if err != nil {
						lv.Appendln(err.Error())
						break
					} else {
						// 写入缓存
						var appPath string
						if appPath, err = filepath.Abs(filepath.Dir(os.Args[0])); err != nil {
							panic(err)
						}
						fmt.Println(appPath)
						fPath := filepath.Join(appPath, "cache.store")
						f, err := os.OpenFile(fPath, os.O_APPEND|os.O_CREATE, 0644)
						if err != nil {
							lv.Appendln("缓存失败！")
						} else {
							os.Truncate(fPath, 0)
							f.WriteString("ps=" + ps + "\r\n")
							f.WriteString("pa=" + pa + "\r\n")
							f.WriteString("pfcb=" + pfcb + "\r\n")
							f.Close()
							lv.Appendln("缓存成功！")
						}
						lv.Appendln("登录成功！")
					}
				case 1:
					// 验证是否已登录
					p := platform.LoginReq{
						ReqUrl: ps,
						Jar:    jar,
					}
					pt, err := platform.NewPlatformLogin(pfcb, &p)
					if err != nil {
						lv.Appendln("初始化平台失败！")
						break
					}
					err = pt.CheckLogon()
					if err != nil {
						lv.Appendln("请重新登录")
						break
					}
					// 清空上一次记录
					mw.NowMD.Items = []*platform.ExportData{}
					mw.NowMD.PublishRowsReset()
					mw.NowTV.SetSelectedIndexes([]int{})
					// 请求数据
					pd := platform.ExportReq{
						ReqUrl:   ps,
						Jar:      jar,
						Password: pp,
					}
					ptd, err := platform.NewPlatformExport(pfcb, &pd)
					if err != nil {
						fmt.Println(err)
						lv.Appendln("平台导出初始化失败！")
						break
					}
					param := platform.ExportParam{
						DataType: cbt,
						StartTime: ss,
						EndTime:   se,
					}
					lv.Appendln("启动完成，开始检出数据...")
					err, datas, msg := ptd.Export(param)
					if err != nil {
						lv.Appendln(err.Error())
						break
					}
					if msg != "" {
						lv.Appendln(err.Error())
					}
					for _, v := range datas {
						model := new(platform.ExportData)
						model.Id = v.Id
						model.Account = v.Account
						model.Amount = v.Amount
						model.DataTime = v.DataTime
						model.Remark = v.Remark

						mw.NowMD.Items = append(mw.NowMD.Items, model)
					}
					mw.NowMD.PublishRowsReset()
					mw.NowTV.SetSelectedIndexes([]int{})
					lv.Appendln("数据查询成功, 正在生成excel...")

					file := xlsx.NewFile()
					sheet, err := file.AddSheet("Sheet1")
					if err != nil {
						fmt.Printf(err.Error())
						lv.Appendln("Excel生成失败(1)")
						return
					}
					row := sheet.AddRow()
					row.AddCell().Value = "ID"
					row.AddCell().Value = "会员账号"
					row.AddCell().Value = "金额"
					row.AddCell().Value = "时间"
					row.AddCell().Value = "备注"

					for _, v := range datas {
						row1 := sheet.AddRow()
						row1.AddCell().Value = v.Id
						row1.AddCell().Value = v.Account
						row1.AddCell().Value = fmt.Sprintf("%.2f", v.Amount)
						row1.AddCell().Value = v.DataTime.Format("2006-01-02 15:04:05")
						row1.AddCell().Value = v.Remark
					}

					var appPath string
					if appPath, err = filepath.Abs(filepath.Dir(os.Args[0])); err != nil {
						fmt.Printf(err.Error())
						lv.Appendln("Excel地址获取失败(2)")
						return
					}
					fPath := filepath.Join(appPath, time.Now().Format("20060102150405") + ".xlsx")
					err = file.Save(fPath)
					if err != nil {
						fmt.Printf(err.Error())
						lv.Appendln("Excel保存失败(3)")
					} else {
						lv.Appendln("导出成功,Excel保存路径：")
						lv.Appendln(fPath)
					}

					mw.sbi.SetText("导出完成")
				}
			}
		}
	}()
}

func LoginBTListener(mw *MyMainWindow) {
	ps = strings.TrimSuffix(strings.TrimSpace(mw.PsTX.Text()), "/")
	pa = strings.TrimSpace(mw.PaTX.Text())
	pp = strings.TrimSpace(mw.PpTX.Text())
	pfcb = mw.PlatformCB.Text()
	if ps == "" || pa == "" || pp == "" || pfcb == "" {
		lv.Appendln("存在空配置项，请检查")
		return
	}
	ch2 <- 0
	lv.Appendln("登录中...")
}
func StartBTListener(mw *MyMainWindow) {
	ps = strings.TrimSuffix(strings.TrimSpace(mw.PsTX.Text()), "/")
	pa = strings.TrimSpace(mw.PaTX.Text())
	pp = strings.TrimSpace(mw.PpTX.Text())
	pfcb = mw.PlatformCB.Text()
	ss = mw.SsDE.Date()
	se = mw.SeDE.Date()
	cbt = mw.DataTypeCB.Text()
	cbts := strings.Split(cbt, "-")
	cbt = cbts[0]
	if ps == "" || pa == "" || pp == "" || pfcb == "" {
		lv.Appendln("存在空配置项，请检查")
		return
	}
	ch2 <- 1
	lv.Appendln("启动中...")
}
