package main

import (
	"bufio"
	"io"
	"os"
	"strings"
	"phagego/plugins/platform"
	_ "phagego/plugins/platform/login"
	_ "phagego/plugins/platform/deposit"
	"fmt"
	"github.com/lxn/walk"
	"github.com/tealeg/xlsx"
	"math"
	"sync/atomic"
	"strconv"
	"net/http/cookiejar"
	"path/filepath"
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
					mw.NowMD.Items = []*platform.Order{}
					mw.NowMD.PublishRowsReset()
					mw.NowTV.SetSelectedIndexes([]int{})
					// 选择文件
					dlg := new(walk.FileDialog)
					dlg.FilePath = mw.path
					dlg.Filter = "Excel files (*.xlsx)|*.xlsx"

					if ok, err := dlg.ShowOpen(mw); err != nil {
						lv.Appendln("打开文件选择失败！")
						break
					} else if !ok {
						lv.Appendln("取消文件选择")
						break
					}
					mw.path = dlg.FilePath
					lv.Appendln("选择文件: " + mw.path)
					// 读取excel
					xlFile, err := xlsx.OpenFile(mw.path)
					if err != nil {
						lv.Appendln("读取excel失败")
						break
					}
					if _, ok := xlFile.Sheet["Sheet1"]; !ok {
						lv.Appendln("Sheet1不存在")
						break
					}
					sheet := xlFile.Sheet["Sheet1"]
					orders := make([]*platform.Order, 0)
					for i, v := range sheet.Rows {
						if i == 0 {
							continue
						}
						if len(v.Cells) < 3 {
							lv.Appendln(fmt.Sprintf("第%d行ID、会员账号、金额必填", i+1))
							break
						}
						order := new(platform.Order)
						order.DepositType = 5
						if id, err := v.Cells[0].Int64(); err != nil {
							lv.Appendln(fmt.Sprintf("第%d行ID必须为数字", i+1))
							break
						} else {
							order.Id = id
						}
						if account := v.Cells[1].Value; account == "" {
							lv.Appendln(fmt.Sprintf("第%d行会员账号为空", i+1))
							break
						} else {
							order.Account = account
						}
						if amount, err := v.Cells[2].Float(); err != nil {
							lv.Appendln(fmt.Sprintf("第%d行金额必须为数字", i+1))
							break
						} else {
							order.Amount = float32(amount)
						}
						if len(v.Cells) > 3 {
							order.PortalMemo = v.Cells[3].Value
						}
						order.Memo = "导入ID：" + v.Cells[0].Value
						orders = append(orders, order)
						mw.NowMD.Items = append(mw.NowMD.Items, order)
					}
					mw.NowMD.PublishRowsReset()
					mw.NowTV.SetSelectedIndexes([]int{})
					// 初始化插件，加款
					pd := platform.DepositReq{
						ReqUrl:   ps,
						Jar:      jar,
						Password: pp,
					}
					ptd, err := platform.NewPlatformDeposit(pfcb, &pd)
					if err != nil {
						fmt.Println(err)
						lv.Appendln("初始化平台加款失败！")
						break
					}
					total := len(orders)
					mw.sbi1.SetText(strconv.Itoa(total))
					size := int(math.Ceil(float64(total) / 5)) // 最多用5个线程处理
					mw.sbi.SetText("导入中...")
					var count uint32 = 0
					for i := 0; i < 5; i++ {
						var end int
						if i*size+size > total {
							end = total
						} else {
							end = i*size + size
						}
						go func(j int, orders []*platform.Order) {
							fmt.Println("线程：", j)
							for k, v := range orders {
								ptd.Send(v)
								mw.NowMD.Items[j+k] = v
								mw.NowTV.UpdateItem(j + k)
								atomic.AddUint32(&count, 1)
							}
						}(i*size, orders[i*size:end])
						if end == total {
							break
						}
					}
					for {
						mw.sbi2.SetText(strconv.Itoa(int(count)))
						if count == uint32(total) {
							mw.sbi2.SetText(strconv.Itoa(int(count)))
							break
						}
					}
					mw.sbi.SetText("导入完成")
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
	if ps == "" || pa == "" || pp == "" || pfcb == "" {
		lv.Appendln("存在空配置项，请检查")
		return
	}
	ch2 <- 1
	lv.Appendln("启动中...")
}
func DownloadBTListener(mw *MyMainWindow) {
	file := xlsx.NewFile()
	sheet, err := file.AddSheet("Sheet1")
	if err != nil {
		fmt.Printf(err.Error())
		lv.Appendln("下载导入格式失败(1)")
		return
	}
	row := sheet.AddRow()
	row.AddCell().Value = "ID(必填)"
	row.AddCell().Value = "会员账号(必填)"
	row.AddCell().Value = "金额(必填,单位：元)"
	row.AddCell().Value = "备注(可空)(用户可见)"

	row1 := sheet.AddRow()
	row1.AddCell().Value = "1"
	row1.AddCell().Value = "zhang3"
	row1.AddCell().Value = "100"
	row1.AddCell().Value = "11日签到赠送"

	var appPath string
	if appPath, err = filepath.Abs(filepath.Dir(os.Args[0])); err != nil {
		fmt.Printf(err.Error())
		lv.Appendln("下载导入格式失败(2)")
		return
	}
	fPath := filepath.Join(appPath, "导入格式.xlsx")
	err = file.Save(fPath)
	if err != nil {
		fmt.Printf(err.Error())
		lv.Appendln("下载导入格式失败(3)")
	} else {
		lv.Appendln("下载成功,路径：")
		lv.Appendln(fPath)
	}
}
