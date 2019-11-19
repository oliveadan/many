package main

import (
	"log"
	. "github.com/lxn/walk/declarative"
	"time"
	"phagego/framewin/logview"
)

// 配置变量
var ps, pa, pp, pfcb string
var ss, se time.Time
var cbt string
var lv *logview.LogView

func NewHomeGui() *MyMainWindow {
	mw := &MyMainWindow{NowMD: new(TableDataModel)}

	if err := (MainWindow{
		AssignTo: &mw.MainWindow,
		Title:    "数据导出程序",
		MinSize:  Size{1000, 800},
		Layout:   HBox{},
		Children: []Widget{
			VSplitter{
				MaxSize: Size{Width: 320, Height: 600},
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
										Model:    []string{"dafa"},
									},

									Label{Text: "后台域名:"},
									LineEdit{
										AssignTo: &mw.PsTX,
										MinSize:  Size{Width: 100},
									},
									Label{Text: "账号:"},
									LineEdit{
										AssignTo: &mw.PaTX,
									},
									Label{Text: "密码:"},
									LineEdit{
										AssignTo:     &mw.PpTX,
										PasswordMode: true,
										Text:         "",
									},
									PushButton{
										AssignTo: &mw.LoginBT,
										Text:     "登录",
									},
								},
							},
							Composite{
								Layout: Grid{Columns: 2, MarginsZero: true},
								Children: []Widget{
									Label{Text: "开始日期:"},
									Label{Text: "结束日期:"},
									DateEdit{
										AssignTo: &mw.SsDE,
									},
									DateEdit{
										AssignTo: &mw.SeDE,
									},
									ComboBox{
										AssignTo: &mw.DataTypeCB,
										Model:    []string{"1-彩票投注", "2-棋牌投注"},
									},
									PushButton{
										AssignTo: &mw.StartBT,
										Text:     "导出",
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
					{Title: "ID", Width: 100},
					{Title: "账号", Width: 120},
					{Title: "金额", Width: 80},
					{Title: "时间", Width: 120},
					{Title: "备注", Width: 200},
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
	if lv, err = logview.NewLogView(mw.Lv); err != nil {
		log.Fatal(err)
	}
	lv.Appendln("日志：")

	mw.LoginBT.Clicked().Attach(func() {
		LoginBTListener(mw)
	})
	mw.StartBT.Clicked().Attach(func() {
		StartBTListener(mw)
	})

	return mw
}
