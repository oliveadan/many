package initial

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"

	frame "phagego/frameweb-v2/initial"
	. "phagego/frameweb-v2/models"
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
		{Id: 90, Creator: 0, Modifior: 0, Version: 0, Pid: 0, Enabled: 1, Display: 1, Description: "微信管理", Url: "", Name: "微信管理", Icon: "#xe677;", Sort: 100},
		{Id: 91, Creator: 0, Modifior: 0, Version: 0, Pid: 90, Enabled: 1, Display: 1, Description: "配置微信", Url: "WechatContorller.Get", Name: "配置微信", Icon: "", Sort: 100},
		{Id: 92, Creator: 0, Modifior: 0, Version: 0, Pid: 90, Enabled: 1, Display: 0, Description: "添加微信", Url: "WechatAddContorller.Get", Name: "添加微信", Icon: "", Sort: 100},
		{Id: 93, Creator: 0, Modifior: 0, Version: 0, Pid: 90, Enabled: 1, Display: 0, Description: "修改微信", Url: "WechatEditContorller.Get", Name: "修改微信", Icon: "", Sort: 100},
		{Id: 94, Creator: 0, Modifior: 0, Version: 0, Pid: 90, Enabled: 1, Display: 0, Description: "删除微信", Url: "WechatContorller.Delone", Name: "删除微信", Icon: "", Sort: 100},
		{Id: 95, Creator: 0, Modifior: 0, Version: 0, Pid: 90, Enabled: 1, Display: 0, Description: "启用/禁用微信", Url: "WechatContorller.Enable", Name: "启用/禁用微信", Icon: "", Sort: 100},
		{Id: 110, Creator: 0, Modifior: 0, Version: 0, Pid: 0, Enabled: 1, Display: 1, Description: "中奖号码管理", Url: "", Name: "中奖号码管理", Icon: "#xe664;", Sort: 100},
		{Id: 111, Creator: 0, Modifior: 0, Version: 0, Pid: 110, Enabled: 1, Display: 1, Description: "配置中奖号码", Url: "WinningnumbersContorller.Get", Name: "配置中奖号码", Icon: "", Sort: 100},
		{Id: 112, Creator: 0, Modifior: 0, Version: 0, Pid: 110, Enabled: 1, Display: 0, Description: "添加中奖号码", Url: "WinningnumbersAddContorller.Get", Name: "添加中奖号码", Icon: "", Sort: 100},
		{Id: 113, Creator: 0, Modifior: 0, Version: 0, Pid: 110, Enabled: 1, Display: 0, Description: "修改中奖号码", Url: "WinningnumbersEditContorller.Get", Name: "修改中奖号码", Icon: "", Sort: 100},
		{Id: 114, Creator: 0, Modifior: 0, Version: 0, Pid: 110, Enabled: 1, Display: 0, Description: "删除中奖号码", Url: "WinningnumbersContorller.Delone", Name: "删除中奖号码", Icon: "", Sort: 100},
	}
	rolePermissions := []RolePermission{
		{Id: 174, RoleId: 2, PermissionId: 90},
		{Id: 175, RoleId: 2, PermissionId: 91},
		{Id: 176, RoleId: 2, PermissionId: 92},
		{Id: 177, RoleId: 2, PermissionId: 93},
		{Id: 178, RoleId: 2, PermissionId: 94},
		{Id: 179, RoleId: 2, PermissionId: 95},
		{Id: 180, RoleId: 2, PermissionId: 110},
		{Id: 181, RoleId: 2, PermissionId: 111},
		{Id: 182, RoleId: 2, PermissionId: 112},
		{Id: 183, RoleId: 2, PermissionId: 113},
		{Id: 184, RoleId: 2, PermissionId: 114},
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
