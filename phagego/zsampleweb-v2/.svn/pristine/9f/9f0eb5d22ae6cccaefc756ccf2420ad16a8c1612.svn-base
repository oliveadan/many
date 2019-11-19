package initial

import (
	frame "phagego/frameweb-v2/initial"
	//. "phagego/frameweb-v2/models"
	//"github.com/astaxie/beego/orm"
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
	/* 示例
	permisions := []Permission{
		{Id: 100, Pid: 0, Enabled: 1, Display: 1, Description: "支付管理", Url: "", Name: "支付管理", Icon: "#xe65e;", Sort: 100},
		{Id: 101, Pid: 100, Enabled: 1, Display: 1, Description: "支付配置", Url: "PaymentConfigIndexController.Get", Name: "支付配置", Icon: "", Sort: 100},
		{Id: 102, Pid: 101, Enabled: 1, Display: 0, Description: "添加支付配置", Url: "PaymentConfigAddController.Get", Name: "添加支付配置", Icon: "", Sort: 100},
		{Id: 103, Pid: 101, Enabled: 1, Display: 0, Description: "编辑支付配置", Url: "PaymentConfigEditController.Get", Name: "编辑支付配置", Icon: "", Sort: 100},
		{Id: 104, Pid: 101, Enabled: 1, Display: 0, Description: "删除支付配置", Url: "PaymentConfigIndexController.Delone", Name: "删除支付配置", Icon: "", Sort: 100},
		{Id: 105, Pid: 101, Enabled: 1, Display: 0, Description: "启用禁用支付配置", Url: "PaymentConfigIndexController.Enabled", Name: "启用禁用支付配置", Icon: "", Sort: 100},
	}
	rolePermissions := []RolePermission{
		{Id: 200, RoleId: 2, PermissionId: 100},
		{Id: 201, RoleId: 2, PermissionId: 101},
		{Id: 202, RoleId: 2, PermissionId: 102},
		{Id: 203, RoleId: 2, PermissionId: 103},
		{Id: 204, RoleId: 2, PermissionId: 104},
		{Id: 205, RoleId: 2, PermissionId: 105},
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
	*/
}
