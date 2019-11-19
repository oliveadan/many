package initial

import (
	"github.com/astaxie/beego/orm"
	frame "phagego/frameweb-v2/initial"
	. "phagego/frameweb-v2/models"

	"github.com/astaxie/beego"
)

// 本项目数据初始化
func InitDbProjectData() {
	// 初始化系统基础数据
	frame.InitDbFrameData()
	// 初始化项目数据
	if beego.AppConfig.DefaultInt("dbautocreate", 0) == 0 {
		return
	}
	beego.Info("Init project data")
	permisions := []Permission{
		{Id: 100, Pid: 0, Enabled: 1, Display: 1, Description: "投票管理", Url: "", Name: "投票管理", Icon: "#xe65e;", Sort: 100},
		{Id: 101, Pid: 100, Enabled: 1, Display: 1, Description: "投票列表", Url: "IndexVoteDetailController.Get", Name: "投票列表", Icon: "", Sort: 100},
		{Id: 102, Pid: 100, Enabled: 1, Display: 0, Description: "删除所有投票", Url: "IndexVoteDetailController.DelBtch", Name: "删除所有投票", Icon: "", Sort: 100},
		{Id: 103, Pid: 100, Enabled: 1, Display: 1, Description: "票数设置", Url: "IndexSetVoteController.Get", Name: "票数设置", Icon: "", Sort: 100},
		{Id: 104, Pid: 100, Enabled: 1, Display: 0, Description: "修改票数设置", Url: "EditSetVoteController.Get", Name: "修改票数设置", Icon: "", Sort: 100},
	}
	rolePermissions := []RolePermission{
		{Id: 200, RoleId: 2, PermissionId: 100},
		{Id: 201, RoleId: 2, PermissionId: 101},
		{Id: 202, RoleId: 2, PermissionId: 102},
		{Id: 203, RoleId: 2, PermissionId: 103},
		{Id: 204, RoleId: 2, PermissionId: 104},
	}
	o := orm.NewOrm()
	for _, v := range permisions {
		if _, _, err := o.ReadOrCreate(&v, "Id"); err != nil {
			beego.Error("InitProjectData Permission error", err)
		}
	}
	for _, v := range rolePermissions {
		if _, _, err := o.ReadOrCreate(&v, "Id"); err != nil {
			beego.Error("InitProjectData RolePermission error", err)
		}
	}
}
