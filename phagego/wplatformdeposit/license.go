package main

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"golang.org/x/sys/windows/registry"
	"fmt"
	"strconv"
	"phagego/common/utils"
	"log"
	"net"
	"strings"
)

type MyLicWindow struct {
	*walk.MainWindow
}

func LicenseValidate(software string) bool {
	// 注册码
	key, exists, err := registry.CreateKey(registry.CURRENT_USER, "SOFTWARE\\"+software, registry.ALL_ACCESS)
	if err != nil {
		log.Fatal(err)
	}
	defer key.Close()

	if exists {
		lic, _, _ := key.GetStringValue("license")
		isValid := validate(lic, utils.Pubsalt)
		if isValid {
			return true
		}
	}

	var inTE *walk.TextEdit

	mw := new(MyLicWindow)

	if err := (MainWindow{
		AssignTo: &mw.MainWindow,
		Title:    "请输入激活码",
		MinSize:  Size{300, 120},
		Layout:   VBox{},
		Children: []Widget{
			HSplitter{
				Children: []Widget{
					TextEdit{AssignTo: &inTE},
				},
			},
			PushButton{
				Text: "激活",
				OnClicked: func() {
					if inTE.Text() == "" {
						walk.MsgBox(mw, "提示", "请填写激活码", walk.MsgBoxIconInformation)
						return
					}
					lic := utils.Md5(inTE.Text() + utils.Pubsalt)
					isValid := validate(lic, utils.Pubsalt)
					if isValid {
						if err := key.SetStringValue("license", lic); err != nil {
							walk.MsgBox(mw, "错误", "注册失败，请重试", walk.MsgBoxIconError)
							return
						} else {
							walk.MsgBox(mw, "提示", "激活成功，请重启程序", walk.MsgBoxIconInformation)
							mw.Close()
						}
					} else {
						walk.MsgBox(mw, "错误", "激活码无效", walk.MsgBoxIconError)
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

func validate(lic string, salt string) bool {
	netInterfaces, err := net.Interfaces()
	if err != nil {
		fmt.Println("Error:", err)
		return false
	}
	for _, netInterface := range netInterfaces {
		macAddr := netInterface.HardwareAddr.String()
		if len(macAddr) == 0 {
			continue
		}
		sign := utils.Md5(strconv.FormatInt(int64(netInterface.Index), 10), salt, strings.ToUpper(netInterface.HardwareAddr.String()))
		runes := []rune(sign)
		var a string
		for _, v := range runes[8:24] {
			x := utils.AnyToDecimal(string(v), 36)
			a = a + strconv.Itoa(x+x%11)
		}
		if utils.Md5(a+salt) == lic {
			return true
		}
	}
	return false
}
