package license

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"golang.org/x/sys/windows/registry"
	"fmt"
	"log"
	"phagego/plugins/license"
	"strings"
	"time"
)

type MyLicWindow struct {
	*walk.MainWindow
}

func Validate(software string, licenseSalt string) bool {
	// 注册码
	key, exists, err := registry.CreateKey(registry.CURRENT_USER, "SOFTWARE\\"+software, registry.ALL_ACCESS)
	if err != nil {
		log.Fatal(err)
	}
	defer key.Close()

	if exists {
		exp, _, err1 := key.GetStringValue("exptime")
		lic, _, err2 := key.GetStringValue("license")
		if err1 == nil && err2 == nil {
			payTime, err := time.ParseInLocation("20060102150405", exp, time.Local)
			if err == nil {
				ok, _ := license.CheckLicense(lic, payTime, true, licenseSalt)
				if ok {
					return true
				}
			}
		}
	}

	var inTE *walk.LineEdit
	var tx *walk.TextEdit

	mw := new(MyLicWindow)

	if err := (MainWindow{
		AssignTo: &mw.MainWindow,
		Title:    "请输入激活码",
		MinSize:  Size{400, 320},
		Layout:   VBox{},
		Children: []Widget{
			LineEdit{AssignTo: &inTE, RowSpan: 1},
			TextEdit{AssignTo: &tx, RowSpan: 3, ReadOnly: true, Text: fmt.Sprintf("%v", "请全选复制后发送给管理员生成激活码\r\n"+strings.Join(license.GetMachineData(licenseSalt), "\r\n"))},
			PushButton{
				RowSpan: 1,
				Text:    "激活",
				OnClicked: func() {
					if inTE.Text() == "" {
						walk.MsgBox(mw, "提示", "请填写激活码", walk.MsgBoxIconInformation)
						return
					}
					ac := strings.TrimSpace(inTE.Text())
					ok, lic, expTime := license.GenLicense(ac, true, licenseSalt)
					if !ok {
						walk.MsgBox(mw, "错误", lic, walk.MsgBoxIconError)
						return
					}
					ok, msg := license.CheckLicense(lic, expTime, true, licenseSalt)
					if ok {
						if err := key.SetStringValue("exptime", expTime.Format("20060102150405")); err != nil {
							walk.MsgBox(mw, "错误", "注册失败，请重试", walk.MsgBoxIconError)
							return
						}
						if err := key.SetStringValue("license", lic); err != nil {
							walk.MsgBox(mw, "错误", "注册失败，请重试", walk.MsgBoxIconError)
							return
						}
						walk.MsgBox(mw, "提示", "激活成功，请重启程序", walk.MsgBoxIconInformation)
						mw.Close()
					} else {
						walk.MsgBox(mw, "错误", msg, walk.MsgBoxIconError)
						return
					}
				},
			},
		},
	}.Create()); err != nil {
		log.Fatal(err)
	}
	mw.Run()
	return false
}
