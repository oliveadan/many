package routers

import (
	"github.com/astaxie/beego"
	_ "phagego/frameweb-v2/routers"
	"phagego/z6hewx/controllers/common/wechat"
	"phagego/z6hewx/controllers/common/winningnumbers"
	"phagego/z6hewx/controllers/front"
	//"github.com/astaxie/beego"
)

func init() {
	// 后台管理系统
	var adminRouter string = beego.AppConfig.String("adminrouter")

	beego.Router("/", &front.FrontController{})
	//微信账号管理
	beego.Router(adminRouter+"/wechat/index", &wechat.WechatContorller{})
	beego.Router(adminRouter+"/wechat/add", &wechat.WechatAddContorller{})
	beego.Router(adminRouter+"/wechat/edit", &wechat.WechatEditContorller{})
	beego.Router(adminRouter+"/wechat/delone", &wechat.WechatContorller{}, "post:Delone")
	beego.Router(adminRouter+"/wechat/enable", &wechat.WechatContorller{}, "post:Enable")
	//中奖号码管理
	beego.Router(adminRouter+"/winningnumbers/index", &winningnumbers.WinningnumbersContorller{})
	beego.Router(adminRouter+"/winningnumbers/add", &winningnumbers.WinningnumbersAddContorller{})
	beego.Router(adminRouter+"/winningnumbers/edit", &winningnumbers.WinningnumbersEditContorller{})
	beego.Router(adminRouter+"/winningnumbers/delone", &winningnumbers.WinningnumbersContorller{}, "post:Delone")
}
