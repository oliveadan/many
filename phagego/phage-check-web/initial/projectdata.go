package initial

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	frame "phagego/phagev2/initial"
	. "phagego/phagev2/models"
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
		{Id: 100, Pid: 0, Enabled: 1, Display: 1, Description: "会员信息管理", Url: "", Name: "会员信息管理", Icon: "#xe65e;", Sort: 100},
		{Id: 101, Pid: 100, Enabled: 1, Display: 1, Description: "查看会员信息", Url: "CheckIndexController.get", Name: "查看会员信息", Icon: "", Sort: 100},
		{Id: 102, Pid: 100, Enabled: 1, Display: 0, Description: "导入会员信息", Url: "CheckIndexController.Import", Name: "导入会员信息", Icon: "", Sort: 100},
		{Id: 103, Pid: 100, Enabled: 1, Display: 0, Description: "删除一条会员信息", Url: "CheckIndexController.Delone", Name: "删除一条会员信息", Icon: "", Sort: 100},
		{Id: 104, Pid: 100, Enabled: 1, Display: 0, Description: "批量删除会员信息", Url: "CheckIndexController.Delbatch", Name: "批量删除", Icon: "", Sort: 100},
		{Id: 105, Pid: 100, Enabled: 1, Display: 0, Description: "导出会员信息", Url: "CheckIndexController.Export", Name: "导出会员信息", Icon: "", Sort: 100},
		{Id: 106, Pid: 100, Enabled: 1, Display: 0, Description: "编辑会员信息", Url: "CheckEditController.get", Name: "编辑会员信息", Icon: "", Sort: 100},
	}
	rolePermissions := []RolePermission{
		{Id: 200, RoleId: 2, PermissionId: 100},
		{Id: 201, RoleId: 2, PermissionId: 101},
		{Id: 202, RoleId: 2, PermissionId: 102},
		{Id: 203, RoleId: 2, PermissionId: 103},
		{Id: 204, RoleId: 2, PermissionId: 104},
		{Id: 205, RoleId: 2, PermissionId: 105},
		{Id: 206, RoleId: 2, PermissionId: 106},
		{Id: 240, RoleId: 3, PermissionId: 100},
		{Id: 241, RoleId: 3, PermissionId: 101},
		{Id: 242, RoleId: 3, PermissionId: 103},
		{Id: 243, RoleId: 3, PermissionId: 104},
		{Id: 244, RoleId: 3, PermissionId: 106},
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
