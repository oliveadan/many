package main

import (
	"github.com/lxn/walk"
	"phagego/plugins/platform"
)

type MyMainWindow struct {
	*walk.MainWindow
	NowTV      *walk.TableView
	LoginBT    *walk.PushButton
	StartBT    *walk.PushButton
	DownloadBT *walk.PushButton
	NowMD      *OrderModel
	LogTX      *walk.TextEdit
	Lv         *walk.Composite
	sbi        *walk.StatusBarItem
	sbi1       *walk.StatusBarItem
	sbi2       *walk.StatusBarItem
	path       string
	// 配置项
	PsTX       *walk.LineEdit
	PaTX       *walk.LineEdit
	PpTX       *walk.LineEdit
	PlatformCB *walk.ComboBox
}

type OrderModel struct {
	walk.TableModelBase
	Items []*platform.Order
}

func (m *OrderModel) RowCount() int {
	return len(m.Items)
}

func (m *OrderModel) Value(row, col int) interface{} {
	item := m.Items[row]
	return formatItem(item, col)
}

func formatItem(item *platform.Order, col int) interface{} {
	switch col {
	case 0:
		return item.Id
	case 1:
		return item.Account
	case 2:
		return item.Amount
	case 3:
		return item.PortalMemo
	case 4:
		switch item.Status {
		case 0:
			return "未处理"
		case 1:
			return "成功"
		default:
			return "失败"
		}
	case 5:
		if item.Msg != "" {
			return item.Msg
		}
	default:
		return "未知数据"
	}
	return ""
}
