package routers

import (
	_ "phagego/phagev2/routers"

	"github.com/astaxie/beego"
	"phagego/phage-check-web/controllers/common/check"
)

func init() {
	//后台管理系统
	var adminRouter string = beego.AppConfig.String("adminrouter")
	beego.Router(adminRouter+"/check/index", &check.CheckIndexController{})
	beego.Router(adminRouter+"/check/import", &check.CheckIndexController{}, "post:Import")
	beego.Router(adminRouter+"/check/delone", &check.CheckIndexController{}, "post:Delone")
	beego.Router(adminRouter+"/check/delbatch", &check.CheckIndexController{}, "post:Delbatch")
	beego.Router(adminRouter+"/check/export", &check.CheckIndexController{}, "post:Export")
	beego.Router(adminRouter+"/check/edit", &check.CheckEditController{})
}
