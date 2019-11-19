package main

import (
	"github.com/lxn/walk"
	"phagego/plugins/platform"
)

type MyMainWindow struct {
	*walk.MainWindow
	NowTV   *walk.TableView
	LoginBT *walk.PushButton
	StartBT *walk.PushButton
	NowMD   *TableDataModel
	LogTX   *walk.TextEdit
	Lv      *walk.Composite
	sbi     *walk.StatusBarItem
	sbi1    *walk.StatusBarItem
	sbi2    *walk.StatusBarItem
	path    string
	// 配置项
	PsTX       *walk.LineEdit
	PaTX       *walk.LineEdit
	PpTX       *walk.LineEdit
	PlatformCB *walk.ComboBox
	// 查询条件
	SsDE       *walk.DateEdit
	SeDE       *walk.DateEdit
	DataTypeCB *walk.ComboBox
}

type TableDataModel struct {
	walk.TableModelBase
	Items []*platform.ExportData
}

func (m *TableDataModel) RowCount() int {
	return len(m.Items)
}

func (m *TableDataModel) Value(row, col int) interface{} {
	item := m.Items[row]
	return formatItem(item, col)
}

func formatItem(item *platform.ExportData, col int) interface{} {
	switch col {
	case 0:
		return item.Id
	case 1:
		return item.Account
	case 2:
		return item.Amount
	case 3:
		return item.DataTime.Format("2006-01-02 15:04:05")
	case 4:
		return item.Remark
	default:
		return "未知数据"
	}
	return ""
}
