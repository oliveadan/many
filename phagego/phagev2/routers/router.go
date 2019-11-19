package routers

import (
	"phagego/phagev2/controllers"
	"phagego/phagev2/controllers/syscommon"
	"phagego/phagev2/controllers/sysmanage"
	"phagego/phagev2/controllers/sysmanage/admin"
	"phagego/phagev2/controllers/sysmanage/index"
	"phagego/phagev2/controllers/sysmanage/login"
	"phagego/phagev2/controllers/sysmanage/permission"
	"phagego/phagev2/controllers/sysmanage/role"
	"phagego/phagev2/controllers/sysmanage/siteconfig"
	"phagego/phagev2/controllers/sysmanage/quicknav"

	"github.com/astaxie/beego"
	"phagego/phagev2/controllers/sysmanage/customer"
)

func init() {
	var adminRouter = beego.AppConfig.String("adminrouter")
	beego.ErrorController(&controllers.ErrorController{})
	beego.Router(adminRouter+"/sys/base", &sysmanage.BaseController{}, "get:Index")
	beego.Router(adminRouter+"/sys/index", &index.SysIndexController{})
	beego.Router("/serversysteminfo", &index.SysIndexController{}, "*:Systeminfo")

	beego.Router(adminRouter+"/syscommon/upload", &syscommon.SyscommonController{}, "post:Upload")
	beego.Router(adminRouter+"/syscommon/mailverify", &syscommon.SyscommonController{}, "post:MailVerify")

	beego.Router(adminRouter+"/admin/index", &admin.AdminIndexController{})
	beego.Router(adminRouter+"/admin/delone", &admin.AdminIndexController{}, "post:Delone")
	beego.Router(adminRouter+"/admin/locked", &admin.AdminIndexController{}, "post:Locked")
	beego.Router(adminRouter+"/admin/LoginVerify", &admin.AdminIndexController{}, "post:LoginVerify")
	beego.Router(adminRouter+"/admin/add", &admin.AdminAddController{})
	beego.Router(adminRouter+"/admin/edit", &admin.AdminEditController{})
	beego.Router(adminRouter+"/changepwd/index", &admin.ChangePwdController{})

	beego.Router(adminRouter+"/role/index", &role.RoleIndexController{})
	beego.Router(adminRouter+"/role/delone", &role.RoleIndexController{}, "post:Delone")
	beego.Router(adminRouter+"/role/add", &role.RoleAddController{})
	beego.Router(adminRouter+"/role/edit", &role.RoleEditController{})

	beego.Router(adminRouter+"/permission/index", &permission.PermissionIndexController{})
	beego.Router(adminRouter+"/permission/delone", &permission.PermissionIndexController{}, "post:Delone")
	beego.Router(adminRouter+"/permission/add", &permission.PermissionAddController{})
	beego.Router(adminRouter+"/permission/edit", &permission.PermissionEditController{})

	beego.Router(adminRouter+"/login", &login.LoginController{})
	beego.Router(adminRouter+"/loginmailverify", &login.LoginController{}, "post:LoginMailVerify")
	beego.Router(adminRouter+"/logout", &login.LoginController{}, "get:Logout")

	beego.Router(adminRouter+"/site/index", &siteconfig.SiteConfigIndexController{})
	beego.Router(adminRouter+"/site/delone", &siteconfig.SiteConfigIndexController{}, "post:Delone")
	beego.Router(adminRouter+"/site/add", &siteconfig.SiteConfigAddController{})
	beego.Router(adminRouter+"/site/edit", &siteconfig.SiteConfigEditController{})

	beego.Router(adminRouter+"/qicknav/index", &quicknav.QuickNavIndexController{})
	beego.Router(adminRouter+"/qicknav/delone", &quicknav.QuickNavIndexController{}, "post:Delone")
	beego.Router(adminRouter+"/qicknav/add", &quicknav.QuickNavAddController{})
	beego.Router(adminRouter+"/qicknav/edit", &quicknav.QuickNavEditController{})

	beego.Router(adminRouter+"/customer/index",&customer.CustomerIndexController{})
	beego.Router(adminRouter+"/customer/edit",&customer.CustomerEditController{})
}
