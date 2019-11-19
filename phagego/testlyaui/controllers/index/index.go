package index

import (
	"github.com/astaxie/beego"
	"strings"
)

type IndexFrontController struct {
	beego.Controller
}

func (this *IndexFrontController) Prepare() {
	this.EnableXSRF = false
}

func (this *IndexFrontController) Get() {
	controllerName, action := this.GetControllerAndAction()
	beego.Info(controllerName,action)
	name := strings.ToLower(controllerName[0 : len(controllerName)-10])
    beego.Info("~~~~~~~~",name)
	s := this.Ctx.Request.RemoteAddr
	beego.Info("ip-----",s)
	l := strings.LastIndex(s, ":")
	beego.Info(s[0:l])
	this.TplName = "index.html"
}
