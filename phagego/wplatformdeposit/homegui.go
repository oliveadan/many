package main

import (
	"log"

	. "github.com/lxn/walk/declarative"
	"github.com/lxn/walk"
)

// 配置变量
var ps, pa, pp, pfcb string
var lv *LogView

func NewHomeGui() *MyMainWindow {
	mw := &MyMainWindow{NowMD: new(OrderModel)}

	if err := (MainWindow{
		AssignTo: &mw.MainWindow,
		Title:    "批量加款程序",
		MinSize:  Size{1000, 600},
		Layout:   HBox{},
		Children: []Widget{
			VSplitter{
				MaxSize: Size{Width: 300, Height: 800},
				Children: []Widget{
					Composite{
						Layout: VBox{},
						Children: []Widget{
							Composite{
								Layout: Grid{Columns: 1, MarginsZero: true},
								Children: []Widget{
									Label{Text: "后台类型:"},
									ComboBox{
										AssignTo: &mw.PlatformCB,
										Model:    []string{"tbk","gpk"},
									},

									Label{Text: "后台域名:"},
									LineEdit{
										AssignTo: &mw.PsTX,
										MinSize:  Size{Width: 100},
										//Text:     "http://qs.ia001.17gpk.com",
									},
									Label{Text: "账号:"},
									LineEdit{
										AssignTo: &mw.PaTX,
										//Text:     "testIT",
									},
									Label{Text: "密码:"},
									LineEdit{
										AssignTo:     &mw.PpTX,
										PasswordMode: true,
										Text:     "123456",
									},
								},
							},
							Composite{
								Layout: Grid{Columns: 3, MarginsZero: true},
								Children: []Widget{
									PushButton{
										AssignTo: &mw.LoginBT,
										Text:     "登录",
									},
									PushButton{
										AssignTo: &mw.StartBT,
										Text:     "导入",
									},
									PushButton{
										AssignTo: &mw.DownloadBT,
										Text:     "下载格式",
									},
								},
							},
						},
					},
					Composite{
						Layout:   VBox{},
						AssignTo: &mw.Lv,
					},
				},
			},
			TableView{
				AssignTo: &mw.NowTV,
				Columns: []TableViewColumn{
					{Title: "ID", Width: 80},
					{Title: "账号", Width: 120},
					{Title: "金额", Width: 80},
					{Title: "备注", Width: 100},
					{Title: "处理状态", Width: 100},
					{Title: "处理结果"},
				},
				StyleCell: func(style *walk.CellStyle) {
					item := mw.NowMD.Items[style.Row()]
					if item.Status == 2 {
						style.BackgroundColor = walk.RGB(250,128,114)
					}
				},
				Model: mw.NowMD,
			},
		},
		StatusBarItems: []StatusBarItem{
			StatusBarItem{
				AssignTo: &mw.sbi,
				Text:     "准备就绪",
				Width:    100,
			},
			StatusBarItem{
				AssignTo: &mw.sbi1,
				Text:     "",
				Width:    80,
			},
			StatusBarItem{
				AssignTo: &mw.sbi2,
				Text:     "",
				Width:    80,
			},
		},
	}.Create()); err != nil {
		log.Fatal(err)
	}
	var err error
	if lv, err = NewLogView(mw.Lv); err != nil {
		log.Fatal(err)
	}
	lv.Appendln("日志：")

	mw.LoginBT.Clicked().Attach(func() {
		LoginBTListener(mw)
	})
	mw.StartBT.Clicked().Attach(func() {
		StartBTListener(mw)
	})
	mw.DownloadBT.Clicked().Attach(func() {
		DownloadBTListener(mw)
	})

	return mw
}
