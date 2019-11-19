package routers

import (
	"github.com/astaxie/beego"
	_ "phagego/frameweb-v2/routers"
	"phagego/testlyaui/controllers/index"
)

func init() {
	// 后台管理系统
	var adminRouter string = beego.AppConfig.String("adminrouter")
	beego.Router("/", &index.IndexFrontController{})
	beego.Info("~~~~~~~~~~~~~", adminRouter)
}
